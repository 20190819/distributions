package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegistyServiceHandler(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
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
