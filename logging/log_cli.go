package logging

import (
	"fmt"
)

type CliLog struct {}

func (cf *CliLog) LogEmerg(message string) string {
	return fmt.Sprintf("emergency: %s", message)
}

func (cf *CliLog) LogAlert(message string) string {
	return fmt.Sprintf("alert: %s", message)
}

func (cf *CliLog) LogCrit(message string) string {
	return fmt.Sprintf("critical: %s", message)
}

func (cf *CliLog) LogErr(message string) string {
	return fmt.Sprintf("error: %s", message)
}

func (cf *CliLog) LogWarning(message string) string {
	return fmt.Sprintf("warning: %s", message)
}

func (cf *CliLog) LogNotice(message string) string {
	return fmt.Sprintf("notice: %s", message)
}

func (cf *CliLog) LogInfo(message string) string {
	return fmt.Sprintf("info: %s", message)
}

func (cf *CliLog) LogDebug(message string) string {
	return fmt.Sprintf("debug: %s", message)
}
