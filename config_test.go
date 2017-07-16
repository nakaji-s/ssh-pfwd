package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddRule(t *testing.T) {
	config := Config{}

	r1 := Rule{Priority: 3}
	r2 := Rule{Priority: 4}
	r3 := Rule{Priority: 5}
	r4 := Rule{Priority: 4}
	r5 := Rule{Priority: 2}
	config.AddRule(r1)
	config.AddRule(r2)
	config.AddRule(r3)
	config.AddRule(r4)
	config.AddRule(r5)

	assert.Equal(t, len(config.Rules), 5)
	assert.Equal(t, config.Rules, []Rule{r5, r1, r2, r4, r3})
}
