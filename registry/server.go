package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3000"
const GetServersUrl = "http://localhost" + ServerPort + "/services"

type registry struct {
	registations []Registration
	mutex        *sync.Mutex
}

func (r registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registations = append(r.registations, reg)
	r.mutex.Unlock()
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

	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
