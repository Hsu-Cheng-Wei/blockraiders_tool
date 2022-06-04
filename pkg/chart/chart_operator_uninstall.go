package chart

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
)

func (chart *Operator[T]) UninstallByIndex(index int) error {
	if err := chart.GetRelease(); err != nil {
		return err
	}

	rel, exist := chart.indexMapping[index]

	if !exist {
		return fmt.Errorf("index %d doesn't exist", index)
	}
	act := action.NewUninstall(chart.Configuration)

	_, err := act.Run(rel.Name)

	return err
}

func (chart *Operator[T]) UninstallByName(name string) error {
	if err := chart.GetRelease(); err != nil {
		return err
	}

	rel, exist := chart.nameMapping[name]

	if !exist {
		return fmt.Errorf("name: %s doesn't exist", name)
	}
	act := action.NewUninstall(chart.Configuration)

	_, err := act.Run(rel.Name)

	return err
}
