package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

/**
日志记录器(Logger)是日志处理的核心组件，具有5种正常级别(Level)。:
1.Level DEBUG :
DEBUG Level指出细粒度信息事件对调试应用程序是非常有帮助的。
2.Level INFO
INFO level表明 消息在粗粒度级别上突出强调应用程序的运行过程。
3.Level WARN
WARN level表明会出现潜在错误的情形。
4.Level ERROR
ERROR level指出虽然发生错误事件，但仍然不影响系统的继续运行。
5.Level FATAL
FATAL level指出每个严重的错误事件将会导致应用程序的退出。
*/

const (
	LevelDebug = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
)

const (
	log_info  = "[INFO]"
	log_warn  = "[WARN]"
	log_debug = "[DEBUG]"
	log_error = "[ERROR]"
	log_fatal = "[FATAL]"
)

var (
	level       = LevelInfo
	abc_ops_log = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
)

func init() {
	logFile := LogDir + PathSeparator + getLogDay(time.Now()) + Log_file_suffix

	if err := MkdirByFile(logFile); nil != err {
		Log_Fatal(err)
	}

	logfile, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Log_Fatal(err)
	}
	abc_ops_log.SetOutput(logfile)
	abc_ops_log.SetFlags(log.LstdFlags)
}

func getLogDate() string {
	t := time.Now()
	return Format_Date(t, "2006/01/02 15:04:05")
}

func getLogDay(time time.Time) string {
	return time.Format("2006-01-02")
}

func Log_level() int {
	return level
}

func Set_log_level(l int) {
	level = l
}

func SetLogger(l *log.Logger) {
	abc_ops_log = l
}

func Log_Info(err ...interface{}) {
	if level <= LevelInfo {
		logPrintln(log_info, err...)
	}
}

func Log_Fatal(err ...interface{}) {
	if level <= LevelFatal {
		logPrintln(log_fatal, err...)
	}
}

func Log_Warn(err ...interface{}) {
	if level <= LevelWarn {
		logPrintln(log_warn, err...)
	}
}

func Log_Debug(err ...interface{}) {
	if level <= LevelDebug {
		logPrintln(log_debug, err...)
	}
}

func Log_Err(err ...interface{}) {
	if level <= LevelError {
		logPrintln(log_error, err...)
	}
}

func logPrintln(msgType string, err ...interface{}) {
	errStr := fmt.Sprintf("%v", err)
	funcName, file, line, ok := runtime.Caller(2)
	if ok {
		name := runtime.FuncForPC(funcName).Name()
		file = filepath.Base(file)
		fmt.Fprintln(os.Stdout, msgType, getLogDate(), file, line, name, errStr)
		abc_ops_log.SetPrefix(msgType)
		abc_ops_log.Println(file, line, name, err)
	} else {
		fmt.Fprintln(os.Stdout, msgType, getLogDate(), errStr)
		abc_ops_log.SetPrefix(msgType)
		abc_ops_log.Println(err)
	}
}
