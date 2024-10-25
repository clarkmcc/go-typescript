package versions

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()
	t.Run("KnownTag", func(t *testing.T) {
		r.Register("v4.2.3", "")
		_, err := r.Get("v4.2.3")
		require.NoError(t, err)
	})
	t.Run("UnknownTag", func(t *testing.T) {
		_, err := r.Get("abc")
		require.Error(t, err)
	})
}

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()
	t.Run("ValidJavascript", func(t *testing.T) {
		r.Register("a", "var a = 10;")
		_, err := r.Get("a")
		require.NoError(t, err)
	})
	t.Run("InvalidJavascript", func(t *testing.T) {
		r.Register("a", "type a struct{}")
		prg, err := r.Get("a")
		require.Nil(t, prg)
		require.Error(t, err)
	})
	t.Run("RegisteredVersions", func(t *testing.T) {
		require.Len(t, r.RegisteredVersions(), 1)
	})
}
