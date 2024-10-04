package typescript

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarkmcc/go-typescript/versions"
	v4_9_3 "github.com/clarkmcc/go-typescript/versions/v4.9.3"
	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
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
	registry := versions.NewRegistry()
	registry.Register("v4.9.3", v4_9_3.Source)

	// This test hits a lot of things:
	//  #1 - We test that we can load the almond AMD module loader
	//  #2 - We test that we can load our own 'evaluate before' script that declares an AMD module
	//  #3 - We test that we can import the AMD module in a type script script and use a function from the module
	t.Run("evaluate with custom AMD module", func(t *testing.T) {
		script := "import { multiply } from 'myModule'; multiply(5, 5)"
		result, err := EvaluateCtx(context.Background(), strings.NewReader(script),
			WithAlmondModuleLoader(),
			WithTranspile(),
			WithTranspileOptions(WithRegistry(registry),
				WithVersion("v4.9.3")),
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

	t.Run("with goja ES6 support", func(t *testing.T) {
		// Template literals
		result, err := Evaluate(strings.NewReader("var a = 10; `${a}`"))
		require.NoError(t, err)
		require.Equal(t, "10", result.Export())

		// Arrow functions
		result, err = Evaluate(strings.NewReader("(() => 10)();"))
		require.NoError(t, err)
		require.Equal(t, int64(10), result.Export())
	})

	t.Run("custom runtime", func(t *testing.T) {
		runtime := goja.New()
		result, err := Evaluate(strings.NewReader("var a = 10; a;"),
			WithEvaluationRuntime(runtime))
		require.NoError(t, err)
		require.Equal(t, int64(10), result.ToInteger())
	})

	t.Run("script hook", func(t *testing.T) {
		script := "var a = 10;"
		_, err := Evaluate(strings.NewReader(script),
			WithScriptHook(func(s string) (string, error) {
				require.Equal(t, script, s)
				return script, nil
			}))
		require.NoError(t, err)
	})

	t.Run("pre-transpile hook", func(t *testing.T) {
		s1 := "let a: number = 10"
		_, err := Evaluate(strings.NewReader(s1),
			WithTranspile(),
			WithTranspileOptions(WithRegistry(registry),
				WithVersion("v4.9.3")),
			WithScriptPreTranspileHook(func(s2 string) (string, error) {
				assert.Equal(t, s1, s2)
				return s2, nil
			}))
		assert.NoError(t, err)
	})
}

var _ io.Reader = &failingReader{}

type failingReader struct{}

func (f *failingReader) Read(p []byte) (n int, err error) { return 0, fmt.Errorf("intentional error") }
