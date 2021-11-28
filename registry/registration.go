package registry

type ServiceName string

type Registration struct {
	ServiceName ServiceName
	ServiceUrl  string
}

const (
	logService = ServiceName("logService")
)
