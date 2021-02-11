package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenRegistry(t *testing.T) {
	registry, err := OpenRegistry()
	if err != nil {
		t.Error(err)
	}

	err = registry.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestRegistry_Remove(t *testing.T) {
	registry, err := OpenRegistry()
	if err != nil {
		t.Error(err)
	}

	registry.Remove("AppImageUpdate-x86_64-old.AppImage")
	_, ok := registry.Entries["AppImageUpdate-x86_64-old.AppImage"]
	assert.False(t, ok)

	err = registry.Close()
	if err != nil {
		t.Error(err)
	}
}
