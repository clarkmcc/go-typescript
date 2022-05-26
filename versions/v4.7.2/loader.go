package v4_7_2

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v4.7.2.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v4.7.2", Source)
}
