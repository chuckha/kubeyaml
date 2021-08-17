package logging

import (
	"fmt"
	"io"
)

type Log struct {
	Writer io.Writer
}

func (l *Log) Debugf(format string, args ...interface{}) {
	wrappedFormat := fmt.Sprintf("[DEBUG]: %s", format)
	fmt.Fprintf(l.Writer, wrappedFormat, args...)
}

func (l *Log) Infof(format string, args ...interface{}) {
	wrappedFormat := fmt.Sprintf("[INFO]: %s", format)
	fmt.Fprintf(l.Writer, wrappedFormat, args...)
}
