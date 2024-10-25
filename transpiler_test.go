package typescript

import (
	"context"
	"github.com/clarkmcc/go-typescript/versions"
	v4_2_3 "github.com/clarkmcc/go-typescript/versions/v4.2.3"
	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCompileVariousScripts(t *testing.T) {
	runtime := goja.New()
	registry := versions.NewRegistry()
	registry.Register("v4.2.3", v4_2_3.Source)

	t.Run("let", func(t *testing.T) {
		compiled, err := TranspileString("let a: number = 10;", WithCompileOptions(map[string]interface{}{
			"module": "none",
		}), WithVersion("v4.2.3"), WithRegistry(registry), WithRuntime(runtime))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", compiled)
	})

	t.Run("arrow function", func(t *testing.T) {
		compiled, err := TranspileString("((): number => 10)()", WithCompileOptions(map[string]interface{}{
			"module": "none",
		}), WithVersion("v4.2.3"), WithRegistry(registry), WithRuntime(runtime))
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
	registry := versions.NewRegistry()
	registry.Register("v4.2.3", v4_2_3.Source)
	output, err := Transpile(strings.NewReader("let a: number = 10;"), WithRegistry(registry), WithVersion("v4.2.3"))
	require.NoError(t, err)
	require.Equal(t, "var a = 10;", output)
}
