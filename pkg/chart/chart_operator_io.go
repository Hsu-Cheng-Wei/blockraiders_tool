package chart

import (
	"blockraiders_tool/templates"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func (chart *Operator[T]) OutputChartByIndex(index int, path string) error {
	err := chart.GetRelease()

	if err != nil {
		return err
	}

	rel, exist := chart.indexMapping[index]

	if !exist {
		return errors.New("index not exist")
	}

	return outPutChart(rel.Chart, path)
}

func (chart *Operator[T]) OutputChartByName(name string, path string) error {
	err := chart.GetRelease()

	if err != nil {
		return err
	}

	rel, exist := chart.nameMapping[name]

	if !exist {
		return errors.New("index not exist")
	}

	return outPutChart(rel.Chart, path)
}

func checkIfDirNotExistCreateNew(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, 0755)
	}
}

func checkIfFileNotExistCreateNew(filePath string) (*os.File, bool, error) {
	if _, err := os.Stat(filePath); os.IsExist(err) {
		fmt.Print("File:" + filePath + " exist do you want to overwrite it?[yes/no]:")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := scanner.Text()

		if strings.ToLower(answer) == "yes" {
			_ = os.Remove(filePath)
		} else {
			return nil, false, nil
		}
	}
	f, err := os.Create(filePath)

	if err != nil {
		return nil, true, err
	}
	return f, true, nil
}

func createFileRecursive(path string, files []*chart.File) error {
	for _, t := range files {
		chartFilePath := filepath.Join(path, t.Name)
		parent := filepath.Dir(chartFilePath)
		checkIfDirNotExistCreateNew(parent)

		f, isNew, err := checkIfFileNotExistCreateNew(chartFilePath)

		if !isNew {
			continue
		}

		if err != nil {
			return err
		}

		_, _ = f.Write(t.Data)
		_ = f.Sync()
		_ = f.Close()
	}
	return nil
}

func outPutChart(chart *chart.Chart, path string) error {
	checkIfDirNotExistCreateNew(path)

	basePath := filepath.Join(path, chart.Name()+"-chart")
	checkIfDirNotExistCreateNew(basePath)

	if err := createFileRecursive(basePath, chart.Templates); err != nil {
		return err
	}

	valuesPath := filepath.Join(basePath, "values.yaml")

	f, isNew, err := checkIfFileNotExistCreateNew(valuesPath)

	if isNew && err != nil {
		return err
	}

	b, _ := yaml.Marshal(chart.Values)
	_, _ = f.Write(b)
	_ = f.Sync()
	_ = f.Close()

	chartPath := filepath.Join(basePath, "Chart.yaml")
	f, isNew, err = checkIfFileNotExistCreateNew(chartPath)

	if isNew && err != nil {
		return err
	}

	b, _ = json.Marshal(*chart.Metadata)

	m := make(map[string]interface{})
	_ = json.Unmarshal(b, &m)

	tmp, err := template.New("chart").Parse(templates.ChartYaml)

	return tmp.Execute(f, m)
}
