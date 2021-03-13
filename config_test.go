package typescript

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_Initialize(t *testing.T) {
	cfg := NewDefaultConfig()
	err := cfg.Initialize()
	require.NoError(t, err)
	_, err = cfg.Runtime.RunString(fmt.Sprintf("%v('not a valid base64 string')", cfg.decoderName))
	require.Error(t, err)
}
