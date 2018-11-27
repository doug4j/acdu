package common

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type SimpleProjectObjectModel struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
}

func GetProjectObjectModelFromCurrentDir() (answer SimpleProjectObjectModel, err error) {
	rawBytes, err := ioutil.ReadFile("pom.xml")
	if err != nil {
		LogError(fmt.Sprintf("Cannot read pom.xml:%v", err))
		return answer, err
	}
	err = xml.Unmarshal(rawBytes, &answer)
	if err != nil {
		LogError(fmt.Sprintf("Cannot parse pom.xml:%v", err))
		return answer, err
	}
	return answer, nil
}
