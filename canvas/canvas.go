package canvas

import "fmt"

type Canvas struct {
	width  uint32
	height uint32
	grid   []uint64
}

func MakeCanvas(width, height uint32) Canvas {
	return Canvas{width, height, make([]uint64, width*height)}
}

func (c *Canvas) Width() uint32 {
	return c.width
}

func (c *Canvas) Height() uint32 {
	return c.height
}

func (c *Canvas) At(x, y uint32) Color {
	if !(x < c.Width()) {
		panic(fmt.Sprintf("x = %d is out of bounds, width is %d", x, c.Width()))
	} else if !(y < c.Height()) {
		panic(fmt.Sprintf("y = %d is out of bounds, height is %d", y, c.Height()))
	}

	return MakeColor(c.grid[c.toIndex(x, y)])
}

func (c *Canvas) WriteAt(x, y uint32, color Color) {
	if !(x < c.Width()) {
		panic(fmt.Sprintf("x = %d is out of bounds, width is %d", x, c.Width()))
	} else if !(y < c.Height()) {
		panic(fmt.Sprintf("y = %d is out of bounds, height is %d", y, c.Height()))
	}

	c.grid[c.toIndex(x, y)] = color.Value()
}

// Convert x and y coordinates to index within c's grid.
func (c *Canvas) toIndex(x, y uint32) int {
	return int((y * c.width) + x)
}
