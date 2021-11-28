package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceUrl       string
	Required         []ServiceName
	ServiceUpdateUrl string
}

type ServiceName string

const (
	logService   = ServiceName("logService")
	gradeService = ServiceName("gradeService")
)

type patchEntry struct {
	Name ServiceName
	Url  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
