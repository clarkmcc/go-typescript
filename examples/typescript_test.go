package examples

import (
	_ "embed"
	"fmt"
	"github.com/clarkmcc/go-typescript"
	"strings"
)

//go:embed typescript-example.ts
var script1 string

var expected = `
var Person = /** @class */ (function () {
    function Person(name) {
        this.name = name;
    }
    Person.prototype.greet = function () {
        return "Hello ".concat(this.name, "!");
    };
    return Person;
}());
var me = new Person("John Doe");
me.greet();`

func ExampleTranspile() {
	// Only transpile the typescript and return transpiled Javascript, don't evaluate
	transpiled, err := typescript.TranspileString(script1)
	if err != nil {
		panic(err)
	}
	if transpiled != expected {
		panic("unexpected transpile result")
	}
	fmt.Println(strings.TrimSpace(transpiled))
}
