# Goja Typescript Transpiler and Evaluator (with AMD module support)
This package provides a simple interface using [github.com/dop251/goja](github.com/dop251/goja) under the hood to allow you to transpile Typescript to Javascript in Go. In addition it provides an evaluator with a built-in AMD module loader which allows you to run Typescript code against a compiled typescript bundle. This package has no direct dependencies besides testing utilities and has a 95% test coverage rate.

Feel free to contribute. This package is fresh and may experience some changes before it's first tagged release.

## Transpiling Examples
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
output, err = typescript.TranspileString(script, typescript.WithCompileOptions(map[string]interface{}{
    "module": "none",
    "strict": true,
}))
```

### Custom Typescript Version
You can optionally specify which typescript version you want to compile using. These versions are based on the Git tags from the Typescript repository. If you're using a version that is supported in this package, you'll need to import the version package as a side-effect and will automatically be registered to the default registry.
```go
import _ "github.com/clarkmcc/go-typescript/versions/v4.2.2"

func main() {
    output, err := typescript.Transpile(reader, typescript.WithVersion("v4.2.2"))
}
```

### Custom Typescript Source
You may want to use a custom typescript version.

```go
func main() {
    output, err := typescript.TranspileString("let a:number = 10;", 
    	WithTypescriptSource("/* source code for typescript*/"))
}
```

## Evaluate Examples
### Basic Evaluation
You can evaluate pure Javascript code with:

```go
result, err := Evaluate(strings.NewReader("var a = 10;")) // returns 10;
```

### Transpile and Evaluate
Or you can transpile first:

```go
result, err := Evaluate(strings.NewReader("let a: number = 10;"), WithTranspile()) // returns 10;
```

### Run Script with AMD Modules
You can load in an AMD module bundle, then execute a Typescript script with access to the modules.

```go
// This is the module we're going to import
modules := strings.TrimSpace(`
    define("myModule", ["exports"], function (exports, core_1) {
        Object.defineProperty(exports, "__esModule", { value: true });
        exports.multiply = void 0;
        var multiply = function (a, b) { return a * b; };
        exports.multiply = multiply;
    });
`)

// This is the script we're going to transpile and evaluate
script := "import { multiply } from 'myModule'; multiply(5, 5)"

// Returns 25
result, err := EvaluateCtx(context.Background(), strings.NewReader(script),
    WithAlmondModuleLoader(),
    WithTranspile(),
    WithEvaluateBefore(strings.NewReader(amdModuleScript)))
```
