package logging

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultLogLevel        string = "error"
	defaultTimeStampFormat string = "2006-01-02 15:04:05"
)

var baseFilePath, _ = filepath.Abs("")

type LogConfig struct {
	FileEnabled    bool   `envconfig:"LOG_FILE_ENABLED" default:"false"`
	FileLocation   string `envconfig:"LOG_FILE_NAME" default:"server.log"`
	ConsoleEnabled bool   `envconfig:"LOG_CONSOLE_ENABLED" default:"true"`
}

type Logger struct {
	*logrus.Logger

	config  *LogConfig
	logFile *os.File
}

func (l *Logger) WithConfig(cfg *LogConfig) error {
	if l.logFile == nil {
		logFile, err := openLogFile(cfg.FileLocation)
		if err != nil {
			return err
		}
		l.logFile = logFile
	}

	// writing to log file
	var writer io.Writer
	writer = io.Writer(l.logFile)
	l.SetOutput(writer)

	// default error log level
	logLevel, _ := logrus.ParseLevel(defaultLogLevel)
	l.SetLevel(logLevel)
	
	l.config = cfg

	l.Info("Logging initialized")

	return nil
}

func openLogFile(loc string) (logFile *os.File, err error) {
	logFile, err = os.OpenFile(loc, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0)
	if err != nil {
		formattedErr := fmt.Errorf("failed to open log file (%v) for logging: %v", loc, err)
		Log.Error(formattedErr.Error())

		return nil, formattedErr
	}

	return logFile, nil
}

type LogFormat struct {
	TimestampFormat string
}

func (l LogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	// init buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// log level added
	// Level will be formatted as [INFO]:
	formattedLevel := fmt.Sprintf("[%-9v", strings.ToUpper(entry.Level.String()+"]"))
	b.WriteString(formattedLevel)

	// Time is formatted as 01-02-2006 15:04:05
	b.WriteString(entry.Time.Format(l.TimestampFormat) + " ")

	// caller func added
	caller := entry.Caller
	methodName := caller.Function
	callerFilePath, _ := filepath.Rel(baseFilePath, caller.File)
	formattedFile := fmt.Sprintf("[%v] ", callerFilePath)
	b.WriteString(formattedFile)
	methodChunks := strings.Split(methodName, "/")
	formattedMethod := fmt.Sprintf("[%v: %v]", methodChunks[len(methodChunks)-1], caller.Line)
	b.WriteString(formattedMethod)

	// entry message added
	if entry.Message != "" {
		formattedMessage := fmt.Sprintf(" - %v", entry.Message)
		b.WriteString(formattedMessage)
	}

	// write log data to buffer
	if len(entry.Data) > 0 {
		b.WriteString(" || ")
	}
	for key, value := range entry.Data {
		formattedData := fmt.Sprintf("%v={ %v }, ", key, value)
		b.WriteString(formattedData)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

var Log = NewLogger()

func NewLogger() *Logger {
	level, _ := logrus.ParseLevel(defaultLogLevel)
	logger := &Logger{
		Logger: &logrus.Logger{
			Out:          os.Stdout,
			Level:        level,
			ReportCaller: true,
			Formatter: &LogFormat{
				TimestampFormat: defaultTimeStampFormat,
			},
			Hooks: make(logrus.LevelHooks),
		},
	}

	return logger
}

func Initialize(cfg *LogConfig) error {
	return Log.WithConfig(cfg)
}
