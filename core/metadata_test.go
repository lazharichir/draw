package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata_Get(t *testing.T) {
	m := Metadata{"foo": "bar"}
	assert.Equal(t, "bar", m.Get("foo"))
}

func TestMetadata_Set(t *testing.T) {
	m := Metadata{}
	assert.Nil(t, m.Get("foo"))
	m.Set("foo", "bar")
	assert.Equal(t, "bar", m.Get("foo"))
}

func TestMetadata_Delete(t *testing.T) {
	m := Metadata{"foo": "bar"}
	m.Delete("foo")
	assert.False(t, m.Has("foo"))
}

func TestMetadata_Has(t *testing.T) {
	m := Metadata{"foo": "bar"}
	assert.True(t, m.Has("foo"))
	assert.False(t, m.Has("baz"))
}

func TestMetadata_Keys(t *testing.T) {
	m := Metadata{"foo": "bar", "baz": "qux"}
	assert.ElementsMatch(t, []string{"foo", "baz"}, m.Keys())
}

func TestMetadata_MarshalJSON(t *testing.T) {
	m := Metadata{"foo": "bar", "baz": 42}
	expected := `{"baz":42,"foo":"bar"}`
	actual, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestMetadata_UnmarshalJSON(t *testing.T) {
	data := []byte(`{"baz":42,"foo":"bar"}`)
	expected := Metadata{"foo": "bar", "baz": float64(42)}
	var actual Metadata
	err := json.Unmarshal(data, &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestMetadata_Value(t *testing.T) {
	m := Metadata{"foo": "bar", "baz": 42}
	expected, err := json.Marshal(m)
	assert.NoError(t, err)
	actual, err := m.Value()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestMetadata_Scan(t *testing.T) {
	data := []byte(`{"baz":42,"foo":"bar"}`)
	expected := Metadata{"foo": "bar", "baz": float64(42)}
	var actual Metadata
	err := actual.Scan(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
