package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/ripemd160"
)

func CreatePublicAddress(privKeyInt *big.Int) string {
	privKeyHex := fmt.Sprintf("%064x", privKeyInt)
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)
	compressedPubKey := privKey.PubKey().SerializeCompressed()
	pubKeyHash := Hash160(compressedPubKey)
	address := EncodeAddress(pubKeyHash, &chaincfg.MainNetParams)

	return address
}

func Hash160(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	sha256Hash := h.Sum(nil)

	r := ripemd160.New()
	r.Write(sha256Hash)
	return r.Sum(nil)
}

func EncodeAddress(pubKeyHash []byte, params *chaincfg.Params) string {
	versionedPayload := append([]byte{params.PubKeyHashAddrID}, pubKeyHash...)
	checksum := DoubleSha256(versionedPayload)[:4]
	fullPayload := append(versionedPayload, checksum...)
	return Base58Encode(fullPayload)
}

func DoubleSha256(b []byte) []byte {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return second[:]
}
