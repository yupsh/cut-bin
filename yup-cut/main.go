package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"

	gloo "github.com/gloo-foo/framework"
	. "github.com/yupsh/cut"
)

const (
	flagDelimiter     = "delimiter"
	flagFields        = "fields"
	flagChars         = "characters"
	flagBytes         = "bytes"
	flagOnlyDelimited = "only-delimited"
)

func main() {
	app := &cli.App{
		Name:  "cut",
		Usage: "remove sections from each line of files",
		UsageText: `cut OPTION... [FILE...]

   Print selected parts of lines from each FILE to standard output.
   With no FILE, or when FILE is -, read standard input.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagDelimiter,
				Aliases: []string{"d"},
				Usage:   "use DELIM instead of TAB for field delimiter",
			},
			&cli.StringFlag{
				Name:    flagFields,
				Aliases: []string{"f"},
				Usage:   "select only these fields (comma-separated list)",
			},
			&cli.StringFlag{
				Name:    flagChars,
				Aliases: []string{"c"},
				Usage:   "select only these characters (comma-separated list)",
			},
			&cli.StringFlag{
				Name:    flagBytes,
				Aliases: []string{"b"},
				Usage:   "select only these bytes (comma-separated list)",
			},
			&cli.BoolFlag{
				Name:    flagOnlyDelimited,
				Aliases: []string{"s"},
				Usage:   "do not print lines not containing delimiters",
			},
		},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "cut: %v\n", err)
		os.Exit(1)
	}
}

func action(c *cli.Context) error {
	var params []any

	// Add file arguments (or none for stdin)
	for i := 0; i < c.NArg(); i++ {
		params = append(params, gloo.File(c.Args().Get(i)))
	}

	// Add flags based on CLI options
	if c.IsSet(flagDelimiter) {
		params = append(params, Delimiter(c.String(flagDelimiter)))
	}
	if c.IsSet(flagFields) {
		fields, err := parseList(c.String(flagFields))
		if err != nil {
			return fmt.Errorf("invalid field list: %w", err)
		}
		params = append(params, Fields(fields))
	}
	if c.IsSet(flagChars) {
		chars, err := parseList(c.String(flagChars))
		if err != nil {
			return fmt.Errorf("invalid character list: %w", err)
		}
		params = append(params, Chars(chars))
	}
	if c.IsSet(flagBytes) {
		bytes, err := parseList(c.String(flagBytes))
		if err != nil {
			return fmt.Errorf("invalid byte list: %w", err)
		}
		params = append(params, Bytes(bytes))
	}
	if c.Bool(flagOnlyDelimited) {
		params = append(params, OnlyDelimited)
	}

	// Create and execute the cut command
	cmd := Cut(params...)
	return gloo.Run(cmd)
}

// parseList parses a comma-separated list of integers and ranges (e.g., "1,3,5-7")
func parseList(s string) ([]int, error) {
	var result []int
	parts := strings.Split(s, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check if it's a range (e.g., "5-7")
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", rangeParts[1])
			}

			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			// Single number
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", part)
			}
			result = append(result, num)
		}
	}

	return result, nil
}
