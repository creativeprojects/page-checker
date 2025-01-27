package lib

import (
	"fmt"

	"github.com/go-rod/rod/lib/utils"
)

type Logger struct{}

func (l *Logger) Println(a ...interface{}) {
	// does nothing for now
}

var _ utils.Logger = &Logger{}

func Verbose(cfg Flags, format string, a ...any) {
	if !cfg.Verbose {
		return
	}
	fmt.Printf(format, a...)
}
