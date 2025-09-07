package crypto

import (
	"math/big"
)

func PrivateKeyToWIF(privKey *big.Int) string {
	prefix := []byte{0x80}
	privKeyBytes := privKey.FillBytes(make([]byte, 32))
	privKeyBytes = append(privKeyBytes, 0x01)
	extendedKey := append(prefix, privKeyBytes...)
	checksum := DoubleSha256(extendedKey)[:4]
	fullKey := append(extendedKey, checksum...)
	wif := Base58Encode(fullKey)
	return wif
}
