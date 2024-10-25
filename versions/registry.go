package versions

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type Registry interface {
	Register(tag string, source string)
	Get(tag string) (*goja.Program, error)
}

// CachingRegistry is a thread-safe registry for storing tagged versions of the typescript source code.
type CachingRegistry struct {
	lock     sync.Mutex
	versions map[string]string
	compiled map[string]*goja.Program
}

// Register registers the provided source to the specified tag in the registry.
func (r *CachingRegistry) Register(tag string, source string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.versions[tag] = source
	delete(r.compiled, tag)
}

// Get attempts to return the typescript source for the specified tag if it exists, otherwise
// it returns an error with a list of typescript versions that are supported by this registry.
func (r *CachingRegistry) Get(tag string) (*goja.Program, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	prg, ok := r.compiled[tag]
	if ok {
		return prg, nil
	}
	src, ok := r.versions[tag]
	if !ok {
		return nil, fmt.Errorf("unsupported version tag '%s', must be one of %v", tag, r.supportedVersionsLocked())
	}
	prg, err := goja.Compile("", src, true)
	if err != nil {
		return nil, fmt.Errorf("compiling registered source for tag '%s': %w", tag, err)
	}
	r.compiled[tag] = prg
	return prg, nil
}

// RegisteredVersions returns an unordered list of the versions that are registered in this registry
func (r *CachingRegistry) RegisteredVersions() (out []string) {
	for k, _ := range r.versions {
		out = append(out, k)
	}
	return
}

// supportedVersionsLocked returns a slice of supported version tags that are registered
// to this registry and can be accessed by calling Get. This function should only be called
// by a caller who has already acquired a lock on the registry.
func (r *CachingRegistry) supportedVersionsLocked() (out []string) {
	for k, _ := range r.versions {
		out = append(out, k)
	}
	return
}

// NewRegistry creates a new instances of a version registry
func NewRegistry() *CachingRegistry {
	return &CachingRegistry{
		versions: make(map[string]string),
		compiled: make(map[string]*goja.Program),
	}
}

// TestSource is a helper function for testing that versions of the Typescript compiler can
// properly be registered.
func TestSource(t *testing.T, version, source string) {
	r := NewRegistry()
	r.Register(version, source)
	_, err := r.Get(version)
	assert.NoErrorf(t, err, "failed to register %v", version)
}
