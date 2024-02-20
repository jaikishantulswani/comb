package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var flip bool
	flag.BoolVar(&flip, "f", false, "")
	flag.BoolVar(&flip, "flip", false, "")

	var separator string
	flag.StringVar(&separator, "s", "", "")
	flag.StringVar(&separator, "separator", "", "")

	flag.Parse()

	if flag.NArg() < 2 || flag.NArg() > 4 {
		flag.Usage()
		os.Exit(1)
	}

	files := make([]io.Reader, flag.NArg())
	var err error
	for i := 0; i < flag.NArg(); i++ {
		if flag.Arg(i) == "-" {
			files[i] = os.Stdin
		} else {
			files[i], err = os.Open(flag.Arg(i))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}
	}

	// Create a slice of slices to store lines from each file
	lines := make([][]string, len(files))
	for i, file := range files {
		sc := bufio.NewScanner(file)
		for sc.Scan() {
			lines[i] = append(lines[i], sc.Text())
		}
	}

	combineLines(lines, flip, separator)
}

func combineLines(lines [][]string, flip bool, separator string) {
	// Generate combinations and print
	for _, combo := range generateCombinations(lines) {
		for i, line := range combo {
			if flip {
				fmt.Print(line, separator)
			} else {
				if i > 0 {
					fmt.Print(separator)
				}
				fmt.Print(line)
			}
		}
		fmt.Println()
	}
}

func generateCombinations(lines [][]string) [][]string {
	if len(lines) == 0 {
		return [][]string{{}}
	}

	combinations := generateCombinations(lines[1:])
	result := make([][]string, 0)

	for _, prefix := range lines[0] {
		for _, combo := range combinations {
			result = append(result, append([]string{prefix}, combo...))
		}
	}

	return result
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Combine the lines from two to four files in every combination. Use '-' to read from stdin.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  comb [OPTIONS] [FILE1|-] [FILE2|-] [FILE3|-] [FILE4|-]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -f, --flip             Flip mode (order by suffix)\n")
		fmt.Fprintf(os.Stderr, "  -s, --separator <str>  String to place between prefix and suffix\n")
	}
}
