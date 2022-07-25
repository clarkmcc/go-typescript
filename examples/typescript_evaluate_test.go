package examples

import (
	_ "embed"
	"fmt"
	"github.com/clarkmcc/go-typescript"
	"strings"
)

//go:embed typescript-example.ts
var script2 string

func ExampleTypescriptEvaluate() {
	// Transpile the typescript and return evaluated result
	result, err := typescript.Evaluate(strings.NewReader(script2), typescript.WithTranspile())
	if err != nil {
		panic(err)
	}
	fmt.Println(result.String())
	// Output:Hello John Doe!
}
