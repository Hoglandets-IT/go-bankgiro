package tools

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

const NormLfChar = 10
const NormCrChar = 13

func AutomaticDecoder(byteString []byte, tryEncodings ...*charmap.Charmap) (string, error) {
	var bodyString string

	for _, chm := range tryEncodings {
		encoder := chm.NewDecoder()
		bodyString, err := encoder.String(string(byteString))
		if err != nil || strings.Contains(bodyString, "ï¿½") || strings.Contains(bodyString, "�") {
			continue
		}
		return bodyString, nil
	}

	bodyString = string(byteString)

	if strings.Contains(bodyString, "�") {
		return "", fmt.Errorf("invalid encoding detected")
	}

	return bodyString, nil
}

func StringIsoEncoder(input string) (string, error) {
	encoder := charmap.ISO8859_1.NewEncoder()
	encodedString, err := encoder.String(input)
	if err != nil {
		return "", err
	}

	return encodedString, nil
}

func BytesIsoEncoder(input []byte) ([]byte, error) {
	encoder := charmap.ISO8859_1.NewEncoder()
	encodedBytes, err := encoder.Bytes(input)
	if err != nil {
		return []byte{}, err
	}

	return encodedBytes, nil
}

func IdentifyInvalidUtf8(text string) error {
	for _, r := range text {
		if r == utf8.RuneError {
			return fmt.Errorf("invalid rune detected")
		}

		if r == '�' {
			return fmt.Errorf("invalid rune detected")
		}
	}

	return nil
}

func IdentifyInvalidUtf8Bytes(byteText []byte) error {
	text := string(byteText)
	for _, r := range text {
		if r == utf8.RuneError {
			return fmt.Errorf("invalid rune detected")
		}

		if r == '�' {
			return fmt.Errorf("invalid rune detected")
		}
	}

	return nil
}

func StringEnsureIso(input string) string {
	encoded, err := StringIsoEncoder(input)
	if err != nil {
		if err := IdentifyInvalidUtf8(input); err != nil {
			return input
		}
	}

	return encoded
}

func BytesEnsureIso(input []byte) []byte {
	encoded, err := BytesIsoEncoder([]byte(input))
	if err != nil {
		if err := IdentifyInvalidUtf8Bytes(input); err != nil {
			return input
		}
	}

	return encoded
}

// Ensure the content is using CRLF line endings
func EnsureCrlfBytes(b []byte) []byte {
	onlyN := bytes.ReplaceAll(b, []byte{NormCrChar, NormLfChar}, []byte{NormLfChar})
	onlyN = bytes.ReplaceAll(onlyN, []byte{NormCrChar}, []byte{NormLfChar})

	return bytes.ReplaceAll(onlyN, []byte{NormLfChar}, []byte{NormCrChar, NormLfChar})
}

// Ensure the content is using CRLF line endings
func EnsureCrlfString(s string) string {
	onlyN := strings.ReplaceAll(s, "\r\n", "\n")
	onlyN = strings.ReplaceAll(onlyN, "\r", "\n")

	return strings.ReplaceAll(onlyN, "\n", "\r\n")
}

func BytesToIsoString(input []byte) (string, error) {
	if err := IdentifyInvalidUtf8(string(input)); err != nil {
		decoded, err := AutomaticDecoder(input, charmap.ISO8859_1, charmap.Windows1252)
		if err != nil {
			return string(input), nil
		}

		return decoded, nil
	}

	bodyString, err := charmap.ISO8859_1.NewEncoder().String(string(input))
	if err != nil {
		return "", fmt.Errorf("error transforming bytes: %w", err)
	}

	return EnsureCrlfString(bodyString), nil
}
