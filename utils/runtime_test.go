package utils

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReturnError(t *testing.T) {
	runtime := goja.New()
	ReturnError(runtime, fmt.Errorf("there's a problem"))
	_, err := runtime.RunString("var a = 10;")
	require.Error(t, err)
}

func TestErrorWrapper(t *testing.T) {
	runtime := goja.New()
	err := runtime.Set("a", ErrorWrapper(runtime, func(call goja.FunctionCall) (interface{}, error) {
		return 10, nil
	}))
	require.NoError(t, err)
	err = runtime.Set("b", ErrorWrapper(runtime, func(call goja.FunctionCall) (interface{}, error) {
		return nil, fmt.Errorf("there's a problem")
	}))
	require.NoError(t, err)
	err = runtime.Set("c", ErrorWrapper(runtime, func(call goja.FunctionCall) (interface{}, error) {
		return nil, nil
	}))
	require.NoError(t, err)

	t.Run("NoError", func(t *testing.T) {
		_, err = runtime.RunString("a()")
		require.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		_, err = runtime.RunString("b()")
		require.Error(t, err)
	})
	t.Run("Undefined", func(t *testing.T) {
		_, err = runtime.RunString("c()")
		require.NoError(t, err)
	})
}
