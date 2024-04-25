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
	keyLen, err := hex.Decode(hm.Key, key)
	if err != nil {
		return fmt.Errorf("invalid key provided: %v", err)
	}

	if keyLen != 32 {
		return fmt.Errorf("invalid key length: %d, expected 32", keyLen)
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
