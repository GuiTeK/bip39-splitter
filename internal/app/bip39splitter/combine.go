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
	"strings"

	"github.com/GuiTeK/bip39-splitter/pkg/shamir"
)

type invalidWordIndexError struct {
	wordIdx  wordIndex
	position int
	language string
}

func (e invalidWordIndexError) Error() string {
	return fmt.Sprintf(
		"invalid word index %d at word position %d in BIP-39 language \"%s\"",
		e.wordIdx,
		e.position,
		e.language,
	)
}

func wordIndicesFromBytes(buffer []byte) []wordIndex {
	var wordIndices []wordIndex
	for i := 0; i < len(buffer); i += 2 {
		wordIndices = append(wordIndices, wordIndex(binary.BigEndian.Uint16(buffer[i:i+2])))
	}
	return wordIndices
}

func mnemonicFromWordIndices(language string, wordIndices []wordIndex) (string, error) {
	languageWords, err := readLanguageFile(language)
	if err != nil {
		return "", fmt.Errorf("failed to read language file: %w", err)
	}

	mnemonic := ""
	for i, wordIdx := range wordIndices {
		if int(wordIdx) >= len(languageWords) {
			return "", invalidWordIndexError{
				wordIdx:  wordIdx,
				position: i,
				language: language,
			}
		}

		if mnemonic == "" {
			mnemonic = languageWords[wordIdx]
		} else {
			mnemonic += fmt.Sprintf(" %s", languageWords[wordIdx])
		}
	}

	return mnemonic, nil
}

func combine(language string) error {
	reader := bufio.NewReader(os.Stdin)

	var parts []string
	i := 1
	for {
		fmt.Printf("Enter part %d and press ENTER (leave empty to finish): ", i)
		share, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		share = strings.TrimSpace(share)
		if len(share) == 0 {
			break
		}
		parts = append(parts, share)
		i += 1
	}

	var partsBytes [][]byte
	for i, part := range parts {
		partBytes, err := hex.DecodeString(part)
		if err != nil {
			return fmt.Errorf("failed to decode part %d: %w", i+1, err)
		}
		partsBytes = append(partsBytes, partBytes)
	}

	wordsBytes, err := shamir.Combine(partsBytes)
	if err != nil {
		return fmt.Errorf("failed to Shamir-combine parts of mnemonic: %w", err)
	}

	wordIndices := wordIndicesFromBytes(wordsBytes)

	mnemonic, err := mnemonicFromWordIndices(language, wordIndices)
	if err != nil {
		return fmt.Errorf("failed to build mnemonic from word indices: %w", err)
	}

	fmt.Printf("\nReconstructed mnemonic (seed phrase):\n%s\n", mnemonic)

	return nil
}
