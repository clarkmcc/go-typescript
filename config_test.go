package typescript

import (
	"fmt"
	"github.com/clarkmcc/go-typescript/versions"
	_ "github.com/clarkmcc/go-typescript/versions/v3.8.3"
	_ "github.com/clarkmcc/go-typescript/versions/v3.9.9"
	_ "github.com/clarkmcc/go-typescript/versions/v4.1.2"
	_ "github.com/clarkmcc/go-typescript/versions/v4.1.3"
	_ "github.com/clarkmcc/go-typescript/versions/v4.1.4"
	_ "github.com/clarkmcc/go-typescript/versions/v4.1.5"
	_ "github.com/clarkmcc/go-typescript/versions/v4.2.2"
	v423 "github.com/clarkmcc/go-typescript/versions/v4.2.3"
	_ "github.com/clarkmcc/go-typescript/versions/v4.2.4"
	_ "github.com/clarkmcc/go-typescript/versions/v4.7.2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_Initialize(t *testing.T) {
	cfg := NewDefaultConfig()
	err := cfg.Initialize()
	require.NoError(t, err)
	_, err = cfg.Runtime.RunString(fmt.Sprintf("%v('not a valid base64 string')", cfg.decoderName))
	require.Error(t, err)
}

func TestVersionLoading(t *testing.T) {
	t.Run("v3.8.3", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v3.8.3"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
	t.Run("v3.9.9", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v3.9.9"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
	t.Run("v4.1.2", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v4.1.2"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
	t.Run("v4.1.3", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v4.1.3"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
	t.Run("v4.1.4", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v4.1.4"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
	t.Run("v4.1.5", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v4.1.5"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
	t.Run("v4.2.2", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v4.2.2"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
	t.Run("v4.2.3", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v4.2.3"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
	t.Run("v4.2.4", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v4.2.4"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
	t.Run("v4.7.2", func(t *testing.T) {
		output, err := TranspileString("let a: number = 10;", WithVersion("v4.7.2"))
		require.NoError(t, err)
		require.Equal(t, "var a = 10;", output)
	})
}

func TestCustomRegistry(t *testing.T) {
	registry := versions.NewRegistry()
	registry.MustRegister("v4.2.3", v423.Source)

	output, err := TranspileString("let a: number = 10;", func(config *Config) {
		config.TypescriptSource = registry.MustGet("v4.2.3")
	})
	require.NoError(t, err)
	require.Equal(t, "var a = 10;", output)
}

func TestWithModuleName(t *testing.T) {
	output, err := TranspileString("let a: number = 10;",
		WithModuleName("myModuleName"),
		WithCompileOptions(map[string]interface{}{
			"module": "amd",
		}))
	require.NoError(t, err)
	require.Contains(t, output, "define(\"myModuleName\"")
}

func TestWithTypescriptSource(t *testing.T) {
	output, err := TranspileString("let a: number = 10;",
		WithTypescriptSource(v423.Source))
	require.NoError(t, err)
	require.Equal(t, "var a = 10;", output)
}
