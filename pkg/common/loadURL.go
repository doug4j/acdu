package common

import (
	"fmt"

	"github.com/skratchdot/open-golang/open"
)

func LoadURLInBrowser(url, oneWordDescription, optionalHint string) error {
	LogWorking(fmt.Sprintf("Loading %v [%v]", oneWordDescription, url))
	err := open.Start(url)
	if err != nil {
		err = fmt.Errorf("Could not start %v URL: %v", oneWordDescription, err)
		return err
	}
	if optionalHint == "" {
		LogOK(fmt.Sprintf("URL [%v] send to open browser.", url))
	} else {
		LogOK(fmt.Sprintf("URL [%v] send to open browser. Hint: %v", url, optionalHint))
	}
	return nil
}
