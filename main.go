package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

const (
	// colorBitShift converts 8-bit RGB (0-255) to 6-bit RGB (0-63)
	colorBitShift = 4
	minRGB        = 0
	maxRGB        = 255
)

var (
	version   = "unknown"
	gitCommit = "unknown"
)

func main() {
	srcName := flag.String("src", "pal.gpl", "Path to the GIMP palette to convert")
	dstName := flag.String("dst", "pal.inc", "Path to the output .inc file")
	flag.Parse()

	fmt.Printf("gpl2asm version: %s, git commit: %s\n", version, gitCommit)

	srcPal, err := os.Open(*srcName)
	if err != nil {
		msg := fmt.Errorf("Error opening source file: %v", err)
		exitWithError(msg)
	}
	defer srcPal.Close()

	dstPal, err := os.Create(*dstName)
	if err != nil {
		msg := fmt.Errorf("Error creating destination file: %v", err)
		exitWithError(msg)
	}
	defer dstPal.Close()

	// GIMP palette format (R G B Name)
	// EX: 0 0 0 Index0
	// Pattern allows spaces and special characters in color names
	re := regexp.MustCompile(`^\s*(\d+)\s+(\d+)\s+(\d+)\s+(.+?)\s*$`)

	palReader := bufio.NewReader(srcPal)

	var i int
	for {
		line, err := palReader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			fmt.Printf("Palette conversion complete, saved in file %s.\n", *dstName)
			break
		} else if err != nil {
			msg := fmt.Errorf("Error reading line: %v", err)
			exitWithError(msg)
		}

		match := re.FindStringSubmatch(line)
		if match != nil {
			red, green, blue := match[1], match[2], match[3]
			comment := match[4]
			r, g, b, err := convertTo6Bits(red, green, blue)
			if err != nil {
				msg := fmt.Errorf("Error converting color on line %d: %v", i+1, err)
				exitWithError(msg)
			}
			if i == 0 {
				fmt.Fprintf(dstPal, "palette LABEL BYTE\n  DB %02d,%02d,%02d ; %s\n", r, g, b, comment)
			} else {
				fmt.Fprintf(dstPal, "  DB %02d,%02d,%02d ; %s\n", r, g, b, comment)
			}
			i++
		}
	}
}

func convertTo6Bits(r, g, b string) (int, int, int, error) {
	red, err := strconv.Atoi(r)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid red value '%s': %w", r, err)
	}
	green, err := strconv.Atoi(g)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid green value '%s': %w", g, err)
	}
	blue, err := strconv.Atoi(b)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid blue value '%s': %w", b, err)
	}

	// Validate RGB ranges
	if red < minRGB || red > maxRGB {
		return 0, 0, 0, fmt.Errorf("red value %d out of range [%d-%d]", red, minRGB, maxRGB)
	}
	if green < minRGB || green > maxRGB {
		return 0, 0, 0, fmt.Errorf("green value %d out of range [%d-%d]", green, minRGB, maxRGB)
	}
	if blue < minRGB || blue > maxRGB {
		return 0, 0, 0, fmt.Errorf("blue value %d out of range [%d-%d]", blue, minRGB, maxRGB)
	}

	// Convert 8-bit (0-255) to 6-bit (0-63) by dividing by 4
	return red / colorBitShift, green / colorBitShift, blue / colorBitShift, nil
}

func exitWithError(err error) {
	fmt.Printf("%v\n", err)
	os.Exit(1)
}
