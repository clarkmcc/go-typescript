package utils

import "github.com/dop251/goja"

// Provides some helper utilities for working with the goja runtime

type FunctionWithError func(call goja.FunctionCall) (interface{}, error)

// ReturnError throws an interrupt error in the middle of a runtime execution
func ReturnError(runtime *goja.Runtime, err error) goja.Value {
	runtime.Interrupt(err)
	return goja.Undefined()
}

// ErrorWrapper wraps goja functions and provides runtime-based error handling.
func ErrorWrapper(runtime *goja.Runtime, in FunctionWithError) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		value, err := in(call)
		if err != nil {
			return ReturnError(runtime, err)
		}
		if value == nil {
			return goja.Undefined()
		}
		return runtime.ToValue(value)
	}
}
