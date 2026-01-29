# gpl2asm

A command-line utility that converts GIMP palette files (`.gpl`) to assembly include files (`.inc`) for retro programming and palette-based systems.

## Overview

`gpl2asm` reads color palettes in GIMP's palette format and converts them to assembly language `DB` (Define Byte) directives. The tool automatically converts 8-bit RGB values (0-255) to 6-bit RGB values (0-63), which is the standard for VGA mode 13h.

## Features

- Converts GIMP palette files to assembly include files
- Automatic 8-bit to 6-bit RGB color conversion
- Cross-platform support (Linux, Windows, macOS)
- Available as standalone binary or Docker container

## Installation

### Download Pre-built Binary

Download the latest release for your platform from the [releases page](https://github.com/yourusername/gpl2asm/releases).

### Using Docker

```bash
docker pull ghcr.io/laghoule/gpl2asm:latest
```

Or from Docker Hub:

```bash
docker pull laghoule/gpl2asm:latest
```

### Build from Source

Requires Go 1.25.6 or later:

```bash
git clone https://github.com/laghoule/gpl2asm.git
cd gpl2asm
go build
```

## Usage

### Basic Usage

```bash
gpl2asm --src input.gpl --dst output.inc
```

### Using Default Files

If no arguments are provided, it uses `pal.gpl` as input and `pal.inc` as output:

```bash
gpl2asm
```

### Command-line Options

- `--src <file>`: Path to the GIMP palette file to convert (default: `pal.gpl`)
- `--dst <file>`: Path to the output assembly include file (default: `pal.inc`)

### Docker Usage

```bash
docker run -v $(pwd):/data ghcr.io/yourusername/gpl2asm:latest \
  --src /data/input.gpl --dst /data/output.inc
```

## Input Format

The tool expects GIMP palette format (`.gpl` files):

```
GIMP Palette
Name: My Palette
Columns: 16
#
  0   0   0 Black
255 255 255 White
255   0   0 Red
  0 255   0 Green
  0   0 255 Blue
```

Each color line consists of:
- Three integers (0-255) representing RGB values
- A name or comment for the color (can include spaces)

## Output Format

The tool generates assembly include files with `DB` directives:

```asm
palette LABEL BYTE
  DB 00,00,00 ; Black
  DB 63,63,63 ; White
  DB 63,00,00 ; Red
  DB 00,63,00 ; Green
  DB 00,00,63 ; Blue
```

Note that RGB values are automatically converted from 8-bit (0-255) to 6-bit (0-63) by dividing by 4.
