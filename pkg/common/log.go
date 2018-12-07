package common

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

//Adding options for not using Emoji (UTF-8/16) on windows due to https://stackoverflow.com/questions/44054983/how-to-output-emoji-to-console-in-node-js-on-windows

//LogOK does Println with conditional OK prefix based on platform
func LogOK(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(formatOK(msg))
		return
	}
	log.Println("👍 " + formatOK(msg))
}

func formatOK(msg string) string {
	return fmt.Sprintf(logPrefixPad("[OK]") + msg)
}

//LogTime does Println with conditional TIME prefix based on platform
func LogTime(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(formatTime(msg))
		return
	}
	log.Println("⏱  " + formatTime(msg))
}

func formatTime(msg string) string {
	return fmt.Sprintf(logPrefixPad("[TIME]") + msg)
}

//LogWorking does Println with conditional WORKING prefix based on platform
func LogWorking(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(formatWorking(msg))
		return
	}
	log.Println("⏳ " + formatWorking(msg))
}

func formatWorking(msg string) string {
	return logPrefixPad("[WORKING]") + msg
}

//LogWaitingForUser does Println with conditional USERINPUT prefix based on platform
func LogWaitingForUser(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(formatWaitingForUser(msg))
		return
	}
	log.Println("⏳ " + formatWaitingForUser(msg))
}

func formatWaitingForUser(msg string) string {
	return logPrefixPad("[USERINPUT]") + msg
}

//LogError does Println with conditional ERROR prefix based on platform
func LogError(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(formatError(msg))
		return
	}
	log.Println("⛔ " + formatError(msg))
}

func formatError(msg string) string {
	return logPrefixPad("[ERROR]") + msg
}

//LogWarn does Println with conditional WARN prefix based on platform
func LogWarn(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(formatWarn(msg))
		return
	}
	log.Println("🤔 " + formatWarn(msg))
}

func formatWarn(msg string) string {
	return logPrefixPad("[WARN]") + msg
}

//LogExit does Println with conditional EXIT prefix based on platform and then ends the process with os.Exit(1)
func LogExit(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(formatExit(msg))
		os.Exit(1)
		return
	}
	log.Println("😢 " + formatExit(msg))
	os.Exit(1)
}

func formatExit(msg string) string {
	return logPrefixPad("[EXIT]") + msg
}

//LogNotImplemented does Println with conditional UNFINISHED prefix based on platform
func LogNotImplemented(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(formatNotImplemented(msg))
		return
	}
	log.Println("🧐 " + formatNotImplemented(msg))
}

func formatNotImplemented(msg string) string {
	return logPrefixPad("[UNFINISHED]") + msg
}

//LogInfo does Println with conditional INFO prefix based on platform
func LogInfo(msg string) {
	if runtime.GOOS == "windows" {
		log.Println(formatInfo(msg))
		return
	}
	log.Println("ℹ️  " + formatInfo(msg))
}

func formatInfo(msg string) string {
	return fmt.Sprintf(logPrefixPad("[INFO]") + msg)
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
