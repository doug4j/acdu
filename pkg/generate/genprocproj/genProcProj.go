package genprocproj

import (
	"errors"
	"fmt"
	"os"

	"github.com/doug4j/acdu/pkg/generate/genproccon"

	"github.com/doug4j/acdu/pkg/generate/genprocbun"

	common "github.com/doug4j/acdu/pkg/common"
)

//ProcessProjectGenerating creates an Activiti Cloud Runtime Bundle and Cloud Connector.
type ProcessProjectGenerating interface {
	GenerateProcessProject(parms Parms) error
}

//NewProcessProjectGenerating create a new instance of ProcessBundleGenerating.
func NewProcessProjectGenerating() ProcessProjectGenerating {
	answer := processProjectGenerator{}
	return answer
}

type processProjectGenerator struct {
}

//Parms are the parameters for the command.
type Parms struct {
	BundleName         string `validate:"min=2,kubeFriendlyName" arg:"required=true,shortname=b,longname=bundlename" help:"Name of the runtime bundle (friendly for kubernetes and jars)."`
	PackageName        string `validate:"min=2,javaPackageName"  arg:"required=true,shortname=p,longname=packagename" help:"Name of java package (friendly for java packages)."`
	ProjectName        string `validate:"min=2" arg:"required=true,shortname=a,longname=projectname" help:"The default Project to use for the process bundle."`
	TagName            string `validate:"min=2" arg:"shortname=t,longname=downloadtag,defaultvalue=${LatestSupportedTag} use version command for more info" help:"Tag name to pull the zip file github."`
	DestinationDir     string `arg:"shortname=d,longname=destdir,defaultValue=./" help:"Destination directory for writing the runtime bundle template. This directory will be appended with the BundleName. Example: a destdir of '/Users/john/projects' and a bundlename 'my-bundle' will results in the runtime bundle being created in a final directory '/Users/john/projects/my-bundle'"`
	ChannelName        string `validate:"min=2,max=255,alphanum,startsWithLowerCaseAsciiAlpha" arg:"required=true,shortname=c,longname=channel" help:"Name of implementation (starting lower case alpha and all alphanum)."`
	ImplementationName string `validate:"min=2,max=255,alphanum,startsWithUpperCaseAsciiAlpha" arg:"required=true,shortname=i,longname=implementation" help:"Name of implementation (starting lower case alpha and all alphanum)."`
}

func (l processProjectGenerator) GenerateProcessProject(parms Parms) error {
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}
	generator, has := generatorsByKey[parms.TagName]
	if !has {
		common.LogWarn(fmt.Sprintf("Unexpected tag [%v], using Default tag [%v] as generation logic. This may or may not work. It has not been tested by this program.", parms.TagName, LatestSupportedTag))
		generator, has = generatorsByKey[parms.TagName]
		if !has {
			common.LogExit(fmt.Sprintf("Fatal error, default downloader [%v] and tag [%v] is not registered in the system", DefaultDownloader, LatestSupportedTag))
		}
		parms.TagName = LatestSupportedTag
	}
	common.LogWorking("TagName '" + parms.TagName + "' requested for downloads")
	if err := generator(parms); err != nil {
		return err
	}
	common.LogInfo("Ready to use the Activiti Cloud Runtime Bundle")
	return nil
}

func ImplementationsString() string {
	var answer string
	for name := range generatorsByKey {
		if answer == "" {
			answer = name
		} else {
			answer = answer + ", " + name
		}
	}
	return answer
}

//Defaults

const LatestSupportedTag = sevenDot0Dot0DotBeta3
const DefaultDownloader = common.DownloaderStereotypeQuickstart

const sevenDot0Dot0DotBeta3 = "7.0.0.Beta3"

type generateRuntimeBundler func(parms Parms) error

var generatorsByKey = map[string]generateRuntimeBundler{
	sevenDot0Dot0DotBeta3: downloaderQuickstartType0,
}

func downloaderQuickstartType0(parms Parms) error {
	//Note: we could easily put this in a goroutine, but for simplicity of following logs and the fact that the
	//user won't practically notice the difference with the size of files and processing we're talking about below.

	destDir := common.EnsureFowardSlashAtStringEnd(parms.DestinationDir) + parms.BundleName + "-processbundle"
	if _, err := os.Stat(destDir); os.IsExist(err) {
		msg := fmt.Sprintf("Output directory '%v' already exists and frther processing might result in loss of data, double check your configuration and try again", destDir)
		return errors.New(msg)
	}
	err := os.Mkdir(destDir, os.ModePerm)
	if err != nil {
		return err
	} else {
		common.LogOK(fmt.Sprintf("Created process bundle directory %v", destDir))
	}

	bundleGenerator := genprocbun.NewProcessBundleGenerating()
	err = bundleGenerator.GenerateRuntimeBundle(genprocbun.Parms{
		BundleName:     parms.BundleName + "-processbundle",
		PackageName:    parms.PackageName,
		ProjectName:    parms.ProjectName,
		TagName:        parms.TagName,
		DestinationDir: destDir,
		Downloader:     genprocbun.DefaultDownloader,
	})
	if err != nil {
		return err
	}
	destDir = common.EnsureFowardSlashAtStringEnd(parms.DestinationDir) + parms.BundleName + "-connector"
	if _, err := os.Stat(destDir); os.IsExist(err) {
		msg := fmt.Sprintf("Output directory '%v' already exists and frther processing might result in loss of data, double check your configuration and try again", destDir)
		return errors.New(msg)
	}
	err = os.Mkdir(destDir, os.ModePerm)
	if err != nil {
		return err
	} else {
		common.LogOK(fmt.Sprintf("Created process bundle directory %v", destDir))
	}
	connectorGenerator := genproccon.NewProcessConnectorGenerating()
	err = connectorGenerator.GenerateConnector(genproccon.Parms{
		BundleName:         parms.BundleName + "-connector",
		PackageName:        parms.PackageName,
		ProjectName:        parms.ProjectName,
		TagName:            parms.TagName,
		DestinationDir:     common.EnsureFowardSlashAtStringEnd(parms.DestinationDir) + parms.BundleName + "-connector",
		ChannelName:        parms.ChannelName,
		ImplementationName: parms.ImplementationName,
	})
	if err != nil {
		return err
	}
	common.LogInfo("Done creating project")
	return nil
}
