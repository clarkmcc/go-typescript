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
func WithTranspile() EvaluateOptionFunc {
	return func(cfg *EvaluateConfig) {
		cfg.Transpile = true
	}
}

// WithTranspileOptions adds options to be passed to the transpiler if the transpiler is applicable
func WithTranspileOptions(opts ...TranspileOptionFunc) EvaluateOptionFunc {
	return func(cfg *EvaluateConfig) {
		cfg.TranspileOptions = append(cfg.TranspileOptions, opts...)
	}
}

// Evaluate calls EvaluateCtx using the default background context
func Evaluate(src io.Reader, opts ...EvaluateOptionFunc) (goja.Value, error) {
	return EvaluateCtx(context.Background(), src, opts...)
}

// EvaluateCtx evaluates the provided src using the specified options and returns the goja value result or an error.
func EvaluateCtx(ctx context.Context, src io.Reader, opts ...EvaluateOptionFunc) (result goja.Value, err error) {
	cfg := &EvaluateConfig{}
	cfg.ApplyDefaults()
	for _, fn := range opts {
		fn(cfg)
	}
	done := make(chan struct{})
	started := make(chan struct{})
	defer close(done)
	go func() {
		// Inform the parent go-routine that we've started, this prevents a race condition where the
		// runtime would beat the context cancellation in unit tests even though the context started
		// out in a 'cancelled' state.
		close(started)
		select {
		case <-ctx.Done():
			cfg.Runtime.Interrupt("context halt")
		case <-done:
			return
		}
	}()
	<-started
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
		opts := []TranspileOptionFunc{
			// We handle our own runtime with our own cancellation
			WithRuntime(cfg.Runtime),
			WithPreventCancellation(),
		}
		for _, opt := range cfg.TranspileOptions {
			opts = append(opts, opt)
		}
		script, err = TranspileCtx(ctx, strings.NewReader(script), opts...)
		if err != nil {
			return nil, fmt.Errorf("transpiling script: %w", err)
		}
	}
	result, err = cfg.Runtime.RunString(script)
	if err != nil {
		if strings.Contains(err.Error(), "context halt") {
			err = context.Canceled
		}
	}
	return
}
