package util

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestLoadResourceSimple(t *testing.T) {
	assert := require.New(t)

	expected := `services:
- debian-console
- ubuntu-console
`
	expected = strings.TrimSpace(expected)

	b, e := LoadResource("https://raw.githubusercontent.com/rancherio/os-services/v0.3.4/index.yml", []string{})

	assert.Nil(e)
	assert.Equal(expected, strings.TrimSpace(string(b)))
}
