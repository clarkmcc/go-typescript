package v4_1_2

import (
	"github.com/clarkmcc/go-typescript/versions"
	"testing"
)

func TestRegister(t *testing.T) {
	versions.TestSource(t, "v4.1.2", Source)
}
