package common

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

const DownloaderStereotypeQuickstart = "stereotype-github-quickstart"
const DownloaderStereotypeExample = "stereotype-github-example"

type GenercExamplesDownloadParms struct {
	BundleName     string
	PackageName    string
	TagName        string
	DestinationDir string
	DownloadURL    string
}

//GenerateFromGenericExampleWithHelmHandler downloads and processes a template application from https://github.com/Activiti/example-runtime-bundle
func GenerateFromGenericExampleWithHelmHandler(parms GenercExamplesDownloadParms) error {
	var err error

	destDir := EnsureFowardSlashAtStringEnd(parms.DestinationDir)

	finalOutDirectory := destDir + parms.BundleName

	finalOutDirectoryExists := true
	if _, err := os.Stat(finalOutDirectory); os.IsNotExist(err) {
		finalOutDirectoryExists = false
	}
	if finalOutDirectoryExists {
		msg := fmt.Sprintf("Final output directory is calculated as '%v' which might result in loss of data, double check your configuration and try again", finalOutDirectory)
		LogError(msg)
		return errors.New(msg)
	}

	tmpDir := destDir + fmt.Sprintf(".acdu-tmp-pkg-%v-bun-%v-tag-%v/", parms.PackageName, parms.BundleName, parms.TagName)

	err = os.RemoveAll(tmpDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		//Don't report on cleanup of directory that didn't exist
	} else {
		LogInfo(fmt.Sprintf("Cleaned up previous temp director %v", tmpDir))
	}

	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = os.Mkdir(tmpDir, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		LogOK(fmt.Sprintf("Temp directory %v already exists", tmpDir))
	} else {
		LogOK(fmt.Sprintf("Created temp directory %v", tmpDir))
	}
	defer func() {
		if !VerboseLogging {
			err = os.RemoveAll(tmpDir)
			if err != nil {
				LogError(fmt.Sprintf("Could not remove temp directory from %v:%v", tmpDir, err))
			}
			LogOK(fmt.Sprintf("Removed temp directory from %v", tmpDir))
		} else {
			LogInfo(fmt.Sprintf("Temp directory %v was not deleted as verbose logging is turned on", tmpDir))
		}
	}()

	tmpZipFile := fmt.Sprintf("%vTmpRtBun.zip", tmpDir)

	downloadURL := fmt.Sprintf("%v/%v.zip", parms.DownloadURL, parms.TagName)

	err = DownloadZipFromURL(downloadURL, tmpZipFile)
	if err != nil {
		LogError(err.Error())
		return err
	}

	outputZipDir, err := UnzipFromFileToDir(tmpZipFile, tmpDir)
	if err != nil {
		log.Fatal(err)
	}
	LogOK(fmt.Sprintf("Unzipped runtime bundle template to %v%v", tmpDir, outputZipDir))

	err = os.Remove(tmpZipFile)
	if err != nil {
		LogError(err.Error())
		return err
	}
	LogOK(fmt.Sprintf("Removed temp zip file %v", tmpZipFile))

	totalTransformsCount := 9

	newOutputUnzipDir := tmpDir + parms.BundleName + "/"
	err = os.Rename(outputZipDir, newOutputUnzipDir)
	if err != nil {
		LogError(err.Error())
		return err
	}
	LogOK(fmt.Sprintf("Transform 1 of %v: Rule rename template directory. Renamed from %v %v", totalTransformsCount, outputZipDir, newOutputUnzipDir))

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
	LogOK(fmt.Sprintf("Transform 2 of %v: Rule adjust group and artifact id in pom.xml. Adjusted with %v and %v with %v bytes at %v", totalTransformsCount, parms.PackageName, parms.BundleName, len(theStr), theFile))

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
	LogOK(fmt.Sprintf("Transform 3 of %v: Rule 'adjust 4 uses in skaffold.yaml. Adjusted to %v with %v bytes at %v", totalTransformsCount, parms.BundleName, len(theStr), theFile))

	theFile = newOutputUnzipDir + "Jenkinsfile"
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "example-runtime-bundle", parms.BundleName, 3)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	LogOK(fmt.Sprintf("Transform 4 of %v: Rule adjust 3 uses in Jenkinsfile. Adjusted to %v with %v bytes at %v", totalTransformsCount, parms.BundleName, len(theStr), theFile))

	renameFrom := newOutputUnzipDir + "charts/example-runtime-bundle"
	renameTo := newOutputUnzipDir + "charts/" + parms.BundleName
	err = os.Rename(renameFrom, renameTo)
	if err != nil {
		LogError(err.Error())
		return err
	}
	LogOK(fmt.Sprintf("Transform 5 of %v: Rule rename the charts folder. Adjusted to %v", totalTransformsCount, renameTo))

	theFile = newOutputUnzipDir + "src/main/resources/application.properties"
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "spring.application.name=${ACT_RB_APP_NAME:rb-my-app}", fmt.Sprintf("spring.application.name=${ACT_RB_APP_NAME:%v}", parms.BundleName), 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	LogOK(fmt.Sprintf("Transform 6 of %v: Rule change the spring.application.name. Adjusted to %v with %v bytes at %v", totalTransformsCount, parms.BundleName, len(theStr), theFile))

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
	LogOK(fmt.Sprintf("Transform 7 of %v: Rule change helm service name and repository to bundle name and tag to 'latest'. Adjusted to %v with %v bytes at %v", totalTransformsCount, parms.BundleName, len(theStr), theFile))

	theFile = newOutputUnzipDir + fmt.Sprintf("charts/%v/Chart.yaml", parms.BundleName)
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "example-runtime-bundle", parms.BundleName, 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return err
	}
	LogOK(fmt.Sprintf("Transform 8 of %v: Rule change helm name. Adjusted to %v with %v bytes at %v", totalTransformsCount, parms.BundleName, len(theStr), theFile))

	err = os.Rename(newOutputUnzipDir, finalOutDirectory)
	if err != nil {
		LogError(err.Error())
		return err
	}
	LogOK(fmt.Sprintf("Transform 9 of %v: Rule move the folder to the desired destination. Adjusted to %v", totalTransformsCount, finalOutDirectory))
	return nil
}

