package examples

import (
	_ "embed"
	"fmt"
	"github.com/clarkmcc/go-typescript"
	"github.com/clarkmcc/go-typescript/versions"
	v4_9_3 "github.com/clarkmcc/go-typescript/versions/v4.9.3"
	"strings"
)

//go:embed javascript-example-module.js
var module string

func ExampleTypescriptAMDModule() {
	registry := versions.NewRegistry()
	registry.Register("v4.9.3", v4_9_3.Source)

	result, err := typescript.Evaluate(strings.NewReader(`import { multiply } from 'myModule'; multiply(5, 5)`),
		typescript.WithTranspile(),
		typescript.WithAlmondModuleLoader(),
		typescript.WithEvaluateBefore(strings.NewReader(module)),
		typescript.WithTranspileOptions(
			typescript.WithRegistry(registry),
			typescript.WithVersion("v4.9.3"),
		))
	if err != nil {
		panic(err)
	}
	fmt.Println(result.ToInteger())
	// Output:25
}
