package chart

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"sort"
	"strconv"
)

func (chart *Operator[T]) add(rel *release.Release) {
	chart.index++

	chart.nameMapping[rel.Name] = rel
	chart.indexMapping[chart.index] = rel
}

func (chart *Operator[T]) GetName(index int) (string, error) {
	ch, exist := chart.indexMapping[index]
	if !exist {
		return "", fmt.Errorf("Index: %d not found", index)
	}

	return ch.Name, nil
}

func (chart *Operator[T]) GetRelease() error {
	if chart.HasGetRelease {
		return nil
	}

	chart.HasGetRelease = true
	client := action.NewList(chart.Configuration)
	// Only list deployed
	client.Deployed = true
	results, err := client.Run()
	if err != nil {
		return err
	}

	for _, rel := range results {
		chart.add(rel)
	}

	return nil
}

func (chart *Operator[T]) GetAndPrintRelease() error {
	err := chart.GetRelease()

	if err != nil {
		return err
	}

	keys := make([]int, 0, len(chart.indexMapping))
	for k := range chart.indexMapping {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		rel := chart.indexMapping[k]
		fmt.Println("(" + strconv.Itoa(k) + ")" + rel.Name)
	}

	return nil
}
