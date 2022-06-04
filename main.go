package main

import (
	"blockraiders_tool/actions"
	"blockraiders_tool/cmds"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type helmArgs struct {
	HasBuild  bool
	HasDeploy bool
	Version   string
	Name      string
}

//goland:noinspection ALL
var globalUsage = ` The Blockraiders-tools

Common actions for Blockraiders-tools:
- blockraiders chart: for chart operator
`

//goland:noinspection ALL
func main() {
	cmd := &cobra.Command{
		Use:          "blockraiders",
		Short:        "The Blockraiders-tools",
		Long:         globalUsage,
		SilenceUsage: true,
	}

	actions.Client = &actions.ClientAction{
		CodegenAct: &actions.CodegenAction{},
		ChartAct:   &actions.ChartAction{},
	}

	cmd.AddCommand(cmds.NewCodegenCmd(), cmds.NewChartCmd(), cmds.NewInitCmd())
	if err := cmd.Execute(); err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}
