package main

import (
	"context"
	"distributions/grades"
	"distributions/log"
	"distributions/registry"
	"distributions/service"
	"fmt"
	stlog "log"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	r := registry.Registration{
		ServiceName:      registry.GradingService,
		ServiceUrl:       serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateUrl: serviceAddress + "/services",
		// HeartbeatURL:     serviceAddress + "/heartbeat",
	}
	ctx, err := service.Start(context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	// 获取依赖的服务
	logProvider, err := registry.GetProvider(registry.LogService)
	if err != nil {
		stlog.Fatal(err)
	} else {
		fmt.Printf("依赖服务发现 Logging service found at: %s\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
