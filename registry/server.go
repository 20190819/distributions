package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const ExportServerPort = ":3000"
const ExportServersUrl = "http://localhost" + ExportServerPort + "/services"

type registry struct {
	registations []Registration
	mutex        *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	// 注册服务
	r.mutex.Lock()
	r.registations = append(r.registations, reg)
	r.mutex.Unlock()
	// 加载依赖的服务
	err := r.sendRequiredService(reg)
	if err != nil {
		return err
	}
	return nil
}

func (r registry) sendRequiredService(reg Registration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var p patch
	for _, serviceReg := range r.registations {
		for _, serviceReq := range reg.RequiredServices {
			if serviceReg.ServiceName == serviceReq {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					Url:  serviceReg.ServiceUrl,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateUrl)
	if err != nil {
		return err
	}
	return nil
}
func (r registry) sendPatch(p patch, url string) error {
	pJson, err := json.Marshal(p)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(pJson))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send patch with code:%v", resp.StatusCode)
	}
	return nil
}

func (r *registry) remove(url string) error {
	for k, item := range r.registations {
		if item.ServiceUrl == url {
			r.registations = append(r.registations[:k], r.registations[k+1:]...)
		}
	}
	return nil
}

var reg = registry{
	registations: make([]Registration, 0),
	mutex:        new(sync.RWMutex),
}

type RegistrationService struct{}

func (s RegistrationService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println("Request Received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding Service:%v with Url:%s\n", r.ServiceName, r.ServiceUrl)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Removing Service with url:%s", string(payload))
		err = reg.remove(string(payload))
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
