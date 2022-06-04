package regex

import (
	"blockraiders_tool/actions"
	"blockraiders_tool/templates"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateRegx struct {
	Args *actions.CodegenAction
	Cfg  *templates.InitCodegen
}

func (regex *TemplateRegx) Regex() {
	var prefix string
	if regex.Args.UseQuery {
		prefix = regex.Cfg.QueryPrefix
	} else {
		prefix = regex.Cfg.CommandPrefix
	}

	appPath := regex.getBasePath(prefix)

	err := os.MkdirAll(appPath, 0755)

	if err != nil {
		panic(err)
	}

	regex.genArgsFile(prefix)

	regex.genHandlerFile(prefix)

}

func (regex *TemplateRegx) getApplicationPaths(prefix string) []string {
	return []string{regex.Cfg.ApplicationPath, regex.Args.Topic, prefix, regex.Args.Name}
}

func (regex *TemplateRegx) getBasePath(prefix string) string {
	path, _ := os.Getwd()

	for _, p := range regex.getApplicationPaths(prefix) {
		path = filepath.Join(path, p)
	}

	return path
}

func (regex *TemplateRegx) getFilePath(prefix string) string {
	appPath := regex.getBasePath(prefix)

	return filepath.Join(appPath, regex.Args.Name+prefix+".cs")
}

func (regex *TemplateRegx) genArgsFile(prefix string) bool {
	return regex.genQueryFileFromTemplate(templates.ArgsTemplate, regex.getFilePath(prefix), map[string]interface{}{
		"namespace": NamespaceRegex(regex.Cfg.Namespace, regex.getApplicationPaths(prefix)),
		"type":      regex.Args.Name + prefix,
	})
}

func (regex *TemplateRegx) getHandlerFilePath(prefix string) string {
	appPath := regex.getBasePath(prefix)

	return filepath.Join(appPath, regex.Args.Name+prefix+"Handler.cs")
}

func (regex *TemplateRegx) genHandlerFile(prefix string) bool {
	return regex.genQueryFileFromTemplate(templates.HandlerTemplate, regex.getHandlerFilePath(prefix), map[string]interface{}{
		"namespace": NamespaceRegex(regex.Cfg.Namespace, regex.getApplicationPaths(prefix)),
		"type":      regex.Args.Name + prefix,
		"typeName":  strings.ToLower(prefix),
		"handler":   regex.Args.Name + prefix + "Handler",
	})
}

func (regex *TemplateRegx) genQueryFileFromTemplate(content string, filePath string, args map[string]interface{}) bool {
	tmp, err := template.New("query").Parse(content)

	if _, err := os.Stat(filePath); err == nil {
		fmt.Print("File:" + filePath + " exist do you want to overwrite it?[yes/no]:")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := scanner.Text()

		if strings.ToLower(answer) == "yes" {
			_ = os.Remove(filePath)
		} else {
			return false
		}
	}

	f, err := os.Create(filePath)

	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		panic(err)
	}

	err = tmp.Execute(f, args)

	if err != nil {
		fmt.Println(err)
	}
	return true
}
