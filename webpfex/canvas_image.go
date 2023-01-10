package webpfex

import (
	"image"
	"image/color"
	"webpfex/canvas"
)

func ImageToCanvas(img image.Image) canvas.Canvas {
	// Bounds don't necessarily start at 0, yes it's hell!
	bounds := img.Bounds()
	width := uint32(bounds.Max.X - bounds.Min.X)
	height := uint32(bounds.Max.Y - bounds.Min.Y)
	cv := canvas.MakeCanvas(width, height)

	for fromY := bounds.Min.Y; fromY < bounds.Max.Y; fromY++ {
		for fromX := bounds.Min.X; fromX < bounds.Max.X; fromX++ {
			toX := uint32(fromX - bounds.Min.X)
			toY := uint32(fromY - bounds.Min.Y)

			r, g, b, a := img.At(fromX, fromY).RGBA()
			color := canvas.MakeColorRgba(uint16(r), uint16(g), uint16(b), uint16(a))

			cv.WriteAt(toX, toY, color)
		}
	}

	return cv
}

func CanvasToImage(cv canvas.Canvas) image.Image {
	ci := MakeCanvasImage(cv)
	return &ci
}

type CanvasImage struct {
	canvas canvas.Canvas
}

func MakeCanvasImage(cv canvas.Canvas) CanvasImage {
	return CanvasImage{cv}
}

func (ci *CanvasImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (ci *CanvasImage) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{int(ci.canvas.Width()), int(ci.canvas.Height())},
	}
}

func (ci *CanvasImage) At(x, y int) color.Color {
	return MakeCanvasImageColor(ci.canvas.At(uint32(x), uint32(y)))
}

type CanvasImageColor struct {
	color canvas.Color
}

func MakeCanvasImageColor(color canvas.Color) CanvasImageColor {
	return CanvasImageColor{color}
}

func (c CanvasImageColor) RGBA() (r, g, b, a uint32) {
	return uint32(c.color.R()), uint32(c.color.G()), uint32(c.color.B()),
		uint32(c.color.A())
}
