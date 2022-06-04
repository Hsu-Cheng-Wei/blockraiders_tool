package chart

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
)

func (chart *Operator[T]) UpdateVersionByIndex(index int, version string) error {
	if err := chart.GetRelease(); err != nil {
		return err
	}

	rel, exist := chart.indexMapping[index]

	if !exist {
		return fmt.Errorf("index: %d doesn't exist", index)
	}

	rel.Chart.Metadata.AppVersion = version
	upgrade := action.NewUpgrade(chart.Configuration)

	_, err := upgrade.Run(rel.Name, rel.Chart, rel.Chart.Values)

	return err
}

func (chart *Operator[T]) UpdateVersionByName(name string, version string) error {
	if err := chart.GetRelease(); err != nil {
		return err
	}

	rel, exist := chart.nameMapping[name]

	if !exist {
		return fmt.Errorf("name: %s doesn't exist", name)
	}

	rel.Chart.Metadata.AppVersion = version
	upgrade := action.NewUpgrade(chart.Configuration)

	_, err := upgrade.Run(rel.Name, rel.Chart, rel.Chart.Values)

	return err
}
