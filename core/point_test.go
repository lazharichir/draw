package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint_Translate(t *testing.T) {
	p := Point{X: 1, Y: 2}
	expected := Point{X: 4, Y: 6}
	actual := p.Translate(3, 4)
	assert.Equal(t, expected, actual)
}

func TestPoint_Add(t *testing.T) {
	p1 := Point{X: 1, Y: 2}
	p2 := Point{X: 3, Y: 4}
	expected := Point{X: 4, Y: 6}
	actual := p1.Add(p2)
	assert.Equal(t, expected, actual)
}

func TestPoint_Subtract(t *testing.T) {
	p1 := Point{X: 4, Y: 6}
	p2 := Point{X: 3, Y: 4}
	expected := Point{X: 1, Y: 2}
	actual := p1.Subtract(p2)
	assert.Equal(t, expected, actual)
}

func TestPoint_Multiply(t *testing.T) {
	p := Point{X: 2, Y: 3}
	expected := Point{X: 4, Y: 6}
	actual := p.Multiply(2)
	assert.Equal(t, expected, actual)
}

func TestPoint_Divide(t *testing.T) {
	p := Point{X: 4, Y: 6}
	expected := Point{X: 2, Y: 3}
	actual := p.Divide(2)
	assert.Equal(t, expected, actual)
}

func TestPoint_Equals(t *testing.T) {
	p1 := Point{X: 1, Y: 2}
	p2 := Point{X: 1, Y: 2}
	assert.True(t, p1.Equals(p2))

	p3 := Point{X: 3, Y: 4}
	assert.False(t, p1.Equals(p3))
}

func TestPoint_IsOrigin(t *testing.T) {
	p1 := Point{X: 0, Y: 0}
	assert.True(t, p1.IsOrigin())

	p2 := Point{X: 1, Y: 2}
	assert.False(t, p2.IsOrigin())
}
