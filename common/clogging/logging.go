package clogging

import (
	"io"
	"os"
	"sync"

	"github.com/op/go-logging"
)

const (
	pkgLogID      = "clogging"
	defaultFormat = "%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}"
	defaultLevel  = logging.INFO
)

var (
	logger        *logging.Logger
	modules       map[string]string
	once          sync.Once
	defaultOutput *os.File
)

func initLogger() {
	once.Do(func() {
		modules = make(map[string]string)
		InitBackend(SetFormat(defaultFormat), os.Stdout)
	})
}

func GetModuleLevel(module string) string {
	level := logging.GetLevel(module).String()
	return level
}

// MustGetLogger is used in place of `logging.MustGetLogger` to allow us to
// store a map of all modules and submodules that have loggers in the system.
func MustGetLogger(module string) *logging.Logger {
	initLogger()
	l := logging.MustGetLogger(module)
	if modules[module] == "" {
		modules[module] = GetModuleLevel(module)
	}
	return l
}

func SetFormat(formatSpec string) logging.Formatter {
	if formatSpec == "" {
		formatSpec = defaultFormat
	}
	return logging.MustStringFormatter(formatSpec)
}

// InitBackend sets up the logging backend based on
// the provided logging formatter and I/O writer.
func InitBackend(formatter logging.Formatter, output io.Writer) {
	backend := logging.NewLogBackend(output, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, formatter)
	logging.SetBackend(backendFormatter).SetLevel(defaultLevel, "")
}

func setModuleLevel(module string, level string) (string, error) {
	logLevel, err := logging.LogLevel(level)
	if err != nil {
		logger.Warningf("Invalid logging level '%s' - ignored", level)
	} else {
		logging.SetLevel(logging.Level(logLevel), module)
		modules[module] = logLevel.String()
		logger.Debugf("Module '%s' logger enabled for log level '%s'", module, logLevel)
	}
	return logLevel.String(), err
}