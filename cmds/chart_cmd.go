package cmds

import (
	"blockraiders_tool/actions"
	"blockraiders_tool/pkg/chart"
	"blockraiders_tool/pkg/config"
	"blockraiders_tool/pkg/docker"
	"blockraiders_tool/templates"
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/cmd/helm/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//goland:noinspection ALL
var chart_listHelp = `
This command include all of blockraiders chart template operator
blockraiders chart output [path] --name <chart-name>
blockraiders chart apply [api/host] [version] --name <chart-name> --config <config-path> --path <ingress-path> --count <instance-count> --build <docker-file>
blockraiders chart update [version] --build <docker-file>
blockraiders chart delete --name <chart-name>
`

type chartOp struct {
	Cfg      *templates.ConfigDto
	Command  *cobra.Command
	Args     []*string
	TypeName *string
}

func NewChartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chart [Action]",
		Short:   "chart operator",
		Long:    chart_listHelp,
		Aliases: []string{"ch"},
		Args:    require.MinimumNArgs(1),
		RunE:    doChartAct,
	}

	addChartFlags(cmd.Flags())
	return cmd
}

func addChartFlags(flag *pflag.FlagSet) {
	flag.StringVarP(&actions.Client.ChartAct.Name, "name", "n", "", "output chart name")
	flag.StringVarP(&actions.Client.ChartAct.ReInstall, "re-install", "r", "", "reinstall")
	flag.StringVar(&actions.Client.ChartAct.ConfigurePath, "config", "", "choose config path")
	flag.StringVarP(&actions.Client.ChartAct.IngressPath, "path", "p", "", "choose ingress path")
	flag.StringVarP(&actions.Client.ChartAct.DockerImage, "build", "b", "", "build and push")
	flag.IntVar(&actions.Client.ChartAct.Count, "count", 1, "instance count")
}

func doChartAct(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Usage()
	}

	action := args[0]

	if err := config.EnsureConfig(); err != nil {
		return err
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	var argsp []*string
	for i := range args {
		argsp = append(argsp, &args[i])
	}

	op := &chartOp{
		Cfg:     cfg,
		Command: cmd,
		Args:    argsp,
	}

	switch action {
	case "output":
		return op.doOutPutChart()
	case "update":
		return op.doUpdateChart()
	case "delete":
		return op.doUninstallChart()
	case "apply":
		return op.doApplyChart()
	}

	return nil
}

func (op *chartOp) doOutPutChart() error {
	action := *op.Args[0]
	if action != "output" {
		return errors.New("parameter is not output")
	}

	path := "."
	if len(op.Args) == 2 {
		path = *op.Args[1]
	}

	ch, err := chart.NewChartOperator[any]("")
	if err != nil {
		return err
	}
	if len(actions.Client.ChartAct.Name) != 0 {
		return ch.OutputChartByName(actions.Client.ChartAct.Name, path)
	}

	index, err := chooseChart(ch)

	if err != nil {
		return err
	}

	return ch.OutputChartByIndex(index, path)
}

func (op *chartOp) doUpdateChart() error {
	action := *op.Args[0]
	if action != "update" {
		return errors.New("parameter is not update")
	}
	if len(op.Args) == 1 {
		return errors.New("parameter version doesn't exist")
	}
	version := *op.Args[1]

	ch, err := chart.NewChartOperator[any]("")
	if err != nil {
		return err
	}

	if len(actions.Client.ChartAct.Name) != 0 {
		return ch.UpdateVersionByName(actions.Client.ChartAct.Name, version)
	}

	index, err := chooseChart(ch)
	if err != nil {
		return err
	}

	op.TypeName = &actions.Client.ChartAct.ReInstall
	if len(actions.Client.ChartAct.ReInstall) != 0 {
		switch *op.TypeName {
		case "api":
			err = ch.UninstallByIndex(index)
			if err != nil {
				return err
			}
			name, err := op.genName(version, "api")
			if err != nil {
				return err
			}
			err = op.applyApi(name, version)
			if err != nil {
				return err
			}
			break
		case "host":
			err = ch.UninstallByIndex(index)
			if err != nil {
				return err
			}
			name, err := op.genName(version, "host")
			if err != nil {
				return err
			}
			err = op.applyHost(name, version)
			if err != nil {
				return err
			}
			break
		default:
			return fmt.Errorf("doesn't include action: %s", actions.Client.ChartAct.ReInstall)
		}

	}

	return ch.UpdateVersionByIndex(index, version)
}

