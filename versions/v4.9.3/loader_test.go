package v4_9_3

import (
	"github.com/clarkmcc/go-typescript/versions"
	"testing"
)

func TestRegister(t *testing.T) {
	versions.TestSource(t, "v4.9.3", Source)
}
