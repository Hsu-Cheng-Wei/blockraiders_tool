package cmds

import (
	"blockraiders_tool/actions"
	"blockraiders_tool/pkg/config"
	"blockraiders_tool/pkg/regex"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/cmd/helm/require"
)

//goland:noinspection ALL
var codegen_listHelp = `
This command include all of blockraiders chart template operator
blockraiders codegen [Option]
`

//var client actions.ChartAction

func NewCodegenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "codegen [option]",
		Short: "code auto generator",
		Long:  codegen_listHelp,
		Args:  require.NoArgs,
		RunE:  doCodegenAct,
	}

	addCodegenFlags(cmd.Flags())
	return cmd
}

func doCodegenAct(_ *cobra.Command, _ []string) error {

	if err := config.EnsureConfig(); err != nil {
		return err
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	args := actions.Client.CodegenAct
	if args.UseQuery {
		args.Type = "query"
	}

	if args.UseCommand {
		args.Type = "command"
	}

	return regex.Regex(cfg.CodegenInfo, args)
}

func addCodegenFlags(flag *pflag.FlagSet) {
	flag.BoolVarP(&actions.Client.CodegenAct.UseQuery, "query", "q", false, "query template")
	flag.BoolVarP(&actions.Client.CodegenAct.UseCommand, "command", "c", false, "command template")
	flag.StringVarP(&actions.Client.CodegenAct.Name, "name", "n", "", "template name")
	flag.StringVarP(&actions.Client.CodegenAct.Topic, "topic", "t", "", "template topic")
	flag.BoolVarP(&actions.Client.CodegenAct.HasValidation, "validate", "v", false, "[Option] generator validate template")
}
