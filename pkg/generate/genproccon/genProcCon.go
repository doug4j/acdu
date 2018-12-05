package genproccon

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	common "github.com/doug4j/acdu/pkg/common"
)

//ProcessConnectorGenerating creates an Activiti Cloud Runtime Bundle.
type ProcessConnectorGenerating interface {
	GenerateConnector(parms Parms) error
}

//NewProcessConnectorGenerating create a new instance of ProcessConnectorGenerating.
func NewProcessConnectorGenerating() ProcessConnectorGenerating {
	answer := processConnectorGenerator{}
	return answer
}

type processConnectorGenerator struct {
}

//Parms are the parameters for the command.
type Parms struct {
	BundleName         string `validate:"min=2,max=255,kubeFriendlyName" arg:"required=true,shortname=b,longname=bundle" help:"Name of the runtime bundle (friendly for kubernetes and jars)."`
	PackageName        string `validate:"min=2,max=255,javaPackageName" arg:"required=true,shortname=p,longname=package" help:"Name of java package (friendly for java packages)."`
	ProjectName        string `validate:"min=2" arg:"required=true,shortname=a,longname=projectname" help:"The default Project to use for the process connector."`
	TagName            string `validate:"min=2" arg:"shortname=t,longname=downloadtag,defaultvalue=${LatestSupportedTag} use version command for more info" help:"Tag name to pull the zip file github."`
	DestinationDir     string `arg:"shortname=d,longname=destdir,defaultValue=./" help:"Destination directory for writing the runtime bundle template. This directory will be appended with the BundleName. Example: a destdir of '/Users/john/projects' and a bundlename 'my-bundle' will results in the runtime bundle being created in a final directory '/Users/john/projects/my-bundle'"`
	ChannelName        string `validate:"min=2,max=255,alphanum,startsWithLowerCaseAsciiAlpha" arg:"required=true,shortname=c,longname=channel" help:"Name of implementation (starting lower case alpha and all alphanum)."`
	ImplementationName string `validate:"min=2,max=255,alphanum,startsWithUpperCaseAsciiAlpha" arg:"required=true,shortname=i,longname=implementation" help:"Name of implementation (starting lower case alpha and all alphanum)."`
}

//Install help deploys and verifies the full Activiti 7 Example application
func (l processConnectorGenerator) GenerateConnector(parms Parms) error {
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}
	generator, has := generatorsByTag[parms.TagName]
	if !has {
		common.LogWarn(fmt.Sprintf("Unexpected tag name [%v], using Default tag [%v] as generation logic. This may or may not work. It has not been tested by this program.", parms.TagName, LatestSupportedTag))
		generator, has = generatorsByTag[LatestSupportedTag]
		if !has {
			common.LogExit(fmt.Sprintf("Fatal error, default tag %v is not registered in the system", LatestSupportedTag))
		}
		parms.TagName = LatestSupportedTag
	}
	common.LogWorking("TagName '" + parms.TagName + "' requested for download")
	if err := generator(parms); err != nil {
		return err
	}
	common.LogInfo("Ready to use the Activiti Cloud Connector")
	return nil
}

func ImplementationsString() string {
	var answer string
	for name := range generatorsByTag {
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

const sevenDot0Dot0DotBeta3 = "7.0.0.Beta3"

type generateRuntimeBundler func(parms Parms) error

var generatorsByTag = map[string]generateRuntimeBundler{
	sevenDot0Dot0DotBeta3: downloaderQuickstartType0,
}

func downloaderQuickstartType0(parms Parms) error {
	genericParms := common.GenercQuickstartDownloadParms{
		BundleName:           parms.BundleName,
		DestinationDir:       parms.DestinationDir,
		PackageName:          parms.PackageName,
		ProjectName:          parms.ProjectName,
		TagName:              parms.TagName,
		DownloadURL:          "https://github.com/Activiti/activiti-cloud-connector-quickstart",
		AdditionalTransCount: 2,
	}
	localTransformsCount, err := common.GenerateFromGenericQuickStartHandler(genericParms)
	if err != nil {
		return err
	}
	totalTransformsCount := genericParms.AdditionalTransCount + localTransformsCount
	destDir := common.EnsureFowardSlashAtStringEnd(parms.DestinationDir)
	finalOutDirectory := destDir + parms.BundleName

	theFile := finalOutDirectory + "/src/main/java/org/activiti/cloud/connector/impl/ExampleConnectorChannels.java"
	//theFile := finalOutDirectory + fmt.Sprintf("/%v/src/main/java/org/activiti/cloud/connector/impl/ExampleConnectorChannels.java", parms.BundleName)
	theBytes, err := ioutil.ReadFile(theFile)
	theStr := string(theBytes)
	//Note: the following code Replaces logically as follows
	//String EXAMPLE_CONNECTOR_CONSUMER = "${parms.ChannelName}";
	//SubscribableChannel ${parms.ChannelName}();
	theStr = strings.Replace(theStr, "exampleConnectorConsumer", parms.ChannelName, 2)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("Transform %v of %v: Rule adjust chanel in ExampleConnectorChannels.java. Adjusted as %v with %v bytes at %v", localTransformsCount+1, totalTransformsCount, parms.ChannelName, len(theStr), theFile))

	theFile = finalOutDirectory + "/src/main/resources/application.properties"
	//theFile := finalOutDirectory + fmt.Sprintf("/%v/src/main/java/org/activiti/cloud/connector/impl/ExampleConnectorChannels.java", parms.BundleName)
	theBytes, err = ioutil.ReadFile(theFile)
	if err != nil {
		return err
	}
	theStr = string(theBytes)
	//Replaces
	// spring.cloud.stream.bindings.${parms.ChannelName}.destination=${parms.ImplementationName}
	// spring.cloud.stream.bindings.${parms.ChannelName}.contentType=application/json
	// spring.cloud.stream.bindings.${parms.ChannelName}.group=${spring.application.name}
	theStr = strings.Replace(theStr, "exampleConnectorConsumer", parms.ChannelName, 3)
	theStr = strings.Replace(theStr, "<<@TODO: ADD YOUR CONNECTOR NAME HERE>>", parms.ImplementationName, 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("Transform %v of %v: Rule application.properties file for channel name and . Adjusted as %v with %v bytes at %v", localTransformsCount+1, totalTransformsCount, parms.ChannelName, len(theStr), theFile))

	return nil
}
