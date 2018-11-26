package genmddoc

import (
	"github.com/doug4j/acdu/pkg/common"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

//InfrastructureInstalling applies core infrastructure supporting Activiti 7 into the given namespace.
type DocumentationGenerating interface {
	GenDoc(parms Parms, command cobra.Command) error
}

//NewDocumentationGenerating create a new instance of DocumentationGenerating.
func NewDocumentationGenerating() DocumentationGenerating {
	answer := documentationGenerater{}
	return answer
}

type documentationGenerater struct {
}

//Parms are the parameters for the command.
type Parms struct {
	DestinationDir string `validate:"" arg:"required=true,shortname=d,longname=destdir" help:"Destination directory for writing markdown documentation."`
}

//Install help deploys and verifies the full Activiti 7 Example application
func (l documentationGenerater) GenDoc(parms Parms, command cobra.Command) error {
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}
	//See https://github.com/spf13/cobra/blob/master/doc/md_docs.md
	return doc.GenMarkdownTree(&command, parms.DestinationDir)
}
