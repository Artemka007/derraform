package errors

import (
    "fmt"
    "runtime"
    "strings"
)

type TerraformError struct {
    Code     string
    Message  string
    Resource string
    Cause    error
    Stack    []string
}

func (e *TerraformError) Error() string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("[%s] %s", e.Code, e.Message))
    
    if e.Resource != "" {
        sb.WriteString(fmt.Sprintf(" (resource: %s)", e.Resource))
    }
    
    if e.Cause != nil {
        sb.WriteString(fmt.Sprintf("\nCaused by: %v", e.Cause))
    }
    
    if len(e.Stack) > 0 {
        sb.WriteString("\nStack trace:")
        for _, frame := range e.Stack {
            sb.WriteString("\n  " + frame)
        }
    }
    
    return sb.String()
}

func NewError(code, message string) *TerraformError {
    return &TerraformError{
        Code:    code,
        Message: message,
        Stack:   captureStack(2), // Пропускаем 2 кадра
    }
}

func WrapError(err error, code, message string) *TerraformError {
    terraErr, ok := err.(*TerraformError)
    if ok {
        terraErr.Message = message + ": " + terraErr.Message
        return terraErr
    }
    
    return &TerraformError{
        Code:    code,
        Message: message,
        Cause:   err,
        Stack:   captureStack(2),
    }
}

func ResourceError(resource, message string, cause error) *TerraformError {
    return &TerraformError{
        Code:     "RESOURCE_ERROR",
        Message:  message,
        Resource: resource,
        Cause:    cause,
        Stack:    captureStack(2),
    }
}

func captureStack(skip int) []string {
    var stack []string
    for i := skip; i < 10; i++ { // Берем 10 кадров
        pc, file, line, ok := runtime.Caller(i)
        if !ok {
            break
        }
        
        fn := runtime.FuncForPC(pc)
        if fn != nil {
            stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
        }
    }
    return stack
}