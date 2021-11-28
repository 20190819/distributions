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
	r := registry.Registration{
		ServiceName: "log Service",
		ServiceUrl:  fmt.Sprintf("http://%s:%s", host, port),
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
