package genproccon

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	ChannelName        string `validate:"min=2,max=255,alphanum,startsWithLowerCaseAsciiAlpha" arg:"required=true,shortname=c,longname=channel" help:"Name of implementation (starting lower case alpha and all alphanum)."`
	ImplementationName string `validate:"min=2,max=255,alphanum,startsWithUpperCaseAsciiAlpha" arg:"required=true,shortname=i,longname=implementation" help:"Name of implementation (starting lower case alpha and all alphanum)."`
	TagName            string `validate:"min=2" arg:"shortname=t,longname=downloadtag,defaultvalue=${LatestSupportedTag} use version command for more info" help:"Tag name to pull the zip file github."`
	DestinationDir     string `arg:"shortname=d,longname=destdir,defaultValue=./" help:"Destination directory for writing the runtime bundle template. This directory will be appended with the BundleName. Example: a destdir of '/Users/john/projects' and a bundlename 'my-bundle' will results in the runtime bundle being created in a final directory '/Users/john/projects/my-bundle'"`
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
	}
	common.LogWorking("TagName '" + parms.TagName + "' requested for download")

	// DownloadURL:         fmt.Sprintf("https://github.com/Activiti/example-cloud-connector/archive/%v.zip", parms.TagName),
	// SourcePackageName:   "org.activiti.cloud.examples",
	// SourceBundleName:    "example-cloud-connector",
	// SourceSpringAppName: "example-connector",

	if err := generator(parms); err != nil {
		return err
	}
	common.LogInfo("Ready to use the Activiti Cloud Connector")
	return nil
}

type generateRuntimeBundler func(parms Parms) error

const sevenDot0Dot0DotBeta3 = "7.0.0.Beta3"

var generatorsByTag = map[string]generateRuntimeBundler{
	sevenDot0Dot0DotBeta3: type0GenerateRuntimeBundle,
}

const LatestSupportedTag = sevenDot0Dot0DotBeta3

func type0GenerateRuntimeBundle(parms Parms) error {
	var err error

	destDir := common.EnsureFowardSlashAtStringEnd(parms.DestinationDir)

	finalOutDirectory := destDir + parms.BundleName

	finalOutDirectoryExists := true
	if _, err := os.Stat(finalOutDirectory); os.IsNotExist(err) {
		finalOutDirectoryExists = false
	}
	if finalOutDirectoryExists {
		msg := fmt.Sprintf("Final output directory is calculated as '%v' which might result in loss of data, double check your configuration and try again", finalOutDirectory)
		common.LogError(msg)
		return errors.New(msg)
	}

	tmpDir := destDir + ".acdu-tmp/"

	err = os.RemoveAll(tmpDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		//Don't report on cleanup of directory that didn't exist
	} else {
		common.LogInfo(fmt.Sprintf("Cleaned up previous temp director %v", tmpDir))
	}

	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = os.Mkdir(tmpDir, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		common.LogOK(fmt.Sprintf("Temp directory %v already exists", tmpDir))
	} else {
		common.LogOK(fmt.Sprintf("Created temp directory %v", tmpDir))
	}
	defer func() {
		if !common.VerboseLogging {
			err = os.RemoveAll(tmpDir)
			if err != nil {
				common.LogError(fmt.Sprintf("Could not remove temp directory from %v:%v", tmpDir, err))
			}
			common.LogOK(fmt.Sprintf("Removed temp directory from %v", tmpDir))
		} else {
			common.LogInfo(fmt.Sprintf("Temp directory %v was not deleted as verbose logging is turned on", tmpDir))
		}
	}()

	tmpZipFile := fmt.Sprintf("%v-acdu-Tmp.zip", tmpDir)
	downloadURL := fmt.Sprintf("https://github.com/Activiti/example-cloud-connector/archive/%v.zip", parms.TagName)

	err = common.DownloadZipFromURL(downloadURL, tmpZipFile)
	if err != nil {
		common.LogError(err.Error())
		return err
	}

	outputZipDir, err := common.UnzipFromFileToDir(tmpZipFile, tmpDir)
	if err != nil {
		log.Fatal(err)
	}
	common.LogOK(fmt.Sprintf("Unzipped runtime bundle template to %v%v", tmpDir, outputZipDir))

	err = os.Remove(tmpZipFile)
	if err != nil {
		common.LogError(err.Error())
		return err
	}
	common.LogOK(fmt.Sprintf("Removed temp zip file %v", tmpZipFile))

	totalTransformsCount := 4

	newOutputUnzipDir := tmpDir + parms.BundleName + "/"
	err = os.Rename(outputZipDir, newOutputUnzipDir)
	if err != nil {
		common.LogError(err.Error())
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 1 of %v: Rule rename template directory. Renamed from %v %v", totalTransformsCount, outputZipDir, newOutputUnzipDir))

	theFile := newOutputUnzipDir + "pom.xml"
	theBytes, err := ioutil.ReadFile(theFile)
	theStr := string(theBytes)
	theStr = strings.Replace(theStr, "<groupId>org.activiti.cloud.examples</groupId>", fmt.Sprintf("<groupId>%v</groupId>", parms.PackageName), 1)
	theStr = strings.Replace(theStr, "<artifactId>example-cloud-connector</artifactId>", fmt.Sprintf("<artifactId>%v</artifactId>", parms.BundleName), 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 2 of %v: Rule adjust group and artifact id in pom.xml. Adjusted with %v and %v with %v bytes at %v", totalTransformsCount, parms.PackageName, parms.BundleName, len(theStr), theFile))

	theFile = newOutputUnzipDir + "src/main/resources/application.properties"
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "spring.application.name=${ACT_RB_APP_NAME:rb-my-app}", fmt.Sprintf("spring.application.name=${ACT_RB_APP_NAME:%v}", parms.BundleName), 1)
	theStr = strings.Replace(theStr, "spring.cloud.stream.bindings.exampleConnectorConsumer", fmt.Sprintf("spring.cloud.stream.bindings.%v", parms.ChannelName), 3)
	theStr = strings.Replace(theStr, "spring.cloud.stream.bindings.exampleConnectorConsumer.destination=Example Connector", fmt.Sprintf("spring.cloud.stream.bindings.exampleConnectorConsumer.destination=%v", parms.ImplementationName), 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 3 of %v: Rule change the spring.application.name, ...stream.bindings.*, and destination. Adjusted to %v with %v bytes at %v", totalTransformsCount, parms.BundleName, len(theStr), theFile))

	err = os.Rename(newOutputUnzipDir, finalOutDirectory)
	if err != nil {
		common.LogError(err.Error())
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 4 of %v: Rule move the folder to the desired destination. Adjusted to %v", totalTransformsCount, finalOutDirectory))
	return nil
}
