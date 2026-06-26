package pixelpuzzle

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"

	_ "golang.org/x/image/webp"
)

type PuzzleJSON struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Width          int        `json:"width"`
	Height         int        `json:"height"`
	Solution       []string   `json:"solution"`
	SkeletonPixels [][]string `json:"skeletonPixels,omitempty"`
	RevealPixels   [][]string `json:"revealPixels,omitempty"`
}

type SpriteSheetOptions struct {
	ID             string
	Title          string
	Source         string
	Out            string
	TileSize       int
	AlphaThreshold uint32
	UseBackground  bool
}

func GenerateSpriteSheet(opts SpriteSheetOptions) (PuzzleJSON, error) {
	if opts.ID == "" || opts.Title == "" || opts.Source == "" || opts.Out == "" || opts.TileSize <= 0 {
		return PuzzleJSON{}, fmt.Errorf("id, title, source, out, and tile size are required")
	}

	img, err := decodeImage(opts.Source)
	if err != nil {
		return PuzzleJSON{}, err
	}
	bounds := img.Bounds()
	tileSize, err := sheetTileSize(bounds, opts.TileSize)
	if err != nil {
		return PuzzleJSON{}, fmt.Errorf("%s is %dx%d: %w", opts.Source, bounds.Dx(), bounds.Dy(), err)
	}

	beforeRect := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+tileSize, bounds.Min.Y+tileSize)
	afterRect := image.Rect(bounds.Min.X+tileSize, bounds.Min.Y, bounds.Min.X+tileSize*2, bounds.Min.Y+tileSize)
	solution := solutionFromImageRect(img, beforeRect, opts.AlphaThreshold, opts.UseBackground)

	puzzle := PuzzleJSON{
		ID:             opts.ID,
		Title:          opts.Title,
		Width:          tileSize,
		Height:         tileSize,
		Solution:       solution,
		SkeletonPixels: pixelRows(img, beforeRect),
		RevealPixels:   pixelRows(img, afterRect),
	}
	if err := writePuzzleJSON(opts.Out, puzzle); err != nil {
		return PuzzleJSON{}, err
	}
	return puzzle, nil
}

func decodeImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func solutionFromImageRect(img image.Image, r image.Rectangle, alphaThreshold uint32, useBackground bool) []string {
	hasTransparency := imageRectHasTransparency(img, r, alphaThreshold)
	bg := detectBackground(img, r)

	rows := make([]string, 0, r.Dy())
	for y := r.Min.Y; y < r.Max.Y; y++ {
		row := make([]byte, 0, r.Dx())
		for x := r.Min.X; x < r.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			filled := a > alphaThreshold*257
			if !hasTransparency && useBackground && bg.matches(r, g, b) {
				filled = false
			}
			if filled {
				row = append(row, '1')
			} else {
				row = append(row, '0')
			}
		}
		rows = append(rows, string(row))
	}
	return rows
}

type backgroundTone int

const (
	backgroundColorTone backgroundTone = iota
	backgroundWhiteTone
	backgroundBlackTone
)

type background struct {
	red   uint32
	green uint32
	blue  uint32
	tone  backgroundTone
}

func detectBackground(img image.Image, r image.Rectangle) background {
	whitePixels, blackPixels := countMonochromePixels(img, r)
	totalPixels := r.Dx() * r.Dy()
	if (whitePixels+blackPixels)*2 >= totalPixels {
		if whitePixels >= blackPixels {
			return background{tone: backgroundWhiteTone}
		}
		return background{tone: backgroundBlackTone}
	}

	red, green, blue := backgroundColor(img, r)
	return background{red: red, green: green, blue: blue, tone: backgroundColorTone}
}

func countMonochromePixels(img image.Image, r image.Rectangle) (int, int) {
	whitePixels := 0
	blackPixels := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			red, green, blue, _ := img.At(x, y).RGBA()
			switch {
			case isWhiteish(red, green, blue):
				whitePixels++
			case isBlackish(red, green, blue):
				blackPixels++
			}
		}
	}
	return whitePixels, blackPixels
}

func (bg background) matches(r, g, b uint32) bool {
	switch bg.tone {
	case backgroundWhiteTone:
		return isWhiteish(r, g, b)
	case backgroundBlackTone:
		return isBlackish(r, g, b)
	default:
		return closeColor(r, g, b, bg.red, bg.green, bg.blue)
	}
}

func sheetTileSize(bounds image.Rectangle, requested int) (int, error) {
	if bounds.Dx() >= requested*2 && bounds.Dy() >= requested {
		return requested, nil
	}
	if bounds.Dx() >= bounds.Dy()*2 {
		return bounds.Dy(), nil
	}
	return 0, fmt.Errorf("expected at least %dx%d or a two-panel sheet", requested*2, requested)
}

func backgroundColor(img image.Image, r image.Rectangle) (uint32, uint32, uint32) {
	var redSum, greenSum, blueSum, count uint64
	add := func(x, y int) {
		red, green, blue, _ := img.At(x, y).RGBA()
		redSum += uint64(red)
		greenSum += uint64(green)
		blueSum += uint64(blue)
		count++
	}
	for x := r.Min.X; x < r.Max.X; x++ {
		add(x, r.Min.Y)
		add(x, r.Max.Y-1)
	}
	for y := r.Min.Y + 1; y < r.Max.Y-1; y++ {
		add(r.Min.X, y)
		add(r.Max.X-1, y)
	}
	if count == 0 {
		red, green, blue, _ := img.At(r.Min.X, r.Min.Y).RGBA()
		return red, green, blue
	}
	return uint32(redSum / count), uint32(greenSum / count), uint32(blueSum / count)
}

func pixelRows(img image.Image, r image.Rectangle) [][]string {
	rows := make([][]string, 0, r.Dy())
	for y := r.Min.Y; y < r.Max.Y; y++ {
		row := make([]string, 0, r.Dx())
		for x := r.Min.X; x < r.Max.X; x++ {
			red, green, blue, alpha := img.At(x, y).RGBA()
			if alpha == 0 {
				row = append(row, "")
				continue
			}
			row = append(row, fmt.Sprintf("#%02X%02X%02X%02X", uint8(red>>8), uint8(green>>8), uint8(blue>>8), uint8(alpha>>8)))
		}
		rows = append(rows, row)
	}
	return rows
}

func writePuzzleJSON(out string, puzzle PuzzleJSON) error {
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(puzzle, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(out, "puzzle.json"), data, 0o644)
}

func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func imageRectHasTransparency(img image.Image, r image.Rectangle, alphaThreshold uint32) bool {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a <= alphaThreshold*257 {
				return true
			}
		}
	}
	return false
}

func closeColor(r, g, b, wantR, wantG, wantB uint32) bool {
	const tolerance = 32 * 257
	return diff(r, wantR) <= tolerance && diff(g, wantG) <= tolerance && diff(b, wantB) <= tolerance
}

func isWhiteish(r, g, b uint32) bool {
	const whiteThreshold = 224 * 257
	return r >= whiteThreshold && g >= whiteThreshold && b >= whiteThreshold
}

func isBlackish(r, g, b uint32) bool {
	const blackThreshold = 32 * 257
	return r <= blackThreshold && g <= blackThreshold && b <= blackThreshold
}

func diff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}
