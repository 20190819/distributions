package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, serviceName, host, port string,
	registerHandlerFunc func()) (context.Context, error) {
	registerHandlerFunc()
	ctx = startService(ctx, serviceName, host, port)
	return ctx, nil
}

func startService(ctx context.Context, serviceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	srv := http.Server{
		Addr: host + ":" + port,
	}

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to shutdown...", serviceName)
		var s string
		fmt.Scanln(&s)
		cancel()

	}()
	return ctx
}
