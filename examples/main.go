package main

import (
	_ "embed"
	typescript "github.com/clarkmcc/go-typescript"
	_ "github.com/clarkmcc/go-typescript/versions/v4.2.2"
	"log"
)

// This is a typescript script that we'll transpile to javascript
//go:embed script.ts
var script string

func main() {
	// The most basic implementation of the transpiler
	output, err := typescript.TranspileString(script, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(output)

	// You can specify a custom typescript version
	output, err = typescript.TranspileString(script, nil, typescript.WithVersion("v4.2.2"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(output)

	// Or provide some typescript compiler options using any of the options here https://www.typescriptlang.org/docs/handbook/compiler-options.html
	output, err = typescript.TranspileString(script, nil, typescript.WithCompileOptions(map[string]interface{}{
		"module": "none",
		"strict": true,
	}))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(output)
}
