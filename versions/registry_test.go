package versions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_Get(t *testing.T) {
	// 5.x test
	t.Run("KnownTag", func(t *testing.T) {
		DefaultRegistry.MustRegister("v5.1.6", "")
		_, err := DefaultRegistry.Get("v5.1.6")
		require.NoError(t, err)
	})
	t.Run("KnownTag", func(t *testing.T) {
		DefaultRegistry.MustRegister("v4.2.3", "")
		_, err := DefaultRegistry.Get("v4.2.3")
		require.NoError(t, err)
	})
	t.Run("UnknownTag", func(t *testing.T) {
		_, err := DefaultRegistry.Get("abc")
		require.Error(t, err)
	})
	t.Run("MustGet", func(t *testing.T) {
		assert.NotPanics(t, func() {
			DefaultRegistry.MustGet("v4.2.3")
		})
	})
}

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()
	t.Run("ValidJavascript", func(t *testing.T) {
		err := r.Register("a", "var a = 10;")
		require.NoError(t, err)
	})
	t.Run("InvalidJavascript", func(t *testing.T) {
		err := r.Register("a", "type a struct{}")
		require.Error(t, err)
	})
	t.Run("RegisteredVersions", func(t *testing.T) {
		require.Len(t, r.RegisteredVersions(), 1)
	})
}

func TestRegistry_MustGet(t *testing.T) {
	r := NewRegistry()
	require.Panics(t, func() {
		r.MustGet("a")
	})
}
