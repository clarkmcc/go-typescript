package typescript

import (
	"context"
	"fmt"
	"github.com/clarkmcc/go-typescript/packages"
	_ "github.com/clarkmcc/go-typescript/versions/v4.2.3"
	"github.com/dop251/goja"
	"io"
	"io/ioutil"
	"strings"
)

type EvaluateOptionFunc func(cfg *EvaluateConfig)

type EvaluateConfig struct {
	// EvaluateBefore are sequentially evaluated in the Javascript runtime before evaluating the provided script.
	EvaluateBefore []io.Reader
	// Transpile indicates whether the script should be transpiled before its evaluated in the runtime.
	Transpile bool
	// TranspileOptions are options passed directly to the transpiler if applicable
	TranspileOptions []TranspileOptionFunc
	// Runtime is the goja runtime used for script execution. If not specified, it defaults to an empty runtime
	Runtime *goja.Runtime
}

// ApplyDefaults applies defaults to the configuration and is called automatically before the config is used
func (cfg *EvaluateConfig) ApplyDefaults() {
	if cfg.Runtime == nil {
		cfg.Runtime = goja.New()
	}
}

func (cfg *EvaluateConfig) HasEvaluateBefore() bool {
	return len(cfg.EvaluateBefore) > 0
}

// WithEvaluateBefore adds scripts that should be evaluated before evaluating the provided script. Each provided script
// is evaluated in the order that it's provided.
func WithEvaluateBefore(sources ...io.Reader) EvaluateOptionFunc {
	return func(cfg *EvaluateConfig) {
		cfg.EvaluateBefore = append(cfg.EvaluateBefore, sources...)
	}
}

// WithAlmondModuleLoader adds the almond module loader to the list of scripts that should be evaluated first
func WithAlmondModuleLoader() EvaluateOptionFunc {
	return WithEvaluateBefore(strings.NewReader(packages.Almond))
}

// WithTranspile indicates whether the provided script should be transpiled before it is evaluated. This does not
// mean that all the evaluate before's will be transpiled as well, only the src provided to EvaluateCtx will be transpiled
func WithTranspile(transpile bool) EvaluateOptionFunc {
	return func(cfg *EvaluateConfig) {
		cfg.Transpile = transpile
	}
}

// WithTranspileOptions adds options to be passed to the transpiler if the transpiler is applicable
func WithTranspileOptions(opts ...TranspileOptionFunc) EvaluateOptionFunc {
	return func(cfg *EvaluateConfig) {
		cfg.TranspileOptions = append(cfg.TranspileOptions, opts...)
	}
}

func EvaluateCtx(ctx context.Context, src io.Reader, opts ...EvaluateOptionFunc) (result goja.Value, err error) {
	cfg := &EvaluateConfig{}
	cfg.ApplyDefaults()
	for _, fn := range opts {
		fn(cfg)
	}
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			cfg.Runtime.Interrupt("halt")
		case <-done:
			return
		}
	}()
	if cfg.HasEvaluateBefore() {
		for _, s := range cfg.EvaluateBefore {
			b, err := ioutil.ReadAll(s)
			if err != nil {
				return nil, fmt.Errorf("reading evaluate befores: %w", err)
			}
			_, err = cfg.Runtime.RunString(string(b))
			if err != nil {
				return nil, fmt.Errorf("evaluating evaluate befores: %w", err)
			}
		}
	}
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("reading src: %w", err)
	}
	script := string(b)
	if cfg.Transpile {
		// This is needed in case the script being transpiled imports other modules. Check if it already exists in case
		// the caller has their own implementation and use of the global exports object.
		if cfg.Runtime.Get("exports") == nil {
			err = cfg.Runtime.Set("exports", cfg.Runtime.NewObject())
			if err != nil {
				return nil, fmt.Errorf("setting exports object: %w", err)
			}
		}
		opts := []TranspileOptionFunc{WithRuntime(cfg.Runtime)}
		for _, opt := range cfg.TranspileOptions {
			opts = append(opts, opt)
		}
		script, err = TranspileCtx(ctx, strings.NewReader(script), opts...)
		if err != nil {
			return nil, fmt.Errorf("transpiling script: %w", err)
		}
	}
	return cfg.Runtime.RunString(script)
}
