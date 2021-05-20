package typescript

import (
	"encoding/base64"
	"fmt"
	"github.com/clarkmcc/go-typescript/utils"
	"github.com/clarkmcc/go-typescript/versions"
	_ "github.com/clarkmcc/go-typescript/versions/v4.2.3"
	"github.com/dop251/goja"
)

// TranspileOptionFunc allows for easy chaining of pre-built config modifiers such as WithVersion.
type TranspileOptionFunc func(*Config)

// Config defines the behavior of the typescript compiler.
type Config struct {
	CompileOptions   map[string]interface{}
	TypescriptSource *goja.Program
	Runtime          *goja.Runtime

	// If a module is exported by the typescript compiler, this is the name the module will be called
	ModuleName string

	// Verbose enables built-in verbose logging for debugging purposes.
	Verbose bool

	// PreventCancellation indicates that the transpiler should not handle context cancellation. This
	// should be used when external runtimes are configured AND cancellation is handled by those runtimes.
	PreventCancellation bool

	// decoderName refers to a random generated string assigned to a function in the runtimes
	// global scope which is analogous to atob(), or a base64 decoding function. This function
	// is needed in the transpile process to ensure that we don't have any issues with string
	// interpolation errors when we pass our source code we want transpiled, into the typescript
	// transpiler. The reason we use a randomly generated string is to avoid the situation where
	// the transpiler caller provides their own runtime with a custom implementation of atob.
	decoderName string

	// Used only for testing to ensure that the compiler can handle config initialization failures
	failOnInitialize bool
}

func (c *Config) Initialize() error {
	if c.failOnInitialize {
		return fmt.Errorf("intentional error")
	}
	c.decoderName = utils.RandomString()
	return c.Runtime.Set(c.decoderName, utils.ErrorWrapper(c.Runtime, func(call goja.FunctionCall) (interface{}, error) {
		bs, err := base64.StdEncoding.DecodeString(call.Argument(0).String())
		if err != nil {
			return nil, err
		}
		return string(bs), nil
	}))
}

// NewDefaultConfig creates a new instance of the Config struct with default values and the latest
// typescript source code.s
func NewDefaultConfig() *Config {
	return &Config{
		Runtime:          goja.New(),
		CompileOptions:   nil,
		TypescriptSource: versions.DefaultRegistry.MustGet("v4.2.3"),
		ModuleName:       "default",
	}
}

// WithVersion loads the provided tagged typescript source from the default registry
func WithVersion(tag string) TranspileOptionFunc {
	return func(config *Config) {
		config.TypescriptSource = versions.DefaultRegistry.MustGet(tag)
	}
}

// WithCompileOptions sets the compile options that will be passed to the typescript compiler.
func WithCompileOptions(options map[string]interface{}) TranspileOptionFunc {
	return func(config *Config) {
		config.CompileOptions = options
	}
}

// WithRuntime allows you to over-ride the default runtime
func WithRuntime(runtime *goja.Runtime) TranspileOptionFunc {
	return func(config *Config) {
		config.Runtime = runtime
	}
}

// WithModuleName determines the module name applied to the typescript module if applicable. This is only needed to
// customize the module name if the typescript module mode is AMD or SystemJS.
func WithModuleName(name string) TranspileOptionFunc {
	return func(config *Config) {
		config.ModuleName = name
	}
}

// WithPreventCancellation prevents the transpiler runtime from handling its own context cancellation.
func WithPreventCancellation() TranspileOptionFunc {
	return func(config *Config) {
		config.PreventCancellation = true
	}
}

// withFailOnInitialize used to test a config initialization failure. This is not exported because
// it's used only for testing.
func withFailOnInitialize() TranspileOptionFunc {
	return func(config *Config) {
		config.failOnInitialize = true
	}
}
