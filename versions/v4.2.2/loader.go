package v4_2_2

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v4.2.2.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v4.2.2", Source)
}
