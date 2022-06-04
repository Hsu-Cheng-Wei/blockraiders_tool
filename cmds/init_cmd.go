package cmds

import (
	"blockraiders_tool/pkg/config"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/cmd/helm/require"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "init",
		Long:    chart_listHelp,
		Aliases: []string{},
		Args:    require.NoArgs,
		RunE:    doInitAct,
	}

	return cmd
}

func doInitAct(_ *cobra.Command, _ []string) error {
	if _, err := config.GetAndCreateConfigPath(); err != nil {
		return err
	}

	if _, err := config.CreateConfig(); err != nil {
		return err
	}

	return nil
}
