package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceUrl       string
	RequiredServices []ServiceName
	ServiceUpdateUrl string
	HeartbeatUrl     string
}

type ServiceName string

const (
	LogService     = ServiceName("logService")
	GradingService = ServiceName("gradingService")
)

type patchEntry struct {
	Name ServiceName
	Url  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
