// txt2bdf converts an ASCII matrix font file to a BDF font file.
//
// Input format:
//
//	FOUNDRY mtraver
//	FAMILY ledmatrix
//	WEIGHT light
//	SLANT r
//	WIDTH normal
//	STYLE sans serif
//	DPI 6  // The DPI of a 4mm pitch LED matrix is 6.4
//	SPACING c
//	CHARSET_REGISTRY ISO8859
//	CHARSET_ENCODING 1
//
//	CHAR 0x41  // "A"
//	.###..
//	#...#.
//	#...#.
//	#####.
//	#...#.
//	#...#.
//	#...#.
//
//	CHAR 0x42  // "B"
//	####..
//	#...#.
//	#...#.
//	####..
//	#...#.
//	#...#.
//	####..
//
// Usage:
//
//	go run main.go input.txt output.bdf
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mtraver/matrixfont/parse"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: txt2bdf input.txt output.bdf\n")
		os.Exit(2)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	font, err := parse.Parse(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed %d glyphs\n", len(font.Glyphs))

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)
	if err := font.WriteBDF(w); err != nil {
		fmt.Fprintf(os.Stderr, "Error converting to BDF: %v\n", err)
		os.Exit(1)
	}

	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Error flushing output file: %v\n", err)
		os.Exit(1)
	}
}
