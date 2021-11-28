package main

import (
	"context"
	"distributions/registry"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/services", &registry.RegistrationService{})
	var srv http.Server
	srv.Addr = registry.ExportServerPort
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Println("Registry service Started. Press any key to shutdown.")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()

	fmt.Println("Shuting down registry service.")
}
