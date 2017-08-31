// utils package provides commands to handle encryption and signing
package utils

import (
	"bytes"
	"encoding/hex"
	"strings"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"golang.org/x/crypto/openpgp"
)

// ProcessRelease  read publickey from given channel and verify the given signature
func ProcessRelease(channel database.Channel, release database.Release, signature string, signed string) (bool, error) {
	keyRingFile := strings.NewReader(channel.PublicKey)
	signedFile := strings.NewReader(signed)
	signatureFile := strings.NewReader(signature)

	keyRing, err := openpgp.ReadArmoredKeyRing(keyRingFile)
	if err != nil {
		return false, err
	}
	_, err = openpgp.CheckArmoredDetachedSignature(keyRing, signedFile, signatureFile)
	if err != nil {
		return false, err
	} else {
		return true, err
	}
}

// Sign create read publickey from given channel and sign the releases stated
func Sign(channel database.Channel, release database.Release, signed string) (string, error) {
	keyRingFile := strings.NewReader(channel.PrivateKey)

	signer, err := openpgp.ReadArmoredKeyRing(keyRingFile)
	if err != nil {
		return "", err
	}

	message := strings.NewReader(signed)
	w := new(bytes.Buffer)

	if err = openpgp.ArmoredDetachSign(w, signer[0], message, nil); err != nil {
		return "", err
	}

	return w.String(), nil
}

// GetFingerprint return the fingerprint from a channel in uppercase letters
func GetFingerprint(channel database.Channel) (string, error) {
	keyRingFile := strings.NewReader(channel.PublicKey)

	keyRing, err := openpgp.ReadArmoredKeyRing(keyRingFile)
	if err != nil {
		return "", err
	}
	fingerprint := strings.ToUpper(hex.EncodeToString(keyRing[0].PrimaryKey.Fingerprint[:]))
	return fingerprint, nil
}

// CheckPGPKey check if the given key follows the basic gpg terms
func CheckPGPKey(key string) error {
	keyRingFile := strings.NewReader(key)
	_, err := openpgp.ReadArmoredKeyRing(keyRingFile)
	return err
}
