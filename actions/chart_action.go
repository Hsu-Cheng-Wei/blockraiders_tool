package actions

type ChartAction struct {
	Name          string
	ConfigurePath string
	IngressPath   string
	DockerImage   string
	ReInstall     string
	Count         int
}
