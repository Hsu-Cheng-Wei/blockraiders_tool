package actions

type ClientAction struct {
	ChartAct   *ChartAction
	CodegenAct *CodegenAction
}

var Client *ClientAction
