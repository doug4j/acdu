package genprocbun

import (
	"fmt"

	common "github.com/doug4j/acdu/pkg/common"
)

//ProcessBundleGenerating creates an Activiti Cloud Runtime Bundle.
type ProcessBundleGenerating interface {
	GenerateRuntimeBundle(parms Parms) error
}

//NewProcessBundleGenerating create a new instance of ProcessBundleGenerating.
func NewProcessBundleGenerating() ProcessBundleGenerating {
	answer := processBundleGenerator{}
	return answer
}

type processBundleGenerator struct {
}

//Parms are the parameters for the command.
type Parms struct {
	BundleName     string `validate:"min=2,kubeFriendlyName" arg:"required=true,shortname=b,longname=bundlename" help:"Name of the runtime bundle (friendly for kubernetes and jars)."`
	PackageName    string `validate:"min=2,javaPackageName"  arg:"required=true,shortname=p,longname=packagename" help:"Name of java package (friendly for java packages)."`
	ProjectName    string `validate:"min=2" arg:"required=true,shortname=a,longname=projectname" help:"The default Project to use for the process bundle."`
	TagName        string `validate:"min=2" arg:"shortname=t,longname=downloadtag,defaultvalue=${LatestSupportedTag} use version command for more info" help:"Tag name to pull the zip file github."`
	DestinationDir string `arg:"shortname=d,longname=destdir,defaultValue=./" help:"Destination directory for writing the runtime bundle template. This directory will be appended with the BundleName. Example: a destdir of '/Users/john/projects' and a bundlename 'my-bundle' will results in the runtime bundle being created in a final directory '/Users/john/projects/my-bundle'"`
	Downloader     string `arg:"defaultValue=${DefaultDownloader}" help:"Downloader options (default not documented)"`
}

func (l processBundleGenerator) GenerateRuntimeBundle(parms Parms) error {
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}
	key := fmt.Sprintf(downloaderKeyFormat, parms.Downloader, parms.TagName)
	generator, has := generatorsByKey[key]
	if !has {
		common.LogWarn(fmt.Sprintf("Unexpected combination of downloader [%v] and tag [%v], using Default tag [%v] as generation logic. This may or may not work. It has not been tested by this program.", parms.Downloader, parms.TagName, LatestSupportedTag))
		key = fmt.Sprintf(downloaderKeyFormat, DefaultDownloader, LatestSupportedTag)
		generator, has = generatorsByKey[key]
		if !has {
			common.LogExit(fmt.Sprintf("Fatal error, default downloader [%v] and tag [%v] is not registered in the system", DefaultDownloader, LatestSupportedTag))
		}
		parms.Downloader = DefaultDownloader
		parms.TagName = LatestSupportedTag
	}
	common.LogWorking("TagName '" + parms.TagName + "' requested for download (using downloader=" + parms.Downloader + ")")
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

const downloaderKeyFormat = "downloader-%v-tag-%v"

const sevenDot0Dot0DotBeta3 = "7.0.0.Beta3"

type generateRuntimeBundler func(parms Parms) error

var generatorsByKey = map[string]generateRuntimeBundler{
	fmt.Sprintf(downloaderKeyFormat, common.DownloaderStereotypeQuickstart, sevenDot0Dot0DotBeta3): downloaderQuickstartType0,

	//This download implementor is probably only useful for beta3, but keeping the code in there as it's likely
	fmt.Sprintf(downloaderKeyFormat, common.DownloaderStereotypeExample, sevenDot0Dot0DotBeta3): downloaderExampleType0,
}

func downloaderQuickstartType0(parms Parms) error {
	genericParms := common.GenercQuickstartDownloadParms{
		BundleName:     parms.BundleName,
		DestinationDir: parms.DestinationDir,
		PackageName:    parms.PackageName,
		ProjectName:    parms.ProjectName,
		TagName:        parms.TagName,
		DownloadURL:    "https://github.com/Activiti/activiti-cloud-runtime-bundle-quickstart",
	}
	_, err := common.GenerateFromGenericQuickStartHandler(genericParms)
	return err
}

func downloaderExampleType0(parms Parms) error {
	genericParms := common.GenercExamplesDownloadParms{
		BundleName:     parms.BundleName,
		DestinationDir: parms.DestinationDir,
		PackageName:    parms.PackageName,
		TagName:        parms.TagName,
		DownloadURL:    "https://github.com/Activiti/example-runtime-bundle",
	}
	return common.GenerateFromGenericExampleWithHelmHandler(genericParms)
}
