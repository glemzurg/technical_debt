package technical_debt

import (
	"fmt"
	"runtime/debug"
)

// errorWithStack is an error message with a stack dump.
type errorWithStack struct {
	message   string // The message of the error.
	stackDump string // The stack the moment of creation.
}

// Error prints the whole error.
func (e *errorWithStack) Error() string {
	return e.message + "\n\n" + e.stackDump
}

// Error creates an new error with stack information.
func Error(err error) (errWithStack error) {

	// No error means nothing to pass through.
	if err == nil {
		return nil
	}

	// The error may already have a stack.
	if _, ok := err.(*errorWithStack); ok {
		// Don't upset the stack.
		return err
	}

	// Bundle the error with a stack.
	return &errorWithStack{
		message:   err.Error(),
		stackDump: string(debug.Stack()),
	}
}

// Errorf creates a new error with stack information using fmt.Errorf parameters.
func Errorf(template string, params ...interface{}) (err error) {
	return Error(fmt.Errorf(template, params...))
}
