package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

type ProgressBar struct {
	*tview.TextView
	max     int
	current int
}

func NewProgressBar(max int) *ProgressBar {
	pb := &ProgressBar{
		TextView: tview.NewTextView(),
		max:      max,
	}
	pb.SetBorder(true).SetTitle("Progress").SetTitleAlign(tview.AlignLeft)
	return pb
}

func (pb *ProgressBar) SetProgress(progress int) {
	pb.current = progress
	pb.SetText(fmt.Sprintf("%d%% [%s%s]", pb.current, string(repeat('#', pb.current/2)), string(repeat('-', 50-pb.current/2))))
}

func repeat(char rune, count int) []rune {
	result := make([]rune, count)
	for i := range result {
		result[i] = char
	}
	return result
}
