package webpfex

import "webpfex/canvas"

// Clear canvas with color.
func ClearCanvas(canvas *canvas.Canvas, color canvas.Color) {
	for y := uint32(0); y < canvas.Height(); y++ {
		for x := uint32(0); x < canvas.Width(); x++ {
			canvas.WriteAt(x, y, color)
		}
	}
}

// Overlay canvas with another by replacing pixels.
func OverlayCanvas(canvas *canvas.Canvas, with *canvas.Canvas, xOffset, yOffset uint32) {
	for y := uint32(0); y < with.Height(); y++ {
		for x := uint32(0); x < with.Width(); x++ {
			canvas.WriteAt(x+xOffset, y+yOffset, with.At(x, y))
		}
	}
}

// Overlay canvas with another by blending the overlay's canvas with the original.
func OverlayBlendCanvas(canvas *canvas.Canvas, with *canvas.Canvas, xOffset, yOffset uint32) {
	for y := uint32(0); y < with.Height(); y++ {
		for x := uint32(0); x < with.Width(); x++ {
			dstX := x + xOffset
			dstY := y + yOffset

			overlay := with.At(x, y)
			color := canvas.At(dstX, dstY)
			blended := OverlayColor(color, overlay)

			canvas.WriteAt(dstX, dstY, blended)
		}
	}
}

func OverlayColor(color canvas.Color, with canvas.Color) canvas.Color {
	wA := float64(with.A()) / 0xFFFF
	bR := uint16((float64(with.R()) * wA) + (float64(color.R()) * (1 - wA)))
	bG := uint16((float64(with.G()) * wA) + (float64(color.G()) * (1 - wA)))
	bB := uint16((float64(with.B()) * wA) + (float64(color.B()) * (1 - wA)))

	return canvas.MakeColorRgba(bR, bG, bB, 0xFFFF)
}
