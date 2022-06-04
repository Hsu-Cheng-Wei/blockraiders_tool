package templates

type InitChart struct {
	Name               string `json:"name"`
	Namespace          string `json:"namespace"`
	ApiRepository      string `json:"apiRepository"`
	GameHostRepository string `json:"gameHostRepository"`
	ApiTemplate        string `json:"apiTemplate"`
	GameHostTemplate   string `json:"gameHostTemplate"`
	Region             string `json:"region"`
}

type InitChartInfo struct {
	Current string                `json:"current"`
	Charts  map[string]*InitChart `json:"charts"`
}

func (info *InitChartInfo) GetCurrentChart() *InitChart {
	return info.Charts[info.Current]
}
