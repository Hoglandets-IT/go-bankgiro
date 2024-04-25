package seal_test

import (
	"bytes"
	"encoding/hex"
	"slices"
	"testing"

	"github.com/hoglandets-it/go-bankgiro/seal"
	"github.com/hoglandets-it/go-bankgiro/tools"
)

const (
	NormalizationTestCharset = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	OutOfBoundsReplacement   = "Ã"
	OutOfBoundsLowerHex      = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
)

var BasicNormalization map[string]string = map[string]string{
	"test":                       "test",
	"1234567890 0":               "1234567890 0",
	"abcdefghijklmnopqrstuvwxyz": "abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"åäöÅÄÖ":                     "}{|][\\",
	"åäö\r\nÅÄÖ\r\n\r\nmultiline-string__": "}{|][\\multiline-string__",
}

// É Ä Ö Å Ü é ä ö å ü
// @ [ \ ] ^ ` { | } ~

func TestBasicNormalization(t *testing.T) {
	for input, expected := range BasicNormalization {
		isoString := tools.StringEnsureIso(input)
		actual := seal.NormalizeContentString(isoString)
		if string(actual) != expected {
			t.Errorf("Normalization failed: got %s expected %s", actual, expected)
		}
	}
}

func TestUnescapedCharacters(t *testing.T) {
	charset := tools.StringEnsureIso(NormalizationTestCharset)
	actual := seal.NormalizeContentString(charset)
	if string(actual) != charset {
		t.Errorf("Normalization failed: got %s expected %s", actual, charset)
	}
}

func TestOutOfBounds(t *testing.T) {
	oobLowerString, err := hex.DecodeString(OutOfBoundsLowerHex)
	if err != nil {
		t.Errorf("Failed to encode out of bounds replacement: %s", err)
	}
	oobLen := len(oobLowerString)

	oobLowerHex := tools.BytesEnsureIso(oobLowerString)
	actual := seal.NormalizeContent(oobLowerHex)

	if slices.Equal(actual, bytes.Repeat([]byte(OutOfBoundsReplacement), oobLen)) {
		t.Errorf("Normalization failed: got %s expected %s", actual, OutOfBoundsLowerHex)
	}
}
