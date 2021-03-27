# Goja Typescript Transpiler
This package provides a simple interface using [github.com/dop251/goja](github.com/dop251/goja) under the hood to allow you to transpile Typescript to Javascript in Go. This package has no direct dependencies besides testing utilities and has a 95% test coverage rate.

Feel free to contribute. This package is fresh and may experience some changes before it's first tagged release.

## Example
For more examples, see the `examples/` directory of this repository
### Transpile Strings
```go
output, err := typescript.TranspileString("let a: number = 10;", nil)
// output: var a = 10;
```

### Transpile Reader
```go
output, err := typescript.Transpile(reader, nil)
```

### Custom Typescript Compile Options
You can optionally specify alternative compiler options that are used by Typescript. Any of the options [https://www.typescriptlang.org/docs/handbook/compiler-options.html](https://www.typescriptlang.org/docs/handbook/compiler-options.html) can be added.
```go
output, err = typescript.TranspileString(script, nil, typescript.WithCompileOptions(map[string]interface{}{
    "module": "none",
    "strict": true,
}))
```

### Custom Typescript Version
#### Default Registry
You can optionally specify which typescript version you want to compile using. These versions are based on the Git tags from the Typescript repository. If you're using a version that is supported in this package, you'll need to import the version package as a side-effect and will automatically be registered to the default registry.
```go
import _ "github.com/clarkmcc/go-typescript/versions/v4.2.2"

func main() {
    output, err := typescript.Transpile(reader, nil, typescript.WithVersion("v4.2.2"))
}
```

#### Custom Registry
You may want to use a custom version registry rather than the default registry.

```go
import version "github.com/clarkmcc/go-typescript/versions/v4.2.2"

func main() {
    registry := versions.NewRegistry()
    registry.MustRegister("v4.2.3", version.Source)
    
    output, err := typescript.TranspileString("let a:number = 10;", &typescript.Config{
        TypescriptSource: program,
    })
}
```

#### Custom Version
Need a different typescript version than the tags we support in this repo? No problem, you can load your own:

```go
program, err := goja.Compile("typescript", "<typescript source code here>", true)
output, err := typescript.Transpile(reader, &typescript.Config{
    CompileOptions:   map[string]interface{}{},
    TypescriptSource: program,
    Runtime:          goja.New(),
})
```
