package v4_9_3

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v4.9.3.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v4.9.3", Source)
}
