package logging

type Logger struct {
        format LogFormat
}

func (l *Logger) SetFormat(format LogFormat) {
        l.format = format
}

func (l *Logger) LogEmerg (message string) string {
        return l.format.LogEmerg(message)
}

func (l *Logger) LogAlert (message string) string {
        return l.format.LogAlert(message)
}

func (l *Logger) LogCrit (message string) string {
        return l.format.LogCrit(message)
}

func (l *Logger) LogErr (message string) string {
        return l.format.LogErr(message)
}

func (l *Logger) LogWarning (message string) string {
        return l.format.LogWarning(message)
}

func (l *Logger) LogNotice (message string) string {
        return l.format.LogNotice(message)
}

func (l *Logger) LogInfo (message string) string {
        return l.format.LogInfo(message)
}

func (l *Logger) LogDebug (message string) string {
        return l.format.LogDebug(message)
}
