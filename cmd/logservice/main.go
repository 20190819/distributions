package main

import (
	"context"
	"distributions/log"
	"distributions/registry"
	"distributions/service"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	serverAddr := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceUrl:       serverAddr,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateUrl: serverAddr + "/services",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()

	fmt.Println("shuting down log_service")
}
