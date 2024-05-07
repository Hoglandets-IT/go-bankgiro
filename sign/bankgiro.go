package sign

import (
	"fmt"

	"github.com/hoglandets-it/go-bankgiro/seal"
	"github.com/hoglandets-it/go-bankgiro/tools"
)

// The type representing an outgoing Bankgiro file
type BankgiroFile struct {
	Content          string
	FormattedContent string
	Seal             seal.HmacSealer
}

// Creates a new Bankgiro file with the given content
func CreateBankgiroFile(content string) BankgiroFile {
	bgf := BankgiroFile{
		Content: content,
		Seal:    seal.HmacSealer{},
	}
	// Formats the content
	// 1. Replaces line endings with CRLF
	// 2. Removes any trailing or contained empty lines
	bgf.FormattedContent = seal.FormatContentString(bgf.Content)

	return bgf
}

// Creates a new Bankgiro file with the given byteslice content
func CreateBankgiroFileBytes(content []byte) (BankgiroFile, error) {
	isoString, err := tools.BytesToIsoString(content)
	if err != nil {
		return BankgiroFile{}, err
	}

	bgf := BankgiroFile{
		Content: isoString,
		Seal:    seal.HmacSealer{},
	}

	bgf.FormattedContent = seal.FormatContentString(bgf.Content)
	bgf.Seal.SetData(bgf.FormattedContent)

	return bgf, nil
}

// Set the key used to seal the Bankgiro file
func (bg *BankgiroFile) SetSealKey(key string) error {
	return bg.Seal.SetKey(key)
}

// Set the key used to seal the Bankgiro file as a byte array
func (bg *BankgiroFile) SetSealKeyBytes(key []byte) error {
	return bg.Seal.SetKeyBytes(key)
}

// Check the KVV value against a provided value
func (bg *BankgiroFile) CheckKvv(kvv string) error {
	return bg.Seal.CheckKvv(kvv)
}

// Set a custom Seal Date
func (bg *BankgiroFile) SetSealDate(date string) {
	bg.Seal.SetSealDate(date)
}

// Check if the file is ready to be signed
func (bg *BankgiroFile) ReadyToSign() bool {
	return bg.Seal.Key != nil && bg.Seal.KeyVer != nil && bg.FormattedContent != "" && bg.Seal.Validate() == nil
}

func (bg *BankgiroFile) Sign() error {
	if !bg.ReadyToSign() {
		return fmt.Errorf("not ready to sign - error")
	}

	err := bg.Seal.Calculate()
	if err != nil {
		return err
	}

	return nil
}

func (bg *BankgiroFile) GetSignedData() string {
	return bg.Seal.GetSignedContent()
}

// TODO: Remove all blank/space-only rows
// TODO: Add check for BG Number (ensureBgNumberCorrect)
// TODO: Add regex \r\n[ ]*\r\n
// TODO: Replace tabulation?
