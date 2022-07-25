// Copyright (C) 2022 Print Tracker, LLC - All Rights Reserved
//
// Unauthorized copying of this file, via any medium is strictly prohibited
// as this source code is proprietary and confidential. Dissemination of this
// information or reproduction of this material is strictly forbidden unless
// prior written permission is obtained from Print Tracker, LLC.

package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomString(t *testing.T) {
	s1 := RandomString()
	assert.Len(t, s1, 24)
	s2 := RandomString()
	assert.NotEqual(t, s1, s2)
}
