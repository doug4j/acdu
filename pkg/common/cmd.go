package common

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func Command(program string, args []string, dir string, shortDesc string) (string, error) {
	start := time.Now()
	fullCommand := strings.Join(append([]string{program}, args...), " ")
	LogWorking(fmt.Sprintf("%v [%v]...", shortDesc, fullCommand))
	cmd := exec.Command(program, args...)
	cmd.Dir = dir
	out, err := exec.Command(program, args...).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Cannot run command [%v]. Output:\n%verror:[%v]", fullCommand, fmt.Sprintf("%s", out), err)
		LogError(err.Error())
		return "", err
	}
	if VerboseLogging {
		LogOK(fmt.Sprintf("%v [%v]. Output:\n%v", fullCommand, shortDesc, fmt.Sprintf("%s", out)))
	} else {
		LogOK(fmt.Sprintf("%v [%v]", shortDesc, fullCommand))
	}
	end := time.Now()
	elapsed := end.Sub(start)
	LogTime(fmt.Sprintf("%v elapsed time: %v", shortDesc, elapsed.Round(time.Millisecond)))
	return fmt.Sprintf("%s", out), nil
}
