package v4_1_4

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v4.1.4.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v4.1.4", Source)
}