func (op *chartOp) doUninstallChart() error {
	action := *op.Args[0]
	if action != "delete" {
		return errors.New("parameter is not delete")
	}

	ch, err := chart.NewChartOperator[any]("")
	if err != nil {
		return err
	}

	if len(actions.Client.ChartAct.Name) != 0 {
		return ch.UninstallByName(actions.Client.ChartAct.Name)
	}

	index, err := chooseChart(ch)

	if err != nil {
		return err
	}

	return ch.UninstallByIndex(index)
}

//goland:noinspection SpellCheckingInspection
func (op *chartOp) doApplyChart() error {
	var err error
	action := *op.Args[0]
	if action != "apply" {
		return errors.New("parameter is not apply")
	}
	if len(op.Args) != 3 {
		return errors.New("parameter invalid")
	}

	if err = config.EnsureConfig(); err != nil {
		return err
	}

	op.TypeName = op.Args[1]
	version := *op.Args[2]

	var name string

	if len(actions.Client.ChartAct.Name) > 0 {
		name = actions.Client.ChartAct.Name
	} else {
		name, err = op.genName(version, *op.TypeName)
		if err != nil {
			return err
		}
	}

	switch *op.TypeName {
	case "api":
		return op.applyApi(name, version)
	case "host":
		return op.applyHost(name, version)
	default:
		return fmt.Errorf("doesn't insluce type name %s", *op.TypeName)
	}
}

func (op *chartOp) genName(version, typeName string) (string, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return "", err
	}
	return cfg.ChartInfo.GetCurrentChart().Name + "-" + typeName + versionSplitRemove(version), nil
}

func (op *chartOp) applyApi(name string, version string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = op.build(wd, version)
	if err != nil {
		return err
	}

	ch, err := chart.NewChartOperator[templates.ApiValuesDto]("")
	if err != nil {
		return err
	}

	path, err := os.Executable()
	if err != nil {
		return err
	}

	var ingressPath string
	if len(actions.Client.ChartAct.IngressPath) == 0 {
		ingressPath = "/" + version
	} else {
		ingressPath = actions.Client.ChartAct.IngressPath
	}

	path = filepath.Dir(path)
	path = filepath.Join(path, op.Cfg.ChartInfo.GetCurrentChart().ApiTemplate)

	app, log4, err := GetConfig()
	if err != nil {
		return err
	}

	fmt.Println("start deploy chart...")
	return ch.Install(path, name, version, func(values *templates.ApiValuesDto) {
		values.Image.Repository = op.Cfg.ChartInfo.GetCurrentChart().ApiRepository
		values.Body.Namespace = op.Cfg.ChartInfo.GetCurrentChart().Namespace
		values.Body.Name = name
		values.Config.Data.Appsettings = app
		values.Config.Data.Log4netConfig = log4
		values.Config.ConfigMapName = name + "-config"
		values.Config.ConfigName = "config"
		values.Count = actions.Client.ChartAct.Count

		for _, mount := range values.VolumeMounts {
			mount.Name = "config"
		}

		for _, host := range values.Ingress.Hosts {
			host.Path = ingressPath
		}

		for _, m := range values.VolumeMounts {
			m.Name = values.Config.ConfigName
		}
	})
}

