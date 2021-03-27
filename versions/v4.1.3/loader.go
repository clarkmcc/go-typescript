package v4_1_3

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v4.1.3.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v4.1.3", Source)
}
