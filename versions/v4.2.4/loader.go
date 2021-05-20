package v4_2_4

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v4.2.4.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v4.2.4", Source)
}