func (op *chartOp) applyHost(name string, version string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = op.build(wd, version)
	if err != nil {
		return err
	}

	ch, err := chart.NewChartOperator[templates.GameHostValuesDto]("")
	if err != nil {
		return err
	}

	path, err := os.Executable()
	if err != nil {
		return err
	}

	path = filepath.Dir(path)
	path = filepath.Join(path, op.Cfg.ChartInfo.GetCurrentChart().GameHostTemplate)

	app, log4, err := GetConfig()
	if err != nil {
		return err
	}

	fmt.Println("start deploy chart...")
	return ch.Install(path, name, version, func(values *templates.GameHostValuesDto) {
		values.Image.Repository = op.Cfg.ChartInfo.GetCurrentChart().GameHostRepository
		values.Body.Namespace = op.Cfg.ChartInfo.GetCurrentChart().Namespace
		values.Body.Name = name
		values.Config.Data.Appsettings = app
		values.Config.Data.Log4netConfig = log4
		values.Config.ConfigMapName = name + "-config"
		values.Config.ConfigName = "config"

		for _, mount := range values.VolumeMounts {
			mount.Name = "config"
		}

		for _, m := range values.VolumeMounts {
			m.Name = values.Config.ConfigName
		}
	})
}

func (op *chartOp) build(path string, version string) error {
	if len(actions.Client.ChartAct.DockerImage) == 0 {
		return nil
	}
	host, err := docker.NewHost()
	if err != nil {
		return err
	}

	var tag string
	if *op.TypeName == "api" {
		tag = op.Cfg.ChartInfo.GetCurrentChart().ApiRepository
	} else {
		tag = op.Cfg.ChartInfo.GetCurrentChart().GameHostRepository
	}
	tag = tag + ":" + version

	dockerFilePath := filepath.Join(path, actions.Client.ChartAct.DockerImage)
	fmt.Printf("Build DockerFile:%s", dockerFilePath)
	if err = host.Build(dockerFilePath, tag, nil); err != nil {
		return err
	}

	token, err := host.Login(op.Cfg.ChartInfo.GetCurrentChart().Region)
	if err != nil {
		return err
	}
	if err = host.Push(tag, token); err != nil {
		return err
	}
	return nil
}

func versionSplitRemove(version string) string {
	res := ""
	for _, s := range strings.Split(version, ".") {
		res += s
	}
	return res
}

func chooseChart[T any](ch *chart.Operator[T]) (int, error) {
	err := ch.GetAndPrintRelease()

	if err != nil {
		return -1, err
	}

	fmt.Print("Select:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return strconv.Atoi(scanner.Text())
}

//goland:noinspection ALL
func GetConfig() (string, string, error) {
	if len(actions.Client.ChartAct.ConfigurePath) == 0 {
		return "", "", errors.New("choose config path")
	}
	m := make(map[string]*fileInfo)
	err := recursiveFilePath(m, actions.Client.ChartAct.ConfigurePath)

	if err != nil {
		return "", "", err
	}

	app, exist := m["appsettings.json"]
	if !exist {
		return "", "", errors.New("config doesn't include appsettings.jaon")
	}

	log4net, exist := m["log4net.config"]
	if !exist {
		return "", "", errors.New("config doesn't include log4net.config")
	}

	appContent, err := ioutil.ReadFile(app.Path)
	if !exist {
		return "", "", err
	}
	log4netContent, err := ioutil.ReadFile(log4net.Path)
	if !exist {
		return "", "", err
	}

	return string(appContent), string(log4netContent), nil
}

func recursiveFilePath(rel map[string]*fileInfo, dir string) error {
	files, err := ioutil.ReadDir(actions.Client.ChartAct.ConfigurePath)

	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if err = recursiveFilePath(rel, filepath.Join(dir, f.Name())); err != nil {
				return err
			}
			continue
		}
		rel[f.Name()] = &fileInfo{
			Name: f.Name(),
			Path: filepath.Join(dir, f.Name()),
		}
	}
	return nil
}

type fileInfo struct {
	Name string
	Path string
}
