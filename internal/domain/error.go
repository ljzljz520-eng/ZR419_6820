package domain

import (
	"fmt"
	"strings"
)

type ErrorLog struct {
	ID      string
	Code    string
	Message string
	Context string
	Level   string
}

func NewErrorLog(id, code, message, context string) ErrorLog {
	return ErrorLog{ID: strings.TrimSpace(id), Code: strings.TrimSpace(code), Message: strings.TrimSpace(message), Context: strings.TrimSpace(context), Level: "error"}
}

func (e ErrorLog) ValidateShape() error {
	if e.ID == "" || e.Code == "" || e.Message == "" {
		return fmt.Errorf("error log %s has incomplete fields", e.ID)
	}
	return nil
}

func (e ErrorLog) Display() string {
	return fmt.Sprintf("[%s] %s: %s (%s)", e.Level, e.Code, e.Message, e.Context)
}
