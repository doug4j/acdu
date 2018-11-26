package common

import (
	"log"
	"os"
	"runtime"
	"strings"
)

//Adding options for not using Emoji (UTF-8/16) on windows due to https://stackoverflow.com/questions/44054983/how-to-output-emoji-to-console-in-node-js-on-windows

//LogOK does Println with conditional OK prefix based on platform
func LogOK(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(logPrefixPad("[OK]") + msg)
		return
	}
	log.Println("👍 " + msg)
}

//LogTime does Println with conditional TIME prefix based on platform
func LogTime(msg string) {
	if runtime.GOOS == "windows" {
		log.Printf(logPrefixPad("[TIME]") + msg)
		return
	}
	log.Printf("⏱ " + msg)
}

//LogWorking does Println with conditional WORKING prefix based on platform
func LogWorking(msg string) {
	if runtime.GOOS == "windows" {
		log.Printf(logPrefixPad("[WORKING]") + msg)
		return
	}
	log.Printf("⏳ " + msg)
}

//LogWaitingForUser does Println with conditional USERINPUT prefix based on platform
func LogWaitingForUser(msg string) {
	if runtime.GOOS == "windows" {
		log.Printf(logPrefixPad("[USERINPUT]") + msg)
		return
	}
	log.Printf("⏳ " + msg)
}

//LogError does Println with conditional ERROR prefix based on platform
func LogError(msg string) {
	if runtime.GOOS == "windows" {
		log.Printf(logPrefixPad("[ERROR]") + msg)
		return
	}
	log.Printf("⛔ " + msg)
}

//LogExit does Println with conditional EXIT prefix based on platform and then ends the process with os.Exit(1)
func LogExit(msg string) {
	if runtime.GOOS == "windows" {
		log.Printf(logPrefixPad("[EXIT]") + msg)
		os.Exit(1)
		return
	}
	log.Printf("😢 " + msg)
	os.Exit(1)
}

//LogLogNotImplementedError does Println with conditional UNFINISHED prefix based on platform
func LogNotImplemented(msg string) {
	if runtime.GOOS == "windows" {
		log.Printf(logPrefixPad("[UNFINISHED]") + msg)
		os.Exit(1)
		return
	}
	log.Printf("🧐" + msg)
}

//LogInfo does Println with conditional INFO prefix based on platform
func LogInfo(msg string) {
	if runtime.GOOS == "windows" {
		log.Printf(logPrefixPad("[INFO]") + msg)
		return
	}
	log.Printf("ℹ️ " + msg)
}

const loggestLogPrefixPad = 13 //this is the longest log prefix plus a space

func logPrefixPad(msg string) string {
	//From https://github.com/git-time-metric/gtm/blob/master/util/string.go#L53-L88
	var overallLen = loggestLogPrefixPad
	var padStr = " "
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = msg + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}
