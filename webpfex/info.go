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
	var lines []string
	for _, l := range strings.Split(info, "\n") {
		if l != "" {
			lines = append(lines, l)
		}
	}

	canvasSizePattern := regexp.MustCompile(`^Canvas size\s*:\s*(\d+) x (\d+)$`)
	canvasSizeMatches := canvasSizePattern.FindStringSubmatch(lines[0])
	if len(canvasSizeMatches) != 3 {
		return AWebpInfo{}, makeParsingError("Failed parsing width and height", lines[0])
	}
	width, err := strconv.ParseUint(canvasSizeMatches[1], 10, 32)
	if err != nil {
		return AWebpInfo{}, makeParsingError("Failed parsing width", lines[0])
	}
	height, err := strconv.ParseUint(canvasSizeMatches[2], 10, 32)
	if err != nil {
		return AWebpInfo{}, makeParsingError("Failed parsing height", lines[0])
	}

	backgroundColorPattern := regexp.MustCompile(`^Background color\s*:\s*(0x[\dA-F]{8})`)
	backgroundColorMatches := backgroundColorPattern.FindStringSubmatch(lines[2])
	if len(backgroundColorMatches) != 2 {
		return AWebpInfo{}, makeParsingError("Failed parsing background color", lines[2])
	}
	backgroundColor, ok := parseHexColor(backgroundColorMatches[1])
	if !ok {
		return AWebpInfo{}, makeParsingError("Failed parsing background color", lines[2])
	}

	frameCountPattern := regexp.MustCompile(`Number of frames\s*:\s*(\d+)`)
	frameCountMatches := frameCountPattern.FindStringSubmatch(lines[3])
	if len(frameCountMatches) != 2 {
		return AWebpInfo{}, makeParsingError("Failed parsing frame count", lines[3])
	}
	frameCount, err := strconv.ParseUint(frameCountMatches[1], 10, 32)
	if err != nil {
		return AWebpInfo{}, makeParsingError("Failed parsing frame count", lines[3])
	}

	frameInfoLines := lines[5:]
	frameInfos := make([]AWebpFrameInfo, 0, len(frameInfoLines))
	for _, fil := range frameInfoLines {
		fi, err := ParseAWebpFrameInfo(fil)
		if err != nil {
			return AWebpInfo{}, makeErrorSub("Failed parsing frame info line", fil, err)
		}

		frameInfos = append(frameInfos, fi)
	}

	return MakeAWebpInfo(
		uint32(width),
		uint32(height),
		backgroundColor,
		uint32(frameCount),
		frameInfos,
	), nil
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

func ParseAWebpFrameInfo(line string) (AWebpFrameInfo, error) {
	fields := strings.Fields(line)

	number, err := strconv.ParseUint(regexp.MustCompile(`\d+`).FindString(fields[0]), 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing frame number", fields[0])
	}

	width, err := strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing width", fields[1])
	}
	height, err := strconv.ParseUint(fields[2], 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing height", fields[2])
	}

	var alpha bool
	switch fields[3] {
	case "yes":
		alpha = true
	case "no":
		alpha = false
	default:
		return AWebpFrameInfo{}, makeParsingError("Failed parsing alpha", fields[3])
	}

	xOffset, err := strconv.ParseUint(fields[4], 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing x offset", fields[4])
	}
	yOffset, err := strconv.ParseUint(fields[5], 10, 32)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing y offset", fields[5])
	}

	durationMs, err := strconv.ParseInt(fields[6], 10, 64)
	if err != nil {
		return AWebpFrameInfo{}, makeParsingError("Failed parsing duration", fields[6])
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

func makeErrorSub(reason string, input string, subError error) *ParsingError {
	e := ParsingError{reason, input, subError}
	return &e
}

func (e ParsingError) Error() string {
	if e.subError != nil {
		return fmt.Sprintf("ParsingError: %s: %q\n%s", e.reason, e.input, e.subError.Error())
	} else {
		return fmt.Sprintf("ParsingError: %s: %q", e.reason, e.input)
	}
}
