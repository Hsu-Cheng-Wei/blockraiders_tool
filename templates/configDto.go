package templates

import (
	"os"
	"path/filepath"
)

type ConfigDto struct {
	ChartInfo   *InitChartInfo `json:"chartInfo"`
	CodegenInfo *InitCodegen   `json:"codegenInfo"`
}

//goland:noinspection ALL
func (cfg *ConfigDto) InitChart() error {
	charts := make(map[string]*InitChart)
	charts["staging"] = &InitChart{
		Name:               "blockraiders-staging-game",
		Namespace:          "blockraiders-staging-game",
		ApiRepository:      "718286959245.dkr.ecr.ap-south-1.amazonaws.com/blockraiders-staging-game-api",
		GameHostRepository: "718286959245.dkr.ecr.ap-south-1.amazonaws.com/blockraiders-staging-game-host",
		ApiTemplate:        "blockraiders-api.tar.tgz",
		GameHostTemplate:   "blockraiders-host.tar.tgz",
		Region:             "ap-south-1",
	}
	charts["prod"] = &InitChart{
		Name:      "blockraiders-prod-game",
		Namespace: "blockraiders-prod-game",
	}

	cfg.ChartInfo = &InitChartInfo{
		Current: "staging",
		Charts:  charts,
	}

	return nil
}

func (cfg *ConfigDto) InitCodegen() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	dir := filepath.Base(path)
	cfg.CodegenInfo = &InitCodegen{
		Namespace:       dir,
		QueryPrefix:     "Query",
		CommandPrefix:   "Command",
		ApplicationPath: "Applications",
	}

	return nil
}
