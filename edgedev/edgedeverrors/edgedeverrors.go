package edgedeverrors

import (
        "fmt"
)
type LoginError struct {
        Message string
}

// TODO add option to format for systemd loki ingestion

func (e LoginError) Error() string {
        return fmt.Sprintf("login failed: %s", e.Message)
}

type SaveError struct {
        Message string
}

func (e SaveError) Error() string {
        return fmt.Sprintf("save failed: %s", e.Message)
}

type BackupError struct {
        Message string
}

func (e BackupError) Error() string {
        return fmt.Sprintf("backup failed: %s", e.Message)
}
