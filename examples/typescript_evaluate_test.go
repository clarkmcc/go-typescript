package examples

import (
	_ "embed"
	"fmt"
	"github.com/clarkmcc/go-typescript"
	"github.com/clarkmcc/go-typescript/versions"
	v4_9_3 "github.com/clarkmcc/go-typescript/versions/v4.9.3"
	"strings"
)

//go:embed typescript-example.ts
var script2 string

func ExampleTypescriptEvaluate() {
	registry := versions.NewRegistry()
	registry.Register("v4.9.3", v4_9_3.Source)

	// Transpile the typescript and return evaluated result
	result, err := typescript.Evaluate(strings.NewReader(script2), typescript.WithTranspile(), typescript.WithTranspileOptions(
		typescript.WithRegistry(registry),
		typescript.WithVersion("v4.9.3"),
	))
	if err != nil {
		panic(err)
	}
	fmt.Println(result.String())
	// Output:Hello John Doe!
}
