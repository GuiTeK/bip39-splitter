// Copyright (c) 2025 Guillaume Truchot (GuiTeK)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE
// OR OTHER DEALINGS IN THE SOFTWARE.
package bip39splitter

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	cmdSplit   = "split"
	cmdCombine = "combine"
)

//go:embed assets/cli/usage.txt
var cliUsage string

//go:embed assets/cli/version.txt
var version string

func printUsage(supportedLanguages []string) func() {
	return func() {
		executableFilename := filepath.Base(os.Args[0])
		_, _ = fmt.Fprintf(
			os.Stderr,
			fmt.Sprintf(cliUsage, executableFilename, version, strings.Join(supportedLanguages, ", ")),
		)
	}
}

func Run() {
	supportedLanguages, err := getSupportedLanguages()
	if err != nil {
		fmt.Printf("Error: failed to get BIP-39 supported languages: %s\n", err)
		os.Exit(1)
	}

	flag.Usage = printUsage(supportedLanguages)

	splitCmd := flag.NewFlagSet(cmdSplit, flag.ExitOnError)
	languageSplit := splitCmd.String(
		"l",
		"",
		fmt.Sprintf(
			"language of the BIP-39 mnemonic (seed phrase). One of:\n    %s",
			strings.Join(supportedLanguages, ", "),
		),
	)
	parts := splitCmd.Int("p", 0, "number of parts to split the mnemonic (seed phrase) into")
	threshold := splitCmd.Int(
		"t",
		0,
		"number of parts required to reconstruct the mnemonic (seed phrase) with command 'combine'",
	)

	combineCmd := flag.NewFlagSet(cmdCombine, flag.ExitOnError)
	languageCombine := combineCmd.String(
		"l",
		"",
		fmt.Sprintf(
			"language of the BIP-39 mnemonic (seed phrase). One of:\n    %s",
			strings.Join(supportedLanguages, ", "),
		),
	)

	help := flag.Bool("help", false, "Show help")
	helpShorthand := flag.Bool("h", false, "Show help (alias of -help)")

	if *help || *helpShorthand {
		flag.Usage()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case cmdSplit:
		if err := splitCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error: failed to parse arguments: %s\n", err)
			os.Exit(1)
		}
		if err := split(*languageSplit, *parts, *threshold); err != nil {
			fmt.Printf("\nError: %s\n", err)
			os.Exit(1)
		}
	case cmdCombine:
		if err := combineCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error: failed to parse arguments: %s\n", err)
			os.Exit(1)
		}
		if err := combine(*languageCombine); err != nil {
			fmt.Printf("\nError: %s\n", err)
			var invalidWordIndexErr invalidWordIndexError
			if errors.As(err, &invalidWordIndexErr) {
				fmt.Println(
					"This error usually means one or more of the parts are incorrect, or you don't have enough parts.",
				)
			}
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(2)
	}
}
