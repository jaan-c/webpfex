package webpfex

import (
	"testing"
	"time"
	"webpfex/canvas"
)

const AWEBP_INFO_DUMMY = `Canvas size: 640 x 640
Features present: animation EXIF metadata transparency
Background color : 0xFFFFFFFF  Loop Count : 0
Number of frames: 8
No.: width height alpha x_offset y_offset duration   dispose blend image_size  compression
  1:   640   640    no        0        0       40       none    no      22710       lossy
  2:   640   577   yes        0       28       40       none   yes      25074       lossy
  3:   640   574   yes        0       28       80       none   yes      24066       lossy
  4:   640   570   yes        0       32       40       none   yes      25690       lossy
  5:   640   573   yes        0       32       40       none   yes      23262       lossy
  6:   640   640    no        0        0       40       none    no      23344       lossy
  7:   640   576   yes        0       26       80       none   yes      25570       lossy
  8:   640   579   yes        0       26       40       none   yes      23316       lossy
Size of the EXIF metadata: 34
`

func TestParseAWebpInfoCanvasSize(t *testing.T) {
	width, height, err := parseAWebpInfoCanvasSize(AWEBP_INFO_DUMMY)
	var expectedWidth uint32 = 640
	var expectedHeight uint32 = 640

	if width != expectedWidth {
		t.Errorf("Expecting %d got %d", width, expectedWidth)
	}
	if height != expectedHeight {
		t.Errorf("Expecting %d got %d", height, expectedHeight)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestParseAWebpInfoBackgroundColor(t *testing.T) {
	backgroundColor, err := parseAWebpInfoBackgroundColor(AWEBP_INFO_DUMMY)
	expectedBackgroundColor := canvas.MakeColor(0xFFFFFFFF)

	if backgroundColor != expectedBackgroundColor {
		t.Errorf("Expecting %X got %X",
			backgroundColor.Value(), expectedBackgroundColor.Value())
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestParseAWebpInfoFrames(t *testing.T) {
	frameCount, frameInfo, err := parseAWebpInfoFrames(AWEBP_INFO_DUMMY)
	var expectedFrameCount uint32 = 8
	expectedFrameInfo1 := MakeAWebpFrameInfo(
		1, 640, 640, false, 0, 0, 40*time.Millisecond, false)
	expectedFrameInfo2 := MakeAWebpFrameInfo(
		2, 640, 577, true, 0, 28, 40*time.Millisecond, true)
	expectedFrameInfo3 := MakeAWebpFrameInfo(
		3, 640, 574, true, 0, 28, 80*time.Millisecond, true)

	if frameCount != expectedFrameCount {
		t.Errorf("Expecting %d got %d", frameCount, expectedFrameCount)
	}
	if frameInfo[0] != expectedFrameInfo1 {
		t.Errorf("Expecting %v got %v", frameInfo[0], expectedFrameInfo1)
	}
	if frameInfo[1] != expectedFrameInfo2 {
		t.Errorf("Expecting %v got %v", frameInfo[0], expectedFrameInfo2)
	}
	if frameInfo[2] != expectedFrameInfo3 {
		t.Errorf("Expecting %v got %v", frameInfo[0], expectedFrameInfo3)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
