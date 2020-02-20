// +build go1.14

package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	as := assert.New(t)

	// args := []string{"hello", "-who=tsingson"}
	args := []string{"./program", "tsingson"}
	var stdout bytes.Buffer

	err := run(args, &stdout)
	as.NoError(err)

	out := stdout.String()
	as.True(strings.Contains(out, "Hello, tsingson"))
}
