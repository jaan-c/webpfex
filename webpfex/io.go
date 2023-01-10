package webpfex

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path"
	"strconv"
	"webpfex/canvas"

	png "image/png"

	_ "golang.org/x/image/webp"
)

func ExtractWebpFramesAsPng(webp string, outdir string) error {
	err := os.Mkdir(outdir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	info, err := ExtractAWebpInfo(webp)
	if err != nil {
		return err
	}

	canvas := canvas.MakeCanvas(info.Width, info.Height)
	ClearCanvas(&canvas, info.BackgroundColor)

	for _, frameInfo := range info.FrameInfos {
		overlay, err := LoadAWebpFrame(webp, frameInfo.Number)
		if err != nil {
			panic(err)
		}

		if frameInfo.Blend {
			OverlayBlendCanvas(&canvas, &overlay, frameInfo.XOffset, frameInfo.YOffset)
		} else {
			OverlayCanvas(&canvas, &overlay, frameInfo.XOffset, frameInfo.YOffset)
		}

		outpath := path.Join(outdir, fmt.Sprintf("%09d.png", frameInfo.Number))
		err = SavePng(canvas, outpath)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func ConvertWebpToMp4(webp string, out string) error {
	frameDir, err := os.MkdirTemp("", "webpfex")
	if err != nil {
		return err
	}
	defer os.RemoveAll(frameDir)

	info, err := ExtractAWebpInfo(webp)
	if err != nil {
		return err
	}

	if err := ExtractWebpFramesAsPng(webp, frameDir); err != nil {
		return err
	}

	frameDuration := info.FrameInfos[0].Duration
	frameRate := int(1000 / frameDuration.Milliseconds())

	cmd := exec.Command("ffmpeg",
		"-framerate", strconv.Itoa(frameRate),
		"-pattern_type", "glob",
		"-i", path.Join(frameDir, "*.png"),
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		out,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg: %s", stderr.String())
	}

	return nil
}

// Extract nth frame from an animated WEBP image; indexing starts at 1. Relies
// on webpmux command.
func LoadAWebpFrame(path string, n uint32) (canvas.Canvas, error) {
	cmd := exec.Command(
		"webpmux",
		"-get", "frame", strconv.FormatUint(uint64(n), 10),
		path,
		"-o", "-")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return canvas.Canvas{}, fmt.Errorf("webpmux: %s", stderr.String())
	}

	img, _, err := image.Decode(bytes.NewReader(stdout.Bytes()))
	if err != nil {
		return canvas.Canvas{}, err
	}

	return ImageToCanvas(img), nil
}

// Extract metadata from an animated WEBP image. Relies on webpmux command.
func ExtractAWebpInfo(path string) (AWebpInfo, error) {
	cmd := exec.Command("webpmux", "-info", path)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return AWebpInfo{}, fmt.Errorf("webpmux: %s", stderr.String())
	}

	if info, err := ParseAWebpInfo(stdout.String()); err != nil {
		return AWebpInfo{}, err
	} else {
		return info, nil
	}
}

func LoadWebp(path string) (canvas.Canvas, error) {
	reader, err := os.Open(path)
	if err != nil {
		return canvas.Canvas{}, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return canvas.Canvas{}, err
	}

	return ImageToCanvas(img), nil
}

func SavePng(cv canvas.Canvas, path string) error {
	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer writer.Close()

	if err := png.Encode(writer, CanvasToImage(cv)); err != nil {
		return nil
	}

	return nil
}
