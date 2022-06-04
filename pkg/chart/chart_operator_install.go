package chart

import (
	"encoding/json"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func (chart *Operator[T]) Install(chartPath, name, version string, valueFunc func(values *T)) error {
	ch, err := loader.Load(chartPath)

	if err != nil {
		return err
	}

	dto, err := toDto[T](ch.Values)

	if err != nil {
		return err
	}

	valueFunc(dto)

	ch.Metadata.Name = name
	ch.Values, err = toValues(dto)
	ch.Metadata.AppVersion = version

	act := action.NewInstall(chart.Configuration)
	act.ReleaseName = name
	act.Namespace = "default"
	_, err = act.Run(ch, map[string]interface{}{})

	fmt.Println("chart installed..")
	return err
}

func toDto[T interface{}](values map[string]interface{}) (*T, error) {
	b, err := json.Marshal(values)

	if err != nil {
		return nil, err
	}

	var res T

	err = json.Unmarshal(b, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func toValues(dto interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(dto)

	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})

	err = json.Unmarshal(b, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
