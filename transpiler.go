package typescript

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// Transpile transpiles the bytes read from reader using the provided config and options
func Transpile(reader io.Reader, opts ...TranspileOptionFunc) (string, error) {
	return TranspileCtx(context.Background(), reader, opts...)
}

// TranspileString compiles the provided typescript string and returns the
func TranspileString(script string, opts ...TranspileOptionFunc) (string, error) {
	return TranspileCtx(context.Background(), strings.NewReader(script), opts...)
}

// TranspileCtx compiles the bytes read from script using the provided context. Note that due to a limitation
// in goja, context cancellation only works while in JavaScript code, it does not interrupt native Go functions.
func TranspileCtx(ctx context.Context, script io.Reader, opts ...TranspileOptionFunc) (string, error) {
	cfg := NewDefaultConfig()
	for _, fn := range opts {
		fn(cfg)
	}
	// Handle context cancellation
	if !cfg.PreventCancellation {
		done := startInterruptable(ctx, cfg.Runtime)
		defer close(done)
	}
	err := cfg.Initialize()
	if err != nil {
		return "", fmt.Errorf("initializing config: %w", err)
	}
	src, err := cfg.Registry.Get(cfg.TypescriptVersion)
	if err != nil {
		return "", fmt.Errorf("getting typescript source: %w", err)
	}
	_, err = cfg.Runtime.RunProgram(src)
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
	s := fmt.Sprintf("ts.transpile(%s('%s'), %s, /*fileName*/ undefined, /*diagnostics*/ undefined, /*moduleName*/ \"%s\")",
		cfg.decoderName, base64.StdEncoding.EncodeToString(scriptBytes), optionBytes, cfg.ModuleName)
	if cfg.Verbose {
		log.Println(s)
	}
	value, err := cfg.Runtime.RunString(s)
	if err != nil {
		return "", fmt.Errorf("running compiler: %w", err)
	}
	return strings.TrimSuffix(value.String(), "\r\n"), nil
}
