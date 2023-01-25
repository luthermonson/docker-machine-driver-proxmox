package proxmox

import (
	"os"
	"strconv"

	"github.com/rancher/machine/libmachine/log"
)

type leveledLogger struct{}

var logger leveledLogger
var debug bool

func init() {
	for _, f := range os.Args {
		if f == "-D" || f == "--debug" || f == "-debug" {
			os.Setenv("MACHINE_DEBUG", "1")
			log.SetDebug(true)
			debug = true
			return
		}
	}

	// check env
	debugEnv := os.Getenv("MACHINE_DEBUG")
	if debugEnv != "" {
		showDebug, err := strconv.ParseBool(debugEnv)
		if err != nil {
			log.Errorf("Error parsing boolean value from MACHINE_DEBUG: %s", err)
			return
		}
		log.SetDebug(showDebug)
		debug = true
	}
}

func (l leveledLogger) Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func (l leveledLogger) Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func (l leveledLogger) Warnf(format string, v ...interface{}) {
	if debug {
		log.Warnf(format, v...)
	}
}

func (l leveledLogger) Debugf(format string, v ...interface{}) {
	if debug {
		log.Debugf(format, v...)
	}
}
