package webpfex

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
	"webpfex/canvas"
)

type AWebpInfo struct {
	Width           uint32
	Height          uint32
	BackgroundColor canvas.Color
	FrameCount      uint32
	FrameInfos      []AWebpFrameInfo
}

func MakeAWebpInfo(
	width, height uint32,
	backgroundColor canvas.Color,
	frameCount uint32,
	frameInfos []AWebpFrameInfo,
) AWebpInfo {
	if frameCount != uint32(len(frameInfos)) {
		panic(fmt.Sprintf(
			"FrameCount %d doesn't match frame infos len %d.", frameCount, len(frameInfos)))
	}

	return AWebpInfo{
		Width:           width,
		Height:          height,
		BackgroundColor: backgroundColor,
		FrameCount:      frameCount,
		FrameInfos:      frameInfos,
	}
}

func ParseAWebpInfo(info string) (AWebpInfo, error) {
	if strings.Contains(info, "No features present.") {
		return AWebpInfo{}, makeParsingError("Not an animated WEBP", info)
	}

	width, height, err := parseAWebpInfoCanvasSize(info)
	if err != nil {
		return AWebpInfo{}, err
	}
	backgroundColor, err := parseAWebpInfoBackgroundColor(info)
	if err != nil {
		return AWebpInfo{}, err
	}
	frameCount, frameInfos, err := parseAWebpInfoFrames(info)
	if err != nil {
		return AWebpInfo{}, err
	}

	return MakeAWebpInfo(
		width,
		height,
		backgroundColor,
		frameCount,
		frameInfos,
	), nil
}

func parseAWebpInfoCanvasSize(info string) (uint32, uint32, error) {
	pattern := regexp.MustCompile(`(?m)^Canvas size: (\d+) x (\d+)$`)
	matches := pattern.FindStringSubmatch(info)

	width, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return 0, 0, makeParsingError("Failed parsing width", info)
	}
	height, err := strconv.ParseUint(matches[2], 10, 32)
	if err != nil {
		return 0, 0, makeParsingError("Failed parsing height", info)
	}

	return uint32(width), uint32(height), nil
}

func parseAWebpInfoBackgroundColor(info string) (canvas.Color, error) {
	pattern := regexp.MustCompile(`(?m)^Background color\s*:\s*(0x[\dA-F]{8})`)
	matches := pattern.FindStringSubmatch(info)
	hex, ok := parseHexColor(matches[1])
	if !ok {
		return canvas.Color{}, makeParsingError("Failed parsing background color", info)
	}

	return hex, nil
}

func parseAWebpInfoFrames(info string) (uint32, []AWebpFrameInfo, error) {
	pattern := regexp.MustCompile(`(?m)^Number of frames: (\d+)`)
	matches := pattern.FindStringSubmatch(info)
	frameCount, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return 0, []AWebpFrameInfo{}, makeParsingError("Failed parsing frame count", info)
	}

	lines := strings.Split(info, "\n")

	frameInfoLines := lines[5 : 5+frameCount]
	var frameInfos []AWebpFrameInfo
	for _, l := range frameInfoLines {
		fi, err := parseAWebpFrameInfo(l)
		if err != nil {
			return 0, []AWebpFrameInfo{}, makeParsingError("Failed parsing frame info", l)
		}

		frameInfos = append(frameInfos, fi)
	}

	return uint32(frameCount), frameInfos, nil
}

func parseHexColor(s string) (canvas.Color, bool) {
	hex, ok := new(big.Int).SetString(s, 0)
	if !ok {
		return canvas.Color{}, false
	}

	if !hex.IsUint64() {
		panic(fmt.Sprintf("Color %s is too big to convert to uint64.", s))
	}

	return canvas.MakeColor(hex.Uint64()), true
}

type AWebpFrameInfo struct {
	Number   uint32
	Width    uint32
	Height   uint32
	Alpha    bool
	XOffset  uint32
	YOffset  uint32
	Duration time.Duration
	Blend    bool
}

func MakeAWebpFrameInfo(
	number,
	width, height uint32,
	alpha bool,
	xOffset, yOffset uint32,
	duration time.Duration,
	blend bool,
) AWebpFrameInfo {
	if number == 0 {
		panic("number cannot be 0.")
	}

	return AWebpFrameInfo{
		Number:   number,
		Width:    width,
		Height:   height,
		Alpha:    alpha,
		XOffset:  xOffset,
		YOffset:  yOffset,
		Duration: duration,
		Blend:    blend,
	}
}

func parseAWebpFrameInfo(line string) (AWebpFrameInfo, error) {
	fields := strings.Fields(line)

	number, err := strconv.ParseUint(regexp.MustCompile(`\d+`).FindString(fields[0]), 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing frame number", line)
	}

	width, err := strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing width", line)
	}
	height, err := strconv.ParseUint(fields[2], 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing height", line)
	}

	var alpha bool
	switch fields[3] {
	case "yes":
		alpha = true
	case "no":
		alpha = false
	default:
		return AWebpFrameInfo{}, makeParsingError("Failed parsing alpha", line)
	}

	xOffset, err := strconv.ParseUint(fields[4], 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing x offset", line)
	}
	yOffset, err := strconv.ParseUint(fields[5], 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing y offset", line)
	}

	durationMs, err := strconv.ParseInt(fields[6], 10, 64)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing duration", line)
	}
	duration := time.Duration(durationMs * int64(time.Millisecond))

	var blend bool
	switch fields[8] {
	case "yes":
		blend = true
	case "no":
		blend = false
	default:
		return AWebpFrameInfo{}, makeParsingError("Failed parsing blend", fields[8])
	}

	return MakeAWebpFrameInfo(
		uint32(number),
		uint32(width), uint32(height),
		alpha,
		uint32(xOffset), uint32(yOffset),
		duration,
		blend,
	), nil
}

type ParsingError struct {
	reason   string
	input    string
	subError error
}

func makeParsingError(reason string, input string) *ParsingError {
	e := ParsingError{reason, input, nil}
	return &e
}

func (e ParsingError) Error() string {
	if e.subError != nil {
		return fmt.Sprintf("ParsingError: %s: %q\n%s", e.reason, e.input, e.subError.Error())
	} else {
		return fmt.Sprintf("ParsingError: %s: %q", e.reason, e.input)
	}
}
