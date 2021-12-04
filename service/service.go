package service

import (
	"context"
	"distributions/registry"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, host, port string, r registry.Registration,
	registerHandlerFunc func()) (context.Context, error) {
	registerHandlerFunc()
	// 启动服务
	ctx = startService(ctx, r.ServiceName, host, port)
	// 注册服务
	registry.RegistyServiceHandler(r)
	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	srv := http.Server{
		Addr: host + ":" + port,
	}
	go func() {
		log.Println(srv.ListenAndServe())
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()
	go func() {
		fmt.Printf("启动服务 %v started. Press any key to shutdown...\n", serviceName)
		var str string
		fmt.Scanln(&str)
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()

	}()
	return ctx
}
