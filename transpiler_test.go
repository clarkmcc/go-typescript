package typescript

import (
	"context"
	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCompileVariousScripts(t *testing.T) {
	runtime := goja.New()

	t.Run("let", func(t *testing.T) {
		compiled, err := TranspileString("let a: number = 10;", WithCompileOptions(map[string]interface{}{
			"module": "none",
		}), WithVersion("v4.2.3"), WithRuntime(runtime))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", compiled)
	})

	t.Run("arrow function", func(t *testing.T) {
		compiled, err := TranspileString("((): number => 10)()", WithCompileOptions(map[string]interface{}{
			"module": "none",
		}), WithVersion("v4.2.3"), WithRuntime(runtime))
		require.NoError(t, err)
		require.Equal(t, "(function () { return 10; })();", compiled)
	})
}

//func TestCompileErrors(t *testing.T) {
//	t.Run("bad syntax", func(t *testing.T) {
//		_, err := TranspileString("asdjaksdhkjasd")
//		require.Error(t, err)
//	})
//}

func TestCancelContext(t *testing.T) {
	runtime := goja.New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := TranspileCtx(ctx, strings.NewReader("let a: number = 10;"), WithRuntime(runtime))
	require.Error(t, err)
}

func TestBadConfig(t *testing.T) {
	_, err := TranspileString("let a: number = 10;", withFailOnInitialize())
	require.Error(t, err)
}

func TestTranspile(t *testing.T) {
	output, err := Transpile(strings.NewReader("let a: number = 10;"))
	require.NoError(t, err)
	require.Equal(t, "var a = 10;", output)
}
