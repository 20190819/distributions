package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegistyServiceHandler(r Registration) error {
	serviceUpdateUrl, err := url.Parse(r.ServiceUpdateUrl)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateUrl.Path, &serviceUpdateHandler{})
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ExportServersUrl, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registry service with status code:%v", res.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

func (supdate *serviceUpdateHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
	}
	prov.Update(p)
}

func ShutdownServiceHandler(url string) error {
	buf := bytes.NewBuffer([]byte(url))
	req, err := http.NewRequest(http.MethodDelete, ExportServersUrl, buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to remove serivce. Registry service responsed with code:%v", res.StatusCode)
	}
	return nil
}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

func (p *providers) Update(pat patch) {
	p.mutex.RLock()
	defer p.mutex.Unlock()
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; ok {
			if !ok {
				p.services[patchEntry.Name] = make([]string, 0)
			}
		} else {
			p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.Url)
		}
	}
	for _, patchEntry := range pat.Removed {
		if patchUrls, ok := p.services[patchEntry.Name]; ok {
			for i := range patchUrls {
				if patchUrls[i] == patchEntry.Url {
					p.services[patchEntry.Name] = append(patchUrls[:i], patchUrls[i+1:]...)
				}
			}
		}
	}
}

func (p providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("No providers for service %v", name)
	}
	idx := int(rand.Float32() * float32(len(providers)))
	return providers[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

var prov = providers{
	services: make(map[ServiceName][]string, 0),
	mutex:    new(sync.RWMutex),
}
