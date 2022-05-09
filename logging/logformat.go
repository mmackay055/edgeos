package logging

type LogFormat interface {
        LogEmerg(message string) string
        LogAlert(message string) string
        LogCrit(message string) string
        LogErr(message string) string
        LogWarning(message string) string
        LogNotice(message string) string
        LogInfo(message string) string
        LogDebug(message string) string
}

