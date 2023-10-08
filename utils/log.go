package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"time"
)

const (
	foreOrg = 0
	foreBlk = iota + 29
	foreRed
	foreGrn
	foreYlw
	foreBle
	forePrp
	foreCyn
	foreWht
)

const (
	backOrg = 0
	backBlk = iota + 39
	backRed
	backGrn
	backYlw
	backBle
	backPrp
	backCyn
	backWht
)

type LogLevel int
const (
	DBUG LogLevel = iota
	INFO
	WARN
	ERRO
	FATL
)
var levelMarker = map[LogLevel]string {
	DBUG: "[DBUG]",
	INFO: "[INFO]",
	WARN: "[WARN]",
	ERRO: "[ERRO]",
	FATL: "[FATL]",
}
var levelColor = map[LogLevel][]int {
	DBUG: { foreBle, backOrg },
	INFO: { foreGrn, backOrg },
	WARN: { foreYlw, backOrg },
	ERRO: { foreRed, backOrg },
	FATL: { forePrp, backOrg },
}

var DumpToFile bool
var VerboseLog bool

var logFile *os.File
var logFileWriter *bufio.Writer
var localTimeLocation *time.Location

func init() {
	// log to file is set to on by default
	DumpToFile = false
	VerboseLog = true

	// use Shanghai time location
	timeLocation, timeErr := time.LoadLocation("Asia/Shanghai")
	if timeErr != nil {
		fmt.Print("failed to set time location")
	}
	localTimeLocation = timeLocation
}

func InitLogFile(logFileDir string) error {
	var logFileName string = fmt.Sprintf(
		"%s.log",
		fmt.Sprintf(
			"%04d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(),
		),
	)

	var openErr error
	logFile, openErr = os.OpenFile(
		fmt.Sprintf("%s/%s", logFileDir, logFileName),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0660,
	)
	if openErr != nil {
		errMkdirAll := os.MkdirAll(logFileDir, 0755)
		if errMkdirAll != nil {
			panic(errMkdirAll.Error())
		}
		logFile, openErr = os.OpenFile(
			fmt.Sprintf("%s/%s", logFileDir, logFileName),
			os.O_WRONLY|os.O_CREATE|os.O_APPEND,
			0666,
		)
		if openErr != nil {
			return errors.New("failed to open log file after creating path")
		}
	}

	return nil
}

func dumpToLogFile(logMsg string) {
	if logFileWriter == nil {
		logFileWriter = bufio.NewWriter(logFile)
	}
	logFileWriter.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
	logFileWriter.Write([]byte(" " + logMsg + "\n"))
	if VerboseLog {
		logFileWriter.Write([]byte(FuncLocString("    ", 3) + "\n"))
	}
	
	if err := logFileWriter.Flush(); err != nil {
		fmt.Println("failed to dump log to file")
	}
}

func colorString(str string, fore, back int) string {
	if back == backOrg {
		return fmt.Sprintf("\033[%dm%s\033[0m", fore, str)
	}
	return fmt.Sprintf("\033[1;%d;%dm%s\033[0m", fore, back, str)
}

func concateMsgs(msgs ...interface{}) string {
	var concateMsgs string = ""
	for idx, msg := range msgs {
		if reflect.TypeOf(msg).Kind() == reflect.String {
			concateMsgs += msg.(string)
		} else {
			concateMsgs += fmt.Sprintf("%#v", msg)
		}
		if idx != len(msgs)-1 {
			concateMsgs += " "
		}
	}
	return concateMsgs
}

func FuncLocString(prefix string, skip int) string {
	var location = ""
	pc, codePath, codeLine, ok := runtime.Caller(skip)
	if !ok {
		location = fmt.Sprintf("%s%s:%s %s", prefix, "unknown location", "unknown line number", "unknown function")
	} else {
		location = fmt.Sprintf("%s%s:%d %s", prefix, codePath, codeLine, runtime.FuncForPC(pc).Name())
	}

	return location
}

func Log(level LogLevel, msgs ...interface{}) {
	concateMsgs := concateMsgs(msgs...)
	logMsg := fmt.Sprintf("%s %s", colorString(
			levelMarker[level],
			levelColor[level][0],
			levelColor[level][1],
		),
		concateMsgs,
	)
	fmt.Println(logMsg)
	if VerboseLog {
		fmt.Println(FuncLocString("    ", 2))
	}
	if DumpToFile {
		dumpToLogFile(fmt.Sprintf("%s %s", levelMarker[level], concateMsgs))
	}
}
