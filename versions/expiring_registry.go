package versions

import (
	"fmt"
	"github.com/dop251/goja"
	"sync"
	"time"
)

// ExpiringRegistry is a thread-safe registry for storing Typescript programs
// that are garbage collected after a certain amount of inactivity. Retrieving
// a program from the registry will reset its expiration time allowing the
// compiled program to stay cached for longer.
type ExpiringRegistry struct {
	lock     sync.Mutex
	versions map[string]string
	compiled map[string]entry

	// A struct is sent on this channel every time the registry is cleaned up.
	Freed chan struct{}

	ttl time.Duration
}

type entry struct {
	value *goja.Program
	exp   time.Time
}

func (r *ExpiringRegistry) Register(tag string, source string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.versions[tag] = source
}

func (r *ExpiringRegistry) Get(tag string) (*goja.Program, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	e, ok := r.compiled[tag]
	if ok && e.exp.After(time.Now()) {
		e.exp = time.Now().Add(r.ttl)
		r.compiled[tag] = e
		return e.value, nil
	}
	delete(r.compiled, tag)

	src, ok := r.versions[tag]
	if !ok {
		return nil, fmt.Errorf("unsupported version tag '%s', must be one of %v", tag, r.RegisteredVersions())
	}
	prg, err := goja.Compile("", src, true)
	if err != nil {
		return nil, fmt.Errorf("compiling registered source for tag '%s': %w", tag, err)
	}
	r.compiled[tag] = entry{value: prg, exp: time.Now().Add(r.ttl)}
	return prg, nil
}

func (r *ExpiringRegistry) RegisteredVersions() (out []string) {
	for k := range r.versions {
		out = append(out, k)
	}
	return
}

func NewExpiringRegistry(ttl time.Duration) *ExpiringRegistry {
	r := &ExpiringRegistry{
		versions: make(map[string]string),
		compiled: make(map[string]entry),
		Freed:    make(chan struct{}),
		ttl:      ttl,
	}

	go func() {
		for {
			time.Sleep(r.ttl)
			r.lock.Lock()
			var deleted bool
			for k, e := range r.compiled {
				if e.exp.Before(time.Now()) {
					deleted = true
					delete(r.compiled, k)
				}
			}
			r.lock.Unlock()

			// Once we've cleaned up, notify any waiting goroutines
			// that we're done. This is useful if callers want to
			// run a manual GC cycle after this.
			if deleted {
				select {
				case r.Freed <- struct{}{}:
				}
			}
		}
	}()

	return r
}
