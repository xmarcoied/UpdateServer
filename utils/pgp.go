package utils

import (
	"bytes"
	"encoding/hex"
	"strings"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"golang.org/x/crypto/openpgp"
)

func ProcessRelease(release database.Release, signature string, signed string) (bool, error) {
	channel := core.GetChannel(release.Channel)
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

func Sign(release database.Release, signed string) (string, error) {
	channel := core.GetChannel(release.Channel)
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

func GetFingerprint(channelname string) (string, error) {
	channel := core.GetChannel(channelname)
	keyRingFile := strings.NewReader(channel.PublicKey)

	keyRing, err := openpgp.ReadArmoredKeyRing(keyRingFile)
	if err != nil {
		return "", err
	}
	fingerprint := strings.ToUpper(hex.EncodeToString(keyRing[0].PrimaryKey.Fingerprint[:]))
	return fingerprint, nil
}
