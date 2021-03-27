package v3_9_9

import (
	_ "embed"
	"github.com/clarkmcc/go-typescript/versions"
)

//go:embed v3.9.9.js
var Source string

func init() {
	versions.DefaultRegistry.MustRegister("v3.9.9", Source)
}
