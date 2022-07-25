package examples

import (
	_ "embed"
	"fmt"
	"github.com/clarkmcc/go-typescript"
	"strings"
)

//go:embed javascript-example-module.js
var module string

func ExampleTypescriptAMDModule() {
	result, err := typescript.Evaluate(strings.NewReader(`import { multiply } from 'myModule'; multiply(5, 5)`),
		typescript.WithTranspile(),
		typescript.WithAlmondModuleLoader(),
		typescript.WithEvaluateBefore(strings.NewReader(module)))
	if err != nil {
		panic(err)
	}
	fmt.Println(result.ToInteger())
	// Output:25
}
