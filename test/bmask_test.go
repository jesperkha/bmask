package test

import (
	"testing"

	"github.com/jesperkha/bmask"
)

func TestBorderedMask(t *testing.T) {
	err := bmask.Draw("test.png", "test_bg2.png", 20, 0)
	if err != nil {
		t.Error(err)
	}
}
