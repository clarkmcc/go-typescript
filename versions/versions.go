package versions

import (
	_ "embed"
)

//go:embed v4.2.3.js
var version423 string

//go:embed v4.2.2.js
var version422 string

func RegisterDefaultsTo(registry *Registry) {
	registry.Register("v4.2.3", version423)
	registry.Register("v4.2.2", version422)
}

func init() {
	// By default, we register all typescripts sources that are packed into this library
	// into the default registry.
	RegisterDefaultsTo(DefaultRegistry)
}
