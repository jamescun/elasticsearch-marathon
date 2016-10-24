package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv2Args(t *testing.T) {
	args := env2args(
		[]string{"PREFIX_FOO=bar", "PREFIX_FOO=baz", "PREFIX_FOO_BAR=foo"},
		"PREFIX_",
		nil,
	)

	assert.Equal(t, []string{"--foo=bar", "--foo=baz", "--foo.bar=foo"}, args)
}
