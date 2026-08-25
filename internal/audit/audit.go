package audit

import (
	"fmt"

	"hospitalportal/internal/domain"
	"hospitalportal/internal/store"
)

type Logger struct {
	store *store.Store
}

func NewLogger(s *store.Store) *Logger {
	return &Logger{store: s}
}

func (l *Logger) Record(code, message, context string) (domain.ErrorLog, error) {
	count, err := l.store.CountErrors()
	if err != nil {
		return domain.ErrorLog{}, err
	}
	log := domain.NewErrorLog(fmt.Sprintf("err-%04d", count+1), code, message, context)
	if err := l.store.SaveErrorLog(log); err != nil {
		return domain.ErrorLog{}, err
	}
	return log, nil
}

func (l *Logger) Recent(limit int) ([]domain.ErrorLog, error) {
	logs, err := l.store.ListErrorLogs()
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit >= len(logs) {
		return logs, nil
	}
	return logs[len(logs)-limit:], nil
}

func (l *Logger) HasCode(code string) (bool, error) {
	logs, err := l.store.ListErrorLogs()
	if err != nil {
		return false, err
	}
	for _, log := range logs {
		if log.Code == code {
			return true, nil
		}
	}
	return false, nil
}
