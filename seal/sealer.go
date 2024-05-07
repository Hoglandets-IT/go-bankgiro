package seal

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
	"time"

	"github.com/hoglandets-it/go-bankgiro/tools"
)

// Key: the hex-decoded HMAC key used to seal the file
// KeyVer (KVV, KeyVerificationValue) is the the value used to verify the key, obtained by sealing the string "00000000"
// HashFunc: the hash function used to calculate the HMAC seal, default is sha256
type HmacSealer struct {
	Key            []byte
	KeyVer         []byte
	Hash           func() hash.Hash
	Mac            []byte
	SealDate       string
	OriginalData   []byte
	PrefixedData   []byte
	FormattedData  []byte
	NormalizedData []byte
}

func (hm *HmacSealer) EnsureNoSignature() error {
	if hm.Mac != nil {
		return fmt.Errorf("cannot change settings after signature has been calculated")
	}

	return nil
}

// Set the hash function used to hash the file contents
func (hm *HmacSealer) SetHashFunction(hashFunc func() hash.Hash) error {
	if err := hm.EnsureNoSignature(); err != nil {
		return err
	}

	hm.Hash = hashFunc

	return nil
}

// Set the key used to seal the file
func (hm *HmacSealer) SetKey(key string) error {
	return hm.SetKeyBytes([]byte(key))
}

// Set the key used to seal the file as a byte array
func (hm *HmacSealer) SetKeyBytes(key []byte) error {
	if err := hm.EnsureNoSignature(); err != nil {
		return err
	}

	hm.Key = make([]byte, 16)
	hexLen, err := hex.Decode(hm.Key, key)
	if err != nil {
		return fmt.Errorf("invalid key provided: %v", err)
	}

	if hexLen != 16 {
		return fmt.Errorf("invalid key length: %d, expected 16", hexLen)
	}

	return hm.GenerateKvv()
}

// Generate and store the KVV Value
func (hm *HmacSealer) GenerateKvv() error {
	if hm.Key == nil {
		return fmt.Errorf("key not set")
	}

	kvvCalcValue := []byte("00000000")

	if hm.Hash == nil {
		hm.SetHashFunction(sha256.New)
	}

	hmhash := hmac.New(hm.Hash, hm.Key)
	hmhash.Write(kvvCalcValue)

	hm.KeyVer = hmhash.Sum([]byte{})

	return nil
}

// Check the KVV value against a provided value
func (hm *HmacSealer) CheckKvv(kvv string) error {
	if hm.KeyVer == nil {
		return fmt.Errorf("key verification value not set")
	}

	storedKvvHex := strings.ToUpper(hex.EncodeToString(hm.KeyVer))
	providedKvvHex := strings.ToUpper(kvv)

	if len(kvv) == 64 {
		if providedKvvHex == storedKvvHex {
			return nil
		}
		return fmt.Errorf("provided and calculated kvv do not match:\r\nProvided:   %s\r\nCalculated: %s", providedKvvHex, storedKvvHex)
	}

	if len(kvv) == 32 {
		if providedKvvHex == storedKvvHex[0:32] {
			return nil
		}
		return fmt.Errorf("provided and calculated kvv do not match:\r\nProvided:   %s\r\nCalculated: %s", providedKvvHex, storedKvvHex[0:32])
	}

	return fmt.Errorf("invalid KVV length: %d, expected 32 or 64", len(kvv))
}

func PrefixContent(content []byte, sealDate string) []byte {
	if len(content) == 0 || bytes.HasPrefix(content, []byte("00")) {
		return content
	}

	prefix := fmt.Sprintf("00%sHMAC%s%s", sealDate, strings.Repeat(" ", 68), "\r\n")

	return append([]byte(prefix), content...)
}

