package main

import (
	"context"
	"distributions/log"
	"distributions/service"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	ctx, err := service.Start(
		context.Background(),
		"log_service",
		host,
		port,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()

	fmt.Println("shuting down log_service")
}
