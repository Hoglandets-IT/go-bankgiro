package send

import (
	"time"

	"github.com/hoglandets-it/go-bankgiro/seal"
)

// The type representing an outgoing Bankgiro file
type BankgiroFile struct {
	SealDate         string
	Content          string
	FormattedContent string
	Seal             seal.HmacSealer
}

// Creates a new Bankgiro file with the given content
func CreateBankgiroFile(content string) BankgiroFile {
	bgf := BankgiroFile{
		SealDate: time.Now().Format("060102"),
		Content:  content,
		Seal:     seal.HmacSealer{},
	}
	// Formats the content
	// 1. Replaces line endings with CRLF
	// 2. Removes any trailing or contained empty lines
	bgf.FormattedContent = seal.FormatContentString(bgf.Content)

	return bgf
}

// Set a different seal date for the Bankgiro file
func (bg *BankgiroFile) SetDate(date time.Time) {
	bg.SealDate = date.Format("060102")
}

// Set the key used to seal the Bankgiro file
func (bg *BankgiroFile) SetSealKey(key string) error {
	return bg.Seal.SetKey(key)
}

// Set the key used to seal the Bankgiro file as a byte array
func (bg *BankgiroFile) SetSealKeyBytes(key []byte) error {
	return bg.Seal.SetKeyBytes(key)
}

// TODO: Remove all blank/space-only rows
// TODO: Add check for BG Number (ensureBgNumberCorrect)
// TODO: Add regex \r\n[ ]*\r\n
// TODO: Replace tabulation?
