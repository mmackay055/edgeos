package logging

import (
        "fmt"
)

type SystemdLokiLog struct {}

func (sl *SystemdLokiLog) LogEmerg(message string) string {
	return fmt.Sprintf("level=emergency msg=\"%s\"", message)
}

func (sl *SystemdLokiLog) LogAlert(message string) string {
	return fmt.Sprintf("level=alert msg=\"%s\"", message)
}

func (sl *SystemdLokiLog) LogCrit(message string) string {
	return fmt.Sprintf("level=critical msg=\"%s\"", message)
}

func (sl *SystemdLokiLog) LogErr(message string) string {
	return fmt.Sprintf("level=error msg=\"%s\"", message)
}

func (sl *SystemdLokiLog) LogWarning(message string) string {
	return fmt.Sprintf("level=warning msg=\"%s\"", message)
}

func (sl *SystemdLokiLog) LogNotice(message string) string {
	return fmt.Sprintf("level=notice msg=\"%s\"", message)
}

func (sl *SystemdLokiLog) LogInfo(message string) string {
	return fmt.Sprintf("level=info msg=\"%s\"", message)
}

func (sl *SystemdLokiLog) LogDebug(message string) string {
	return fmt.Sprintf("level=debug msg=\"%s\"", message)
}
