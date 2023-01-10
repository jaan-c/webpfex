package canvas

import "testing"

func TestResolution(t *testing.T) {
	var width uint32 = 16
	var height uint32 = 9
	canvas := MakeCanvas(width, height)

	if w := canvas.Width(); w != width {
		t.Errorf("Expecting Width to be %q, got %q", width, w)
	}
	if h := canvas.height; h != height {
		t.Errorf("Expecting Height to be %q, got %q", height, h)
	}
}

func TestWriteAndAt(t *testing.T) {
	canvas := MakeCanvas(2, 2)

	canvas.WriteAt(0, 0, MakeColor(1))
	canvas.WriteAt(1, 0, MakeColor(2))
	canvas.WriteAt(0, 1, MakeColor(3))
	canvas.WriteAt(1, 1, MakeColor(4))

	if p := canvas.At(0, 0).Value(); p != 1 {
		t.Errorf("Expecting 0,0 to be %q, got %q", 1, p)
	}
	if p := canvas.At(1, 0).Value(); p != 2 {
		t.Errorf("Expecting 0,0 to be %q, got %q", 2, p)
	}
	if p := canvas.At(0, 1).Value(); p != 3 {
		t.Errorf("Expecting 0,0 to be %q, got %q", 3, p)
	}
	if p := canvas.At(1, 1).Value(); p != 4 {
		t.Errorf("Expecting 0,0 to be %q, got %q", 4, p)
	}
}
