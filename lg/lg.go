package lg

import (
	"io"
	"log"
	"os"

	"github.com/natefinch/lumberjack"
)

var (
	logSink io.Writer
	dl      *log.Logger
	el      *log.Logger
	il      *log.Logger
	wl      *log.Logger
)

func init() {
	if logSink == nil {
		logSink = os.Stderr
	}
	dl = log.New(logSink, "DEBUG: ", log.Ldate|log.Ltime)
	el = log.New(logSink, "ERROR: ", log.Ldate|log.Ltime)
	il = log.New(logSink, "INFO: ", log.Ldate|log.Ltime)
	wl = log.New(logSink, "WARNING: ", log.Ldate|log.Ltime)
}

// Init specifies the name of the log file that stores the log.
//
// The `level` parameter determines which log statements are sent to the log
// file in addition to STDOUT.
func Init(filename, level string) {
	if filename != "" {
		lj := &lumberjack.Logger{
			Filename: filename,
			MaxSize:  1,
		}
		logSink = io.MultiWriter(lj, os.Stdout)
	}
	switch level {
	case "DEBUG":
		dl.SetOutput(logSink)
		il.SetOutput(logSink)
		wl.SetOutput(logSink)
		el.SetOutput(logSink)
	case "INFO":
		dl.SetOutput(os.Stdout)
		il.SetOutput(logSink)
		wl.SetOutput(logSink)
		el.SetOutput(logSink)
	case "WARNING":
		dl.SetOutput(os.Stdout)
		il.SetOutput(os.Stdout)
		wl.SetOutput(logSink)
		el.SetOutput(logSink)
	default: //inclusive of ERROR
		dl.SetOutput(os.Stdout)
		il.SetOutput(os.Stdout)
		wl.SetOutput(os.Stdout)
		el.SetOutput(logSink)
	}
}

// Debug writes an error entry in the manner of fmt.Printf, prefixed with
// "DEBUG: ".
func Debug(format string, v ...interface{}) {
	dl.Printf(format+"\r", v...)
}

// Error writes an error entry in the manner of fmt.Printf, prefixed with
// "ERROR: ".
func Error(format string, v ...interface{}) {
	el.Printf(format+"\r", v...)
}

// Info writes an error entry in the manner of fmt.Printf, prefixed with
// "INFO: ".
func Info(format string, v ...interface{}) {
	il.Printf(format+"\r", v...)
}

// Warning writes an error entry in the manner of fmt.Printf, prefixed with
// "WARNING: ".
func Warning(format string, v ...interface{}) {
	wl.Printf(format+"\r", v...)
}
