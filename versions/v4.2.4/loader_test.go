package v4_2_4

import (
	"github.com/clarkmcc/go-typescript/versions"
	"testing"
)

func TestRegister(t *testing.T) {
	versions.TestSource(t, "v4.2.4", Source)
}
