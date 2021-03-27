package v3_8_3

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v3.8.3.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v3.8.3", Source)
}
