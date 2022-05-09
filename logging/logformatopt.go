package logging

import (
	"errors"
	"fmt"
)

const (
	systemdLoki = `systemdloki`
	cli         = `cli`
)

// EdgeDevOpt option definition
type LogFormatOpt struct {
	name   string
	format LogFormat
}

func (f *LogFormatOpt) String() string {
	return f.name
}

func (f *LogFormatOpt) Type() string {
	return "LogFormat"
}

func (f *LogFormatOpt) Set(val string) error {
	if len(val) < 1 {
		return nil
	}
	return f.setFormat(val)
}

func AvailableEdgeDevs() []string {
	available := [2]string{systemdLoki, cli}
	return available[:]
}

func (f *LogFormatOpt) setFormat(devname string) error {
	var format LogFormat

	switch devname {
	case systemdLoki:
		format = &SystemdLokiLog{}
		f.name = systemdLoki
        case cli:
                format = &CliLog{}
                f.name = cli
	default:
		return errors.New(fmt.Sprintf("format %s is not suppored",
			devname))
	}
	f.format = format
	return nil
}

func (f LogFormatOpt) Get() LogFormat {
	return f.format
}

func (f *LogFormatOpt) IsSet() bool {
	return f.format != nil
}
