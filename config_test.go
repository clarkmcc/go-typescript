package typescript

import (
	"fmt"
	"github.com/clarkmcc/go-typescript/versions"
	v3_8_3 "github.com/clarkmcc/go-typescript/versions/v3.8.3"
	v3_9_9 "github.com/clarkmcc/go-typescript/versions/v3.9.9"
	v4_1_2 "github.com/clarkmcc/go-typescript/versions/v4.1.2"
	v4_1_3 "github.com/clarkmcc/go-typescript/versions/v4.1.3"
	v4_1_4 "github.com/clarkmcc/go-typescript/versions/v4.1.4"
	v4_1_5 "github.com/clarkmcc/go-typescript/versions/v4.1.5"
	v4_2_2 "github.com/clarkmcc/go-typescript/versions/v4.2.2"
	v4_2_3 "github.com/clarkmcc/go-typescript/versions/v4.2.3"
	v4_2_4 "github.com/clarkmcc/go-typescript/versions/v4.2.4"
	v4_7_2 "github.com/clarkmcc/go-typescript/versions/v4.7.2"
	v4_9_3 "github.com/clarkmcc/go-typescript/versions/v4.9.3"
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
	registry := versions.NewRegistry()

	sources := map[string]string{
		"v3.8.3": v3_8_3.Source,
		"v3.9.9": v3_9_9.Source,
		"v4.1.2": v4_1_2.Source,
		"v4.1.3": v4_1_3.Source,
		"v4.1.4": v4_1_4.Source,
		"v4.1.5": v4_1_5.Source,
		"v4.2.2": v4_2_2.Source,
		"v4.2.3": v4_2_3.Source,
		"v4.2.4": v4_2_4.Source,
		"v4.7.2": v4_7_2.Source,
		"v4.9.3": v4_9_3.Source,
	}

	for tag, source := range sources {
		registry.Register(tag, source)
	}

	for tag, _ := range sources {
		t.Run(tag, func(t *testing.T) {
			output, err := TranspileString("let a: number = 10;", WithRegistry(registry), WithVersion(tag))
			require.NoError(t, err)
			require.Equal(t, "var a = 10;", output)
		})
	}
}

func TestWithModuleName(t *testing.T) {
	registry := versions.NewRegistry()
	registry.Register("v4.9.3", v4_9_3.Source)
	output, err := TranspileString("let a: number = 10;",
		WithModuleName("myModuleName"),
		WithRegistry(registry),
		WithVersion("v4.9.3"),
		WithCompileOptions(map[string]interface{}{
			"module": "amd",
		}))
	require.NoError(t, err)
	require.Contains(t, output, "define(\"myModuleName\"")
}
