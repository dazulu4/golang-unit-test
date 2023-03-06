package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddSuccess(t *testing.T) {
	result := Add(20, 2)
	expected := 22

	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestAddSuccessRequire(t *testing.T) {
	c := require.New(t)
	result := Add(20, 2)
	expected := 22

	c.Equal(expected, result)
}
