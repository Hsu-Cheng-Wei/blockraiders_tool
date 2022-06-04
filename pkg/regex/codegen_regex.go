package regex

import (
	"blockraiders_tool/actions"
	"blockraiders_tool/templates"
	"errors"
)

func Regex(cfg *templates.InitCodegen, args *actions.CodegenAction) error {
	if len(args.Type) == 0 {
		return errors.New("please enter the type [query|command]")
	}

	(&TemplateRegx{
		Args: args,
		Cfg:  cfg,
	}).Regex()

	return nil
}
