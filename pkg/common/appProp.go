package common

import (
	"io/ioutil"
	"regexp"
)

type SimpleProp struct {
	SpringAppName string
}

func GetSimpleProp(filePath string) (SimpleProp, error) {
	theBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return SimpleProp{}, err
	}
	theStr := string(theBytes)
	springAppNameRegex := regexp.MustCompile(`(?ms)(.*spring.application.name=)([\$|\{|\}|_|A-Z|:|a-z|\.|\-|0-9]*)(.*)`)
	springAppNameMatch := springAppNameRegex.FindStringSubmatch(theStr)
	answer := SimpleProp{
		SpringAppName: springAppNameMatch[2],
	}
	return answer, nil
}
