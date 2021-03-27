package v4_2_3

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v4.2.3.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v4.2.3", Source)
}
