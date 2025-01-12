package ui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProgressBar(t *testing.T) {
	max := 100
	pb := NewProgressBar(max)

	assert.NotNil(t, pb)
	assert.Equal(t, max, pb.max)
	assert.Equal(t, 0, pb.current)
	assert.NotNil(t, pb.TextView)
	assert.Equal(t, "Progress", pb.GetTitle())
}

func TestSetProgress(t *testing.T) {
	max := 100
	pb := NewProgressBar(max)

	progress := 50
	pb.SetProgress(progress)

	expectedText := fmt.Sprintf("%d%% [%s%s]", progress, string(repeat('#', progress/2)), string(repeat('-', 50-progress/2)))
	assert.Equal(t, progress, pb.current)
	assert.Equal(t, expectedText, pb.GetText(true))
}

func TestRepeat(t *testing.T) {
	char := '#'
	count := 5
	expected := []rune{'#', '#', '#', '#', '#'}

	result := repeat(char, count)
	assert.Equal(t, expected, result)
}
