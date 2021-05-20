package typescript

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

var (
	amdModuleScript = strings.TrimSpace(`
		define("myModule", ["exports"], function (exports, core_1) {
			Object.defineProperty(exports, "__esModule", { value: true });
			exports.multiply = void 0;
			var multiply = function (a, b) { return a * b; };
			exports.multiply = multiply;
		});
	`)
)

func TestEvaluateCtx(t *testing.T) {
	// This test hits a lot of things:
	//  #1 - We test that we can load the almond AMD module loader
	//  #2 - We test that we can load our own 'evaluate before' script that declares an AMD module
	//  #3 - We test that we can import the AMD module in a type script script and use a function from the module
	t.Run("evaluate with custom AMD module", func(t *testing.T) {
		script := "import { multiply } from 'myModule'; multiply(5, 5)"
		result, err := EvaluateCtx(context.Background(), strings.NewReader(script),
			WithAlmondModuleLoader(),
			WithTranspile(),
			WithEvaluateBefore(strings.NewReader(amdModuleScript)),
			WithTranspileOptions(func(config *Config) {
				config.Verbose = true
			}))
		require.NoError(t, err)
		require.Equal(t, int64(25), result.ToInteger())
	})

	// Ensures the context cancellation works correctly with the goja runtime
	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err := EvaluateCtx(ctx, strings.NewReader("var a = 10;"))
		fmt.Println(err)
		require.True(t, errors.Is(err, context.Canceled))
		require.Nil(t, result)
	})

	// A syntax error in the evaluate befores should return an error
	t.Run("evaluate 'evaluate before' error", func(t *testing.T) {
		_, err := Evaluate(strings.NewReader("var a = 10;"),
			WithEvaluateBefore(strings.NewReader("let a: number = 10;")))
		require.Error(t, err)
		require.Contains(t, err.Error(), "evaluating evaluate befores: SyntaxError")
	})

	t.Run("unreadable 'evaluate before'", func(t *testing.T) {
		_, err := Evaluate(strings.NewReader("var a = 10;"),
			WithEvaluateBefore(&failingReader{}))
		require.Error(t, err)
		require.Contains(t, err.Error(), "reading evaluate befores")
	})
}

var _ io.Reader = &failingReader{}

type failingReader struct{}

func (f *failingReader) Read(p []byte) (n int, err error) { return 0, fmt.Errorf("intentional error") }
