package api

import (
	"fmt"
	"strings"
)

func Render(result Result) string {
	var builder strings.Builder
	if result.OK {
		builder.WriteString("OK ")
	} else {
		builder.WriteString("ERROR ")
	}
	builder.WriteString(result.Kind)
	builder.WriteString(": ")
	builder.WriteString(result.Message)
	if result.Error != "" {
		builder.WriteString(" [")
		builder.WriteString(result.Error)
		builder.WriteString("]")
	}
	for _, row := range result.Rows {
		builder.WriteString("\n- ")
		builder.WriteString(row)
	}
	return builder.String()
}

func RenderTable(title string, rows []string) string {
	var builder strings.Builder
	builder.WriteString(title)
	for index, row := range rows {
		builder.WriteString(fmt.Sprintf("\n%02d. %s", index+1, row))
	}
	return builder.String()
}

func RenderErrorLogs(rows []string) string {
	if len(rows) == 0 {
		return "No error logs"
	}
	return RenderTable("Error logs", rows)
}