// Update the values of the Formatted/Normalized data fields
// Formatted will ensure the correct line endings are present
// Normalized will normalize the content for HMAC signature
func (hm *HmacSealer) UpdateFormatted() error {
	if err := hm.EnsureNoSignature(); err != nil {
		return err
	}

	if hm.SealDate == "" {
		hm.SealDate = time.Now().Format("060102")
	}
	hm.PrefixedData = PrefixContent(hm.OriginalData, hm.SealDate)
	hm.FormattedData = FormatContent(hm.PrefixedData)
	hm.NormalizedData = NormalizeContent(hm.PrefixedData)

	return nil
}

// Replace the data with a given string
func (hm *HmacSealer) SetData(data string) error {
	if err := hm.EnsureNoSignature(); err != nil {
		return err
	}
	hm.OriginalData = []byte(data)
	return hm.UpdateFormatted()
}

// Append given string to the data
func (hm *HmacSealer) AddData(data string) error {
	if err := hm.EnsureNoSignature(); err != nil {
		return err
	}

	hm.OriginalData = append(hm.OriginalData, []byte(data)...)
	return hm.UpdateFormatted()
}

// Replace the data with a given byte slice
func (hm *HmacSealer) SetDataBytes(data []byte) error {
	if err := hm.EnsureNoSignature(); err != nil {
		return err
	}
	hm.OriginalData = data
	return hm.UpdateFormatted()
}

// Append given byte slice to the data
func (hm *HmacSealer) AddDataBytes(data []byte) error {
	if err := hm.EnsureNoSignature(); err != nil {
		return err
	}
	hm.OriginalData = append(hm.OriginalData, data...)
	return hm.UpdateFormatted()
}

// Set a custom seal date
func (hm *HmacSealer) SetSealDate(date string) error {
	if err := hm.EnsureNoSignature(); err != nil {
		return err
	}

	hm.SealDate = date
	return hm.UpdateFormatted()
}

// Validate that everything is in place to calculate the HMAC seal
func (hm *HmacSealer) Validate() (err error) {
	if len(hm.OriginalData) == 0 {
		return fmt.Errorf("verification failed: no data present to be signed")
	}

	if len(hm.Key) == 0 {
		return fmt.Errorf("verification failed: no key present to sign data")
	}

	return
}

// Calculate the HMAC seal MAC for the data
func (hm *HmacSealer) Calculate() (err error) {
	if err = hm.Validate(); err != nil {
		return err
	}

	if hm.Hash == nil {
		hm.SetHashFunction(sha256.New)
	}

	hmhash := hmac.New(hm.Hash, hm.Key)
	hmhash.Write(hm.NormalizedData)

	hm.Mac = hmhash.Sum([]byte{})

	return
}

// Get the full MAC as a hex string (64 characters)
func (hm *HmacSealer) GetMac() string {
	return strings.ToUpper(hex.EncodeToString(hm.Mac))
}

// Get the calculated HMAC seal MAC as a hex string (32 characters) formatted for use in BG files
func (hm *HmacSealer) GetMacBgFormat() string {
	return hm.GetMac()[0:32]
}

// Get the calculated HMAC seal MAC as a bytestring
func (hm *HmacSealer) GetMacBytes() []byte {
	return hm.Mac
}

// Get the KVV value as a hex string (64 characters)
func (hm *HmacSealer) GetKvv() string {
	return strings.ToUpper(hex.EncodeToString(hm.KeyVer))
}

// Get the KVV value as a hex string (32 characters) formatted for use in BG files
func (hm *HmacSealer) GetKvvBgFormat() string {
	return hm.GetKvv()[0:32]
}

func (hm *HmacSealer) GetSignedContent() string {
	signature := fmt.Sprintf(
		"99%s%s%s%s",
		time.Now().Format("060102"),
		hm.GetKvvBgFormat(),
		hm.GetMacBgFormat(),
		strings.Repeat(" ", 8),
	)

	return tools.StringEnsureIso(
		strings.Join(
			[]string{
				string(hm.PrefixedData),
				signature,
			},
			"\r\n",
		),
	)
}
