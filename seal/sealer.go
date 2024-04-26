package seal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

// Key: the hex-decoded HMAC key used to seal the file
// KeyVer (KVV, KeyVerificationValue) is the the value used to verify the key, obtained by sealing the string "00000000"
// HashFunc: the hash function used to calculate the HMAC seal, default is sha256
type HmacSealer struct {
	Key            []byte
	KeyVer         []byte
	Hash           func() hash.Hash
	Mac            []byte
	OriginalData   []byte
	FormattedData  []byte
	NormalizedData []byte
}

// Set the hash function used to hash the file contents
func (hm *HmacSealer) SetHashFunction(hashFunc func() hash.Hash) {
	hm.Hash = hashFunc
}

// Set the key used to seal the file
func (hm *HmacSealer) SetKey(key string) error {
	return hm.SetKeyBytes([]byte(key))
}

// Set the key used to seal the file as a byte array
func (hm *HmacSealer) SetKeyBytes(key []byte) error {
	hexLen, err := hex.Decode(hm.Key, key)
	if err != nil {
		return fmt.Errorf("invalid key provided: %v", err)
	}

	if hexLen != 16 {
		return fmt.Errorf("invalid key length: %d, expected 16", hexLen)
	}

	return nil
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

// Update the values of the Formatted/Normalized data fields
// Formatted will ensure the correct line endings are present
// Normalized will normalize the content for HMAC signature
func (hm *HmacSealer) UpdateFormatted() {
	hm.FormattedData = FormatContent(hm.OriginalData)
	hm.NormalizedData = NormalizeContent(hm.OriginalData)
}

// Replace the data with a given string
func (hm *HmacSealer) SetData(data string) {
	hm.OriginalData = []byte(data)
	hm.UpdateFormatted()
}

// Append given string to the data
func (hm *HmacSealer) AddData(data string) {
	hm.OriginalData = append(hm.OriginalData, []byte(data)...)
	hm.UpdateFormatted()
}

// Replace the data with a given byte slice
func (hm *HmacSealer) SetDataBytes(data []byte) {
	hm.OriginalData = data
	hm.UpdateFormatted()
}

// Append given byte slice to the data
func (hm *HmacSealer) AddDataBytes(data []byte) {
	hm.OriginalData = append(hm.OriginalData, data...)
	hm.UpdateFormatted()
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
