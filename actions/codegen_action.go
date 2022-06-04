package actions

type CodegenAction struct {
	Type          string
	UseQuery      bool
	UseCommand    bool
	Name          string
	Topic         string
	HasAuth       bool
	HasValidation bool
}
