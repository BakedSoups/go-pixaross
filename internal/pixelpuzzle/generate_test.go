package pixelpuzzle

import (
	"image"
	"image/color"
	"testing"
)

func TestSolutionFromImageRectTreatsWhiteBackgroundAsEmpty(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	fill(img, color.RGBA{255, 255, 250, 255})
	img.SetRGBA(1, 1, color.RGBA{3, 2, 1, 255})

	got := solutionFromImageRect(img, img.Bounds(), 128, true)
	want := []string{
		"000",
		"010",
		"000",
	}
	assertRows(t, got, want)
}

func TestSolutionFromImageRectTreatsBlackBackgroundAsEmpty(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	fill(img, color.RGBA{2, 3, 4, 255})
	img.SetRGBA(1, 1, color.RGBA{250, 251, 252, 255})

	got := solutionFromImageRect(img, img.Bounds(), 128, true)
	want := []string{
		"000",
		"010",
		"000",
	}
	assertRows(t, got, want)
}

func fill(img *image.RGBA, c color.RGBA) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.SetRGBA(x, y, c)
		}
	}
}

func assertRows(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d rows, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("row %d = %q, want %q", i, got[i], want[i])
		}
	}
}
