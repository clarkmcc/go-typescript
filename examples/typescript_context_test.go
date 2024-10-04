package examples

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/clarkmcc/go-typescript"
	"github.com/clarkmcc/go-typescript/versions"
	v4_9_3 "github.com/clarkmcc/go-typescript/versions/v4.9.3"
	"strings"
)

//go:embed typescript-example.ts
var script3 string

func ExampleContext() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	registry := versions.NewRegistry()
	registry.Register("v4.9.3", v4_9_3.Source)

	_, err := typescript.TranspileCtx(ctx,
		strings.NewReader(script3),
		typescript.WithRegistry(registry),
		typescript.WithVersion("v4.9.3"))
	if err == nil {
		panic("expected error")
	}
	fmt.Println(err)
	// Output:running typescript compiler: context halt at <eval>:1:1(0)
}
