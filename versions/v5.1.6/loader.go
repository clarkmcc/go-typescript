package v5_1_6

import (
	_ "embed"

	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v5.1.6.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v5.1.6", Source)
}
