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
	"embed"
	"fmt"
	"path"
	"strings"
)

type wordIndex uint16

const (
	languageFilesDirectory = "assets/bip39"
)

//go:embed assets/bip39
var bip39Dir embed.FS

func getSupportedLanguages() ([]string, error) {
	dirEntries, err := bip39Dir.ReadDir(languageFilesDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded directory: %w", err)
	}

	var supportedLanguages []string
	for _, dirEntry := range dirEntries {
		supportedLanguages = append(supportedLanguages, strings.TrimSuffix(dirEntry.Name(), ".txt"))
	}

	return supportedLanguages, nil
}

func readLanguageFile(language string) ([]string, error) {
	languageFilename := fmt.Sprintf("%s.txt", language)

	dirEntries, err := bip39Dir.ReadDir(languageFilesDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded directory: %w", err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.Name() != languageFilename {
			continue
		}

		data, err := bip39Dir.ReadFile(path.Join(languageFilesDirectory, languageFilename))
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded file: %w", err)
		}

		return strings.Split(string(data), "\n"), nil
	}

	return nil, fmt.Errorf("unsupported language \"%s\"", language)
}
