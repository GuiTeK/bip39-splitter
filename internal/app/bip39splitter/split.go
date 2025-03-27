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
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/GuiTeK/bip39-splitter/pkg/shamir"
)

func mnemonicToWordIndices(language string, mnemonic string) ([]wordIndex, error) {
	languageWords, err := readLanguageFile(language)
	if err != nil {
		return nil, fmt.Errorf("failed to read language file: %w", err)
	}

	words := strings.Split(mnemonic, " ")
	var wordIndices []wordIndex
	for _, word := range words {
		word := strings.TrimSpace(word)
		if len(word) == 0 {
			continue
		}

		i := slices.Index(languageWords, word)
		if i == -1 {
			return nil, fmt.Errorf("word \"%s\" does not exist in BIP-39 language \"%s\"", word, language)
		}

		if i >= len(languageWords) {
			return nil, fmt.Errorf(
				"invalid word index %d for word \"%s\" in BIP-39 language \"%s\"",
				i,
				word,
				language,
			)
		}

		wordIndices = append(wordIndices, wordIndex(i))
	}

	return wordIndices, nil
}

func wordIndicesToBytes(wordIndices []wordIndex) []byte {
	var buffer []byte
	for _, value := range wordIndices {
		buffer = binary.BigEndian.AppendUint16(buffer, uint16(value))
	}
	return buffer
}

func split(language string, parts int, threshold int) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your BIP-39 mnemonic (seed phrase): ")
	mnemonic, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	mnemonic = strings.TrimSpace(mnemonic)

	wordIndices, err := mnemonicToWordIndices(language, mnemonic)
	if err != nil {
		return fmt.Errorf("failed to convert mnemonic to word indices: %w", err)
	}

	wordBytes := wordIndicesToBytes(wordIndices)

	partsBytes, err := shamir.Split(wordBytes, parts, threshold)
	if err != nil {
		return fmt.Errorf("failed to Shamir-split mnemonic: %w", err)
	}

	separatorLine := strings.Repeat("=", 98)
	fmt.Printf("\nParts:\n%s\n", separatorLine)
	for _, partBytes := range partsBytes {
		fmt.Println(hex.EncodeToString(partBytes))
	}
	fmt.Println(separatorLine)

	return nil
}
