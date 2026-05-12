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
	"flag"
	"fmt"
	"os"

	"github.com/mtraver/matrixfont/log"
	"github.com/mtraver/matrixfont/parse"
)

// Flags.
var (
	flagVerbosity int
)

func init() {
	flag.IntVar(
		&flagVerbosity, "v", log.LevelWarn,
		fmt.Sprintf("log verbosity from %d (debug) to %d (error)", log.LevelDebug, log.LevelError),
	)
}

func parseFlags() error {
	flag.Parse()
	return nil
}

func main() {
	if err := parseFlags(); err != nil {
		fmt.Printf("argument error: %v\n", err)
		os.Exit(2)
	}

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}
	inputPath := args[0]
	outputPath := args[1]

	logger := log.New(os.Stdout, flagVerbosity)

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	font, err := parse.Parse(inputFile, parse.WithLogVerbosity(flagVerbosity))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input: %v\n", err)
		os.Exit(1)
	}

	logger.Infof("Parsed %d glyphs", len(font.Glyphs))

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