// //GenerateFromGenericQuickStartHandling is the interface for downloading and processing a template application from https://github.com/Activiti/activiti-cloud-runtime-bundle-quickstart/ or https://github.com/Activiti/activiti-cloud-connector-quickstart/
// type GenerateFromGenericQuickStartHandling func(parms GenercQuickstartDownloadParms) error

//GenercQuickstartDownloadParms are parameters for download handling
type GenercQuickstartDownloadParms struct {
	BundleName           string
	PackageName          string
	ProjectName          string
	TagName              string
	DestinationDir       string
	DownloadURL          string
	AdditionalTransCount int
}

//GenerateFromGenericQuickStartHandler downloads and processes a template application from https://github.com/Activiti/activiti-cloud-runtime-bundle-quickstart/ or https://github.com/Activiti/activiti-cloud-connector-quickstart/
func GenerateFromGenericQuickStartHandler(parms GenercQuickstartDownloadParms) (localTransformsCount int, err error) {
	destDir := EnsureFowardSlashAtStringEnd(parms.DestinationDir)

	finalOutDirectory := destDir + parms.BundleName

	finalOutDirectoryExists := true
	if _, err := os.Stat(finalOutDirectory); os.IsNotExist(err) {
		finalOutDirectoryExists = false
	}
	if finalOutDirectoryExists {
		msg := fmt.Sprintf("Final output directory is calculated as '%v' which might result in loss of data, double check your configuration and try again", finalOutDirectory)
		return 0, errors.New(msg)
	}

	tmpDir := destDir + fmt.Sprintf(".acdu-tmp-prj-%v-pkg-%v-bun-%v-tag-%v/", parms.ProjectName, parms.PackageName, parms.BundleName, parms.TagName)

	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = os.Mkdir(tmpDir, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, err
		}
		LogOK(fmt.Sprintf("Temp directory %v already exists", tmpDir))
	} else {
		LogOK(fmt.Sprintf("Created temp directory %v", tmpDir))
	}
	defer func() {
		if !VerboseLogging {
			err = os.RemoveAll(tmpDir)
			if err != nil {
				LogError(fmt.Sprintf("Could not remove temp directory from %v:%v", tmpDir, err))
			}
			LogOK(fmt.Sprintf("Removed temp directory from %v", tmpDir))
		} else {
			LogInfo(fmt.Sprintf("Temp directory %v was not deleted as verbose logging is turned on", tmpDir))
		}
	}()

	tmpZipFile := fmt.Sprintf("%vTmp.zip", tmpDir)

	//downloadURL := fmt.Sprintf("%v/%v.zip", parms.DownloadURL, parms.TagName)
	downloadURL := fmt.Sprintf("%v/archive/%v.zip", parms.DownloadURL, "master") //Watch the git repo, if it changes to tagging we should use 'parms.TagName' for the second parm
	err = DownloadZipFromURL(downloadURL, tmpZipFile)
	if err != nil {
		LogError(err.Error())
		return 0, err
	}

	outputZipDir, err := UnzipFromFileToDir(tmpZipFile, tmpDir)
	if err != nil {
		log.Fatal(err)
	}
	LogOK(fmt.Sprintf("Unzipped runtime bundle template to %v%v", tmpDir, outputZipDir))

	if !VerboseLogging {
		err = os.Remove(tmpZipFile)
		if err != nil {
			LogError(err.Error())
			return 0, err
		}
		LogOK(fmt.Sprintf("Removed temp zip file %v", tmpZipFile))
	} else {
		LogInfo(fmt.Sprintf("Temp zip %v was not deleted as verbose logging is turned on", tmpZipFile))
	}

	//Begin Transformation
	localTransformsCount = 7
	totalTransformsCount := localTransformsCount + parms.AdditionalTransCount

	newOutputUnzipDir := tmpDir + parms.BundleName + "/"
	err = os.Rename(outputZipDir, newOutputUnzipDir)
	if err != nil {
		return 0, err
	}
	LogOK(fmt.Sprintf("RT Bundle Transform 1 of %v: Rule rename template directory. Renamed from %v %v", totalTransformsCount, outputZipDir, newOutputUnzipDir))

	theFile := newOutputUnzipDir + "pom.xml"
	theBytes, err := ioutil.ReadFile(theFile)
	if err != nil {
		return 0, err
	}
	theStr := string(theBytes)
	theStr = strings.Replace(theStr, "REPLACE_ME_APP_NAME", parms.BundleName, -1)
	theStr = strings.Replace(theStr, "<groupId>org.activiti.cloud</groupId>", fmt.Sprintf("<groupId>%v</groupId>", parms.PackageName), 1)
	regex := regexp.MustCompile(`(?ms)(.*<activiti-cloud-dependencies\.version>)[a-z|\.|\-|0-9]*(<\/activiti-cloud-dependencies\.version>.*)`)
	match := regex.FindStringSubmatch(theStr)
	theStr = fmt.Sprintf("%v%v%v", match[1], parms.TagName, match[2])
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return 0, err
	}
	LogOK(fmt.Sprintf("RT Bundle Transform 2 of %v: Rule adjust group id, artifact id, REPLACE_ME_APP_NAME, and activiti-cloud-dependencies  in pom.xml. Adjusted with %v, %v, %v, and %v writing %v bytes at %v", totalTransformsCount, parms.PackageName, parms.BundleName, parms.ProjectName, parms.TagName, len(theStr), theFile))

	theFile = newOutputUnzipDir + "src/main/resources/application.properties"
	theBytes, err = ioutil.ReadFile(theFile)
	if err != nil {
		return 0, err
	}
	theStr = string(theBytes)
	// theStr = strings.Replace(theStr, "spring.application.name=${ACT_RB_APP_NAME:rb-my-app}", fmt.Sprintf("spring.application.name=${ACT_RB_APP_NAME:%v}", parms.BundleName), 1)
	regex = regexp.MustCompile(`(?ms)(.*spring.application.name=)[\$|\{|\}|_|A-Z|:|a-z|\.|\-|0-9]*(.*)`)
	match = regex.FindStringSubmatch(theStr)
	theStr = fmt.Sprintf("%v%v%v", match[1], parms.ProjectName, match[2])
	regex = regexp.MustCompile(`(?ms)(.*activiti.cloud.application.name=)[\$|\{|\}|_|A-Z|:|a-z|\.|\-|0-9]*(.*)`)
	match = regex.FindStringSubmatch(theStr)
	theStr = fmt.Sprintf("%v%v%v", match[1], fmt.Sprintf("${ACT_RB_APP_NAME:%v}", parms.ProjectName), match[2])
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return 0, err
	}
	LogOK(fmt.Sprintf("RT Bundle Transform 3 of %v: Rule change the spring.application.name and activiti.cloud.application.name. Adjusted to %v and %v with %v bytes at %v", totalTransformsCount, parms.BundleName, parms.ProjectName, len(theStr), theFile))

	renameFrom := newOutputUnzipDir + "charts/REPLACE_ME_APP_NAME"
	renameTo := newOutputUnzipDir + "charts/" + parms.BundleName
	err = os.Rename(renameFrom, renameTo)
	if err != nil {
		LogError(err.Error())
		return 0, err
	}
	LogOK(fmt.Sprintf("RT Bundle Transform 4 of %v: Rule rename the charts folder. Adjusted to %v", totalTransformsCount, renameTo))

	theFile = newOutputUnzipDir + fmt.Sprintf("charts/%v/values.yaml", parms.BundleName)
	theBytes, err = ioutil.ReadFile(theFile)
	if err != nil {
		return 0, err
	}
	theStr = string(theBytes)
	//theStr = strings.Replace(theStr, "name: rb-my-app", fmt.Sprintf("name: %v", parms.BundleName), 1)
	theStr = strings.Replace(theStr, "repository: draft", fmt.Sprintf("repository: %v", parms.BundleName), 1)
	theStr = strings.Replace(theStr, "tag: dev", "tag: latest", 1)
	theStr = strings.Replace(theStr, "REPLACE_ME_APP_NAME", fmt.Sprintf("%v-%v", parms.ProjectName, parms.BundleName), 1) //replaces 'service.name'
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return 0, err
	}
	LogOK(fmt.Sprintf("RT Bundle Transform 5 of %v: Rule change helm service name and repository to bundle name and tag to 'latest'. Adjusted to %v with %v bytes at %v", totalTransformsCount, parms.BundleName, len(theStr), theFile))

	theFile = newOutputUnzipDir + fmt.Sprintf("charts/%v/Chart.yaml", parms.BundleName)
	theBytes, err = ioutil.ReadFile(theFile)
	theStr = string(theBytes)
	theStr = strings.Replace(theStr, "REPLACE_ME_APP_NAME", parms.BundleName, 1)
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(theFile, []byte(theStr), os.ModePerm)
	if err != nil {
		return 0, err
	}
	LogOK(fmt.Sprintf("RT Bundle Transform 6 of %v: Rule change helm name. Adjusted to %v with %v bytes at %v", totalTransformsCount, parms.BundleName, len(theStr), theFile))

	err = os.Rename(newOutputUnzipDir, finalOutDirectory)
	if err != nil {
		LogError(err.Error())
		return 0, err
	}
	LogOK(fmt.Sprintf("RT Bundle Transform 7 of %v: Rule move the folder to the desired destination. Adjusted to %v", totalTransformsCount, finalOutDirectory))
	return localTransformsCount, nil
}
