package examples

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/clarkmcc/go-typescript"
	"strings"
)

//go:embed typescript-example.ts
var script3 string

func ExampleContext() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := typescript.TranspileCtx(ctx, strings.NewReader(script3))
	if err == nil {
		panic("expected error")
	}
	fmt.Println(err)
	// Output:running typescript compiler: context halt at <eval>:1:1(0)
}
