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

func TestRegistry_Set(t *testing.T) {
	registry, err := OpenRegistry()
	if err != nil {
		t.Error(err)
	}

	_ = registry.Add("AppImageUpdate-x86_64-old.AppImage", "appimagehub:23942034")
	_, ok := registry.Entries["AppImageUpdate-x86_64-old.AppImage"]
	assert.True(t, ok)

	err = registry.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestRegistry_Get(t *testing.T) {
	registry, err := OpenRegistry()
	if err != nil {
		t.Error(err)
	}

	entry, _ := registry.Get("AppImageUpdate-x86_64-old.AppImage")
	assert.Equal(t, entry, registry.Entries["AppImageUpdate-x86_64-old.AppImage"])

	err = registry.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestRegistry_Unset(t *testing.T) {
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
