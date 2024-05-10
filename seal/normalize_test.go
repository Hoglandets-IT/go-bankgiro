package seal_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/hoglandets-it/go-bankgiro/seal"
	"github.com/hoglandets-it/go-bankgiro/tools"
)

const (
	NormalizationTestCharset = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmno\r\npqrstuvwxyz{|}~"
	OutOfBoundsReplacement   = "Ã"
	OutOfBoundsLowerHex      = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	OutOfBoundsUpperHex      = "7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9f"
)

var BasicNormalization map[string]string = map[string]string{
	"test":                       "test",
	"1234567890 0":               "1234567890 0",
	"abcdefghijklmnopqrstuvwxyz": "abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"åäöÅÄÖ":                     "}{|][\\",
	"åäö\r\nÅÄÖ\r\n\r\nmultiline-string__": "}{|][\\multiline-string__",
}

// Ignore Chars
// \r\n

// Default Out of Bounds replacement:
// Ã

// Special Replacements:
// É Ä Ö Å Ü é ä ö å ü
// @ [ \ ] ^ ` { | } ~

func ErrorPrintComparison(t *testing.T, got interface{}, expected interface{}, more ...interface{}) {
	errStr := ""
	for _, a := range more {
		errStr += fmt.Sprintf("\r\n%v", a)
	}
	t.Errorf(
		"Normalization failed: \r\nGOT: %s \r\nEXP: %s %s",
		got,
		expected,
		errStr,
	)
}

func TestBasicNormalization(t *testing.T) {
	for input, expected := range BasicNormalization {
		isoString := tools.StringEnsureIso(input)
		actual := seal.NormalizeContentString(isoString)
		if string(actual) != expected {
			ErrorPrintComparison(t, string(actual), expected)
		}
	}
}

func TestUnescapedCharacters(t *testing.T) {
	charset := tools.StringEnsureIso(NormalizationTestCharset)
	actual := seal.NormalizeContentString(charset)
	expected := charset[0:80] + charset[82:]

	if string(actual) != expected {
		ErrorPrintComparison(t, actual, expected)
	}
}

func TestOutOfBounds(t *testing.T) {
	oobLowerString, err := hex.DecodeString(OutOfBoundsLowerHex)
	if err != nil {
		t.Errorf("Failed to encode out of bounds replacement: %s", err)
		return
	}

	oobUpperString, err := hex.DecodeString(OutOfBoundsUpperHex)
	if err != nil {
		t.Errorf("Failed to encode out of bounds replacement: %s", err)
		return
	}

	oobLowerHex := tools.BytesEnsureIso(oobLowerString)
	oobUpperHex := tools.BytesEnsureIso(oobUpperString)

	actualLower := seal.NormalizeContent(oobLowerHex)
	actualUpper := seal.NormalizeContent(oobUpperHex)

	oobLowerLen := len(actualLower)
	oobUpperLen := len(actualUpper)

	expectedLower := tools.BytesEnsureIso(bytes.Repeat([]byte(OutOfBoundsReplacement), oobLowerLen))
	expectedUpper := tools.BytesEnsureIso(bytes.Repeat([]byte(OutOfBoundsReplacement), oobUpperLen))

	if !tools.BytesliceEqual(actualLower, expectedLower) {
		ErrorPrintComparison(t, actualLower, expectedLower)
	}

	if !tools.BytesliceEqual(actualUpper, expectedUpper) {
		ErrorPrintComparison(t, actualUpper, expectedUpper)
	}
}

var FileNormalization []string = []string{
	"andringslista-new",
	"andringslista-old",
	"avvisade-new",
	"avvisade-old",
	"betalningsspec-new",
	"betalningsspec-old",
	"bevakningsreg-new",
	"bevakningsreg-old",
	"medgivande-new",
	"medgivande-old",
	"medgivandeavi-new",
	"medgivandeavi-old",
	"medgivandereg-new",
	"medgivandereg-old",
}

func TestFileNormalization(t *testing.T) {

	for _, name := range FileNormalization {
		fileContent, err := os.ReadFile(fmt.Sprintf("../tests/normalization/%s.txt", name))
		if err != nil {
			t.Errorf("Failed to read file: %s", err)
		}

		expectedContent, err := os.ReadFile(fmt.Sprintf("../tests/normalization/%s-expected.txt", name))
		if err != nil {
			t.Errorf("Failed to read file: %s", err)
		}

		isoContent := tools.BytesEnsureIso(fileContent)
		isoExpectedContent := tools.BytesEnsureIso(expectedContent)

		actualContent := seal.NormalizeContent(isoContent)

		if !tools.BytesliceEqual(actualContent, isoExpectedContent) {
			t.Errorf("Normalization failed for file %s: got/expected \r\n%s \r\n%s", name, actualContent, isoExpectedContent)
		}
	}

}
