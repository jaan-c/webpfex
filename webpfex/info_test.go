package webpfex

import (
	"testing"
	"time"
	"webpfex/canvas"
)

func TestParseAWebpMetadata(t *testing.T) {
	stdout := `Canvas size: 400 x 400
Features present: animation transparency
Background color : 0xFF000000  Loop Count : 0
Number of frames: 12
No.: width height alpha x_offset y_offset duration   dispose blend image_size  compression
  1:   400   400    no        0        0       70       none    no       5178    lossless
  2:   400   400   yes        0        0       70       none   yes       1386    lossless
  3:   400   400   yes        0        0       70       none   yes       1472    lossless
  4:   400   394   yes        0        6       70       none   yes       3212    lossless
  5:   371   394   yes        0        6       70       none   yes       1888    lossless
  6:   394   382   yes        6        6       70       none   yes       3346    lossless
  7:   400   388   yes        0        0       70       none   yes       3786    lossless
  8:   394   383   yes        0        0       70       none   yes       1858    lossless
  9:   394   394   yes        0        6       70       none   yes       3794    lossless
 10:   372   394   yes       22        6       70       none   yes       3458    lossless
 11:   400   400    no        0        0       70       none    no       5270    lossless
 12:   320   382   yes        0        6       70       none   yes       2506    lossless
`
	var width uint32 = 400
	var height uint32 = 400
	backgroundColor := canvas.MakeColor(4278190080)
	var frameCount uint32 = 12

	info, err := ParseAWebpInfo(stdout)
	if err != nil {
		t.Error(err.Error())
	}
	if info.Width != width {
		t.Errorf("Expecting %d, got %d", width, info.Width)
	}
	if info.Height != height {
		t.Errorf("Expecting %d, got %d", height, info.Height)
	}
	if info.BackgroundColor != backgroundColor {
		t.Errorf("Expecting %q, got %q", backgroundColor, info.BackgroundColor)
	}
	if info.FrameCount != frameCount {
		t.Errorf("Expecting %d, got %d", frameCount, info.FrameCount)
	}
}

func TestParseAWebpFrameMetadata(t *testing.T) {
	//       No.: width height alpha x_offset y_offset duration   dispose blend   image_size  compression
	line := "2:   450   400    no        12        14       70       none    yes       5178    lossless"
	var number uint32 = 2
	var width uint32 = 450
	var height uint32 = 400
	alpha := false
	var xOffset uint32 = 12
	var yOffset uint32 = 14
	duration := 70 * time.Millisecond
	var blend bool = true

	meta, err := ParseAWebpFrameInfo(line)
	if err != nil {
		t.Error(err.Error())
	}
	if meta.Number != number {
		t.Errorf("Expecting %d, got %d", number, meta.Number)
	}
	if meta.Width != width {
		t.Errorf("Expecting %d, got %d", width, meta.Width)
	}
	if meta.Height != height {
		t.Errorf("Expecting %d, got %d", height, meta.Height)
	}
	if meta.Alpha != alpha {
		t.Errorf("Expecting %t, got %t", alpha, meta.Alpha)
	}
	if meta.XOffset != xOffset {
		t.Errorf("Expecting %d, got %d", xOffset, meta.XOffset)
	}
	if meta.YOffset != yOffset {
		t.Errorf("Expecting %d, got %d", yOffset, meta.YOffset)
	}
	if meta.Duration != duration {
		t.Errorf("Expecting %d, got %d", duration, meta.Duration)
	}
	if meta.Blend != blend {
		t.Errorf("Expecting %t, got %t", blend, meta.Blend)
	}
}
