package typescript

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"io"
	"io/ioutil"
	"strings"
)

// Transpile transpiles the bytes read from reader using the provided config and options
func Transpile(reader io.Reader, cfg *Config, opts ...OptionFunc) (string, error) {
	return TranspileCtx(context.Background(), reader, cfg, opts...)
}

// TranspileString compiles the provided typescript string and returns the
func TranspileString(script string, cfg *Config, opts ...OptionFunc) (string, error) {
	return TranspileCtx(context.Background(), strings.NewReader(script), cfg, opts...)
}

// TranspileCtx compiles the bytes read from script using the provided context. Note that due to a limitation
// in goja, context cancellation only works while in JavaScript code, it does not interrupt native Go functions.
func TranspileCtx(ctx context.Context, script io.Reader, cfg *Config, opts ...OptionFunc) (string, error) {
	if cfg == nil {
		cfg = NewDefaultConfig()
	}
	if cfg.Runtime == nil {
		cfg.Runtime = goja.New()
	}
	for _, fn := range opts {
		fn(cfg)
	}
	// Handle context cancellation
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
	err := cfg.Initialize()
	if err != nil {
		return "", fmt.Errorf("initializing config: %w", err)
	}
	_, err = cfg.Runtime.RunProgram(cfg.TypescriptSource)
	if err != nil {
		return "", fmt.Errorf("running typescript compiler: %w", err)
	}
	optionBytes, err := json.Marshal(cfg.CompileOptions)
	if err != nil {
		return "", fmt.Errorf("marshalling compile options: %w", err)
	}
	scriptBytes, err := ioutil.ReadAll(script)
	if err != nil {
		return "", fmt.Errorf("reading script from reader: %w", err)
	}
	value, err := cfg.Runtime.RunString(fmt.Sprintf("ts.transpile(%s('%s'), %s, /*fileName*/ undefined, /*diagnostics*/ undefined, /*moduleName*/ \"myModule\")",
		cfg.decoderName, base64.StdEncoding.EncodeToString(scriptBytes), optionBytes))
	if err != nil {
		return "", fmt.Errorf("running compiler: %w", err)
	}
	return strings.TrimSuffix(value.String(), "\r\n"), nil
}
