package chart

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"log"
	"os"
)

type Operator[T interface{}] struct {
	Configuration *action.Configuration
	CliSettings   *cli.EnvSettings

	nameMapping   map[string]*release.Release
	indexMapping  map[int]*release.Release
	index         int
	HasGetRelease bool
}

func NewChartOperator[T any](namespace string) (*Operator[T], error) {
	res := &Operator[T]{}
	err := initChartOperator(res, namespace)

	return res, err
}

func initChartOperator[T interface{}](chart *Operator[T], namespace string) error {

	chart.index = 0
	chart.nameMapping = make(map[string]*release.Release)
	chart.indexMapping = make(map[int]*release.Release)

	chart.CliSettings = cli.New()
	chart.CliSettings.SetNamespace(namespace)
	chart.Configuration = new(action.Configuration)
	if err := chart.Configuration.Init(chart.CliSettings.RESTClientGetter(), chart.CliSettings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return err
	}
	return nil
}
