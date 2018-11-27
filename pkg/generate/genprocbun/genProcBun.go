package genprocbun

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	PackageName    string `validate:"min=2,javaPackageName" arg:"required=true,shortname=p,longname=packagename" help:"Name of java package (friendly for java packages)."`
	TagName        string `validate:"min=2" arg:"shortname=t,longname=downloadtag,defaultvalue=${LatestSupportedTag} use version command for more info" help:"Tag name to pull the zip file github."`
	DestinationDir string `arg:"shortname=d,longname=destdir,defaultValue=./" help:"Destination directory for writing the runtime bundle template. This directory will be appended with the BundleName. Example: a destdir of '/Users/john/projects' and a bundlename 'my-bundle' will results in the runtime bundle being created in a final directory '/Users/john/projects/my-bundle'"`
}

//Install help deploys and verifies the full Activiti 7 Example application
func (l processBundleGenerator) GenerateRuntimeBundle(parms Parms) error {

	generator, has := generatorsByTag[parms.TagName]
	if !has {
		common.LogWarn(fmt.Sprintf("Unexpected tag name [%v], using Default tag [%v] as generation logic. This may or may not work. It has not been tested by this program.", parms.TagName, LatestSupportedTag))
		generator, has = generatorsByTag[LatestSupportedTag]
		if !has {
			common.LogExit(fmt.Sprintf("Fatal error, default tag %v is not registered in the system", LatestSupportedTag))
		}
	}
	common.LogWorking("TagName '" + parms.TagName + "' requested for download")
	return generator(parms)

}

type generateRuntimeBundler func(parms Parms) error

const sevenDot0Dot0DotBeta3 = "7.0.0.Beta3"

var generatorsByTag = map[string]generateRuntimeBundler{
	sevenDot0Dot0DotBeta3: type0GenerateRuntimeBundle,
}

const LatestSupportedTag = sevenDot0Dot0DotBeta3

func type0GenerateRuntimeBundle(parms Parms) error {
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}
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

	tmpRtBunZipFile := fmt.Sprintf("%vTmpRtBun.zip", tmpDir)

	rtBundleURL := fmt.Sprintf("https://github.com/Activiti/example-runtime-bundle/archive/%v.zip", parms.TagName)

	err = common.DownloadZipFromURL(rtBundleURL, tmpRtBunZipFile)
	if err != nil {
		common.LogError(err.Error())
		return err
	}

	outputZipDir, err := common.UnzipFromFileToDir(tmpRtBunZipFile, tmpDir)
	if err != nil {
		log.Fatal(err)
	}
	common.LogOK(fmt.Sprintf("Unzipped runtime bundle template to %v%v", tmpDir, outputZipDir))

	err = os.Remove(tmpRtBunZipFile)
	if err != nil {
		common.LogError(err.Error())
		return err
	}
	common.LogOK(fmt.Sprintf("Removed temp zip file %v", tmpRtBunZipFile))

	rtBundleTemplateTransCount := 9

	newOutputUnzipDir := tmpDir + parms.BundleName + "/"
	err = os.Rename(outputZipDir, newOutputUnzipDir)
	if err != nil {
		common.LogError(err.Error())
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 1 of %v: Rule rename template directory. Renamed from %v %v", rtBundleTemplateTransCount, outputZipDir, newOutputUnzipDir))

	theFile := newOutputUnzipDir + "pom.xml"
	theBytes, err := ioutil.ReadFile(theFile)
	theStr := string(theBytes)
	theStr = strings.Replace(theStr, "<groupId>org.activiti.cloud.examples</groupId>", fmt.Sprintf("<groupId>%v</groupId>", parms.PackageName), 1)
	theStr = strings.Replace(theStr, "<artifactId>example-runtime-bundle</artifactId>", fmt.Sprintf("<artifactId>%v</artifactId>", parms.BundleName), 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 2 of %v: Rule adjust group and artifact id in pom.xml. Adjusted with %v and %v with %v bytes at %v", rtBundleTemplateTransCount, parms.PackageName, parms.BundleName, len(theStr), theFile))

	theFile = newOutputUnzipDir + "skaffold.yaml"
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "{{.DOCKER_REGISTRY}}/activiti/example-runtime-bundle", fmt.Sprintf("{{.DOCKER_REGISTRY}}/activiti/%v", parms.BundleName), 3)
	theStr = strings.Replace(theStr, "changeme", parms.BundleName, 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 3 of %v: Rule 'adjust 4 uses in skaffold.yaml. Adjusted to %v with %v bytes at %v", rtBundleTemplateTransCount, parms.BundleName, len(theStr), theFile))

	theFile = newOutputUnzipDir + "Jenkinsfile"
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "example-runtime-bundle", parms.BundleName, 3)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 4 of %v: Rule adjust 3 uses in Jenkinsfile. Adjusted to %v with %v bytes at %v", rtBundleTemplateTransCount, parms.BundleName, len(theStr), theFile))

	renameFrom := newOutputUnzipDir + "charts/example-runtime-bundle"
	renameTo := newOutputUnzipDir + "charts/" + parms.BundleName
	err = os.Rename(renameFrom, renameTo)
	if err != nil {
		common.LogError(err.Error())
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 5 of %v: Rule rename the charts folder. Adjusted to %v", rtBundleTemplateTransCount, renameTo))

	theFile = newOutputUnzipDir + "src/main/resources/application.properties"
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "spring.application.name=${ACT_RB_APP_NAME:rb-my-app}", fmt.Sprintf("spring.application.name=${ACT_RB_APP_NAME:%v}", parms.BundleName), 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 6 of %v: Rule change the spring.application.name. Adjusted to %v with %v bytes at %v", rtBundleTemplateTransCount, parms.BundleName, len(theStr), theFile))

	theFile = newOutputUnzipDir + fmt.Sprintf("charts/%v/values.yaml", parms.BundleName)
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "name: rb-my-app", fmt.Sprintf("name: %v", parms.BundleName), 1)
	theStr = strings.Replace(theStr, "repository: draft", fmt.Sprintf("repository: %v", parms.BundleName), 1)
	theStr = strings.Replace(theStr, "tag: dev", "tag: latest", 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 7 of %v: Rule change helm service name and repository to bundle name and tag to 'latest'. Adjusted to %v with %v bytes at %v", rtBundleTemplateTransCount, parms.BundleName, len(theStr), theFile))

	theFile = newOutputUnzipDir + fmt.Sprintf("charts/%v/Chart.yaml", parms.BundleName)
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "example-runtime-bundle", parms.BundleName, 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 8 of %v: Rule change helm name. Adjusted to %v with %v bytes at %v", rtBundleTemplateTransCount, parms.BundleName, len(theStr), theFile))

	err = os.Rename(newOutputUnzipDir, finalOutDirectory)
	if err != nil {
		common.LogError(err.Error())
		return err
	}
	common.LogOK(fmt.Sprintf("RT Bundle Transform 9 of %v: Rule move the folder to the desired destination. Adjusted to %v", rtBundleTemplateTransCount, finalOutDirectory))

	common.LogInfo("Ready to use the Activiti Cloud Runtime Bundle")

	return nil
}
