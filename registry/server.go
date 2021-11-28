package registry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const ExportServerPort = ":3000"
const ExportServersUrl = "http://localhost" + ExportServerPort + "/services"

type registry struct {
	registations []Registration
	mutex        *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registations = append(r.registations, reg)
	r.mutex.Unlock()
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
	mutex:        new(sync.Mutex),
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
