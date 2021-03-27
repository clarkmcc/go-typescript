package versions

import (
	"fmt"
	"github.com/dop251/goja"
	"sync"
)

// DefaultRegistry is the default instance of the typescript tagged version registry.
var DefaultRegistry = NewRegistry()

// Registry is a thread-safe registry for storing tagged versions of the typescript source code.
type Registry struct {
	lock     sync.Mutex
	versions map[string]*goja.Program
}

// Register registers the provided source to the specified tag in the registry.
func (r *Registry) Register(tag string, source string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	program, err := goja.Compile("", source, true)
	if err != nil {
		return fmt.Errorf("compiling registered source for tag '%s': %w", tag, err)
	}
	r.versions[tag] = program
	return nil
}

// MustRegister calls Register and panics if we're unable to register the version.
func (r *Registry) MustRegister(tag string, source string) {
	err := r.Register(tag, source)
	if err != nil {
		panic(err)
	}
}

// Get attempts to return the typescript source for the specified tag if it exists, otherwise
// it returns an error with a list of typescript versions that are supported by this registry.
func (r *Registry) Get(tag string) (*goja.Program, error) {
	src, ok := r.versions[tag]
	if !ok {
		return nil, fmt.Errorf("unsupported version tag '%s', must be one of %v", tag, r.supportedVersionsLocked())
	}
	return src, nil
}

// MustGet calls Get with the specified tag, but panics if the tag cannot be found.
func (r *Registry) MustGet(tag string) *goja.Program {
	source, err := r.Get(tag)
	if err != nil {
		panic(err)
	}
	return source
}

// supportedVersionsLocked returns a slice of supported version tags that are registered
// to this registry and can be accessed by calling Get. This function should only be called
// by a caller who has already acquired a lock on the registry.
func (r *Registry) supportedVersionsLocked() (out []string) {
	for k, _ := range r.versions {
		out = append(out, k)
	}
	return
}

// NewRegistry creates a new instances of a version registry
func NewRegistry() *Registry {
	return &Registry{
		versions: map[string]*goja.Program{},
	}
}
