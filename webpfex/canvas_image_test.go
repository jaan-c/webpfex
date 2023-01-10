package webpfex

import (
	"testing"
	"webpfex/canvas"
)

func TestMakeImageCanvasColor(t *testing.T) {
	color := canvas.MakeColorRgba(200, 150, 100, 50)
	cr, cg, cb, ca := color.Rgba()
	ir, ig, ib, ia := MakeCanvasImageColor(color).RGBA()

	if uint32(cr) != ir {
		t.Errorf("%d != %d", cr, ir)
	}
	if uint32(cg) != ig {
		t.Errorf("%d != %d", cg, ig)
	}
	if uint32(cb) != ib {
		t.Errorf("%d != %d", cb, ib)
	}
	if uint32(ca) != ia {
		t.Errorf("%d != %d", ca, ia)
	}
}
