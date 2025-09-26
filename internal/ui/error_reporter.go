package ui

import (
	"fmt"
	"strings"

	"github.com/Artemka007/derraform/internal/errors"
	"github.com/fatih/color"
)

type ErrorReporter struct {
	errors []*errors.TerraformError
}

func NewErrorReporter() *ErrorReporter {
	return &ErrorReporter{}
}

func (e *ErrorReporter) AddError(err *errors.TerraformError) {
	e.errors = append(e.errors, err)
}

func (e *ErrorReporter) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *ErrorReporter) PrintSummary() {
	if len(e.errors) == 0 {
		color.Green("✓ All operations completed successfully!")
		return
	}

	color.Red("\n✗ Deployment completed with %d error(s):", len(e.errors))
	fmt.Println()

	for i, err := range e.errors {
		color.Red("%d. [%s] %s", i+1, err.Code, err.Message)
		if err.Resource != "" {
			fmt.Printf("   Resource: %s\n", err.Resource)
		}
		if err.Cause != nil {
			fmt.Printf("   Details: %v\n", err.Cause)
		}
		fmt.Println()
	}
}

func (e *ErrorReporter) PrintDetailedError(err *errors.TerraformError) {
	color.Red("╭─ Error: %s", err.Message)
	color.Red("│ Code: %s", err.Code)

	if err.Resource != "" {
		color.Red("│ Resource: %s", err.Resource)
	}

	if err.Cause != nil {
		causeLines := strings.Split(err.Cause.Error(), "\n")
		for i, line := range causeLines {
			if i == 0 {
				color.Red("│ Caused by: %s", line)
			} else {
				color.Red("│            %s", line)
			}
		}
	}

	if len(err.Stack) > 0 {
		color.Red("│ Stack trace:")
		for _, frame := range err.Stack[:3] { // Показываем только 3 верхних кадра
			color.Red("│   %s", frame)
		}
	}

	color.Red("╰─")
}
