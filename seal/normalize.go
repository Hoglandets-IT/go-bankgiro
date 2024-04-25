package seal

import (
	"bytes"

	"github.com/hoglandets-it/go-bankgiro/tools"
)

const NormLowerLimit = 32
const NormUpperLimit = 126

const NormLfChar = 10
const NormCrChar = 13

const NormOutOfRangeReplacement = 195

var NormSpecialReplacement map[int]int = map[int]int{
	201: 64,
	196: 91,
	214: 92,
	197: 93,
	220: 94,
	233: 96,
	228: 123,
	246: 124,
	229: 125,
	252: 126,
}

// Normalize a single character
func NormalizeByte(b byte) byte {
	intChar := int(b)

	// Carriage returns and line feeds should be entirely ignored
	if NormLfChar == intChar || NormCrChar == intChar {
		return byte(0)
	}

	// Values outside of the printable ASCII range should be replaced
	if intChar < NormLowerLimit || intChar > NormUpperLimit {

		// Some characters want special replacements, others are replaced with a generic replacement char (195)
		if NormSpecialReplacement[intChar] != 0 {
			return byte(NormSpecialReplacement[intChar])
		}

		return byte(195)
	}

	return b
}

// Normalize a range of bytes
func NormalizeBytes(b []byte, buf *bytes.Buffer) {
	for _, ub := range b {
		nb := NormalizeByte(ub)
		if nb != 0 {
			buf.WriteByte(nb)
		}
	}
}

// Ensure that the first line of the content is a maximum of 80 characters long
func NormalizeFirstLine(b []byte) []byte {
	lines := bytes.Split(b, []byte{NormCrChar, NormLfChar})
	if len(lines) == 0 || len(lines[0]) == 80 {
		return b
	}

	lines[0] = lines[0][:80]
	return bytes.Join(lines, []byte{NormCrChar, NormLfChar})
}

func NormalizeContent(content []byte) []byte {
	var buf bytes.Buffer
	leReplaced := tools.EnsureCrlfBytes(content)

	NormalizeBytes(leReplaced, &buf)

	return buf.Bytes()
}

func NormalizeContentString(content string) []byte {
	return NormalizeContent([]byte(content))
}
