package v3_9_9

import (
	"github.com/clarkmcc/go-typescript/versions"
	"testing"
)

func TestRegister(t *testing.T) {
	versions.TestSource(t, "v3.9.9", Source)
}
