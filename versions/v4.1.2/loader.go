package v4_1_2

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v4.1.2.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v4.1.2", Source)
}
