package sign_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/hoglandets-it/go-bankgiro/sign"
)

const (
	SignedBy     = "1234567890ABCDEF1234567890ABCDEF"
	SignedByKvv  = "FF365893D899291C3BF505FB3175E880"
	SignedOnDate = "240429"
)

var TEST_FILES = []string{
	"basic",
	"blank-rows",
}

// Test: Ensure error when incorrect date
// Test: Ensure error when incorrect key
// Test: Ensure error when incorrect kvv

func TestAgainstExpected(t *testing.T) {
	for _, file := range TEST_FILES {
		content, err := os.ReadFile("../tests/sealFile/" + file + ".txt")
		if err != nil {
			t.Fatal(err)
		}

		expected, err := os.ReadFile("../tests/sealFile/" + file + "-signed.txt")
		if err != nil {
			t.Fatal(err)
		}

		bgf, err := sign.CreateBankgiroFileBytes(content)
		if err != nil {
			t.Error(err)
		}

		bgf.SetSealKey(SignedBy)
		bgf.CheckKvv(SignedByKvv)
		bgf.SetSealDate(SignedOnDate)

		if !bgf.ReadyToSign() {
			t.Error("Not ready to sign")
		}

		err = bgf.Sign()
		if err != nil {
			t.Error(err)
		}

		signed := bgf.GetSignedData()
		if string(signed) != string(expected) {
			splitSigned := bytes.Split([]byte(signed), []byte("\r\n"))
			splitExpected := bytes.Split(expected, []byte("\r\n"))
			if len(splitExpected) == 1 {
				splitExpected = bytes.Split(expected, []byte("\n"))
			}
			for i, line := range splitSigned {
				if !bytes.Equal(line, splitExpected[i]) {
					t.Errorf("Signed content does not match expected content\r\nline%s: %s\r\nline%s: %s", fmt.Sprint(i), string(splitSigned[i]), fmt.Sprint(i), string(splitExpected[i]))
				}
			}
			if !bytes.Equal([]byte(signed), expected) {
				t.Errorf("Signed content does not match expected content\r\n%s\r\n%s", string(signed), string(expected))
			}
			// t.Errorf("Signed content does not match expected content\r\n%s\r\n%s", string(signed), string(signedContent))
		}

	}
}
