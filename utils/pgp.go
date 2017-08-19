package utils

import (
	"encoding/hex"
	"log"
	"os"
	"strings"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"golang.org/x/crypto/openpgp"
)

func ProcessRelease(release database.Release) bool {

	ReleaseDir := "static/releases/tmp"
	SignatureDir := "static/signatures/tmp.asc"

	signed, err := os.Open(ReleaseDir)
	if err != nil {
		return false
	}
	defer signed.Close()

	signature, err := os.Open(SignatureDir)
	if err != nil {
		return false
	}
	defer signature.Close()

	pub := "static/channels/public/" + release.Channel + ".asc"
	keyRingFile, err := os.Open(pub)
	if err != nil {
		log.Println(err)
		return false
	}
	defer keyRingFile.Close()
	keyRing, err := openpgp.ReadArmoredKeyRing(keyRingFile)
	if err != nil {
		log.Println(err)
		return false
	}

	isvalid := verify(keyRing, ReleaseDir, SignatureDir)

	return isvalid
}

func verify(keyRing openpgp.EntityList, signed, signature string) bool {
	signedFile, err := os.Open(signed)
	if err != nil {
		log.Println(err)
		return false
	}
	defer signedFile.Close()
	signatureFile, err := os.Open(signature)
	if err != nil {
		log.Println(err)
		return false
	}
	defer signatureFile.Close()
	_, err = openpgp.CheckArmoredDetachedSignature(keyRing, signedFile, signatureFile)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func Sign(PrivateKeyFileName string, ReleaseFileName string, SignatureFileName string) (err error) {

	keyRingFile, err := os.Open(PrivateKeyFileName)
	if err != nil {
		log.Println(err)
		return err
	}
	defer keyRingFile.Close()

	signer, err := openpgp.ReadArmoredKeyRing(keyRingFile)
	if err != nil {
		log.Println(err)
		return err
	}

	var message *os.File
	if message, err = os.Open(ReleaseFileName); err != nil {
		return err
	}
	defer message.Close()

	var w *os.File
	if w, err = os.Create(SignatureFileName); err != nil {
		return err
	}
	defer w.Close()

	if err = openpgp.ArmoredDetachSign(w, signer[0], message, nil); err != nil {
		return err
	}

	return nil
}

func GetFingerprint(channel string) (string, error) {
	pub := "static/channels/public/" + channel + ".asc"
	keyRingFile, err := os.Open(pub)
	if err != nil {
		return "", err
	}
	defer keyRingFile.Close()
	keyRing, err := openpgp.ReadArmoredKeyRing(keyRingFile)
	if err != nil {
		return "", err
	}
	fingerprint := strings.ToUpper(hex.EncodeToString(keyRing[0].PrimaryKey.Fingerprint[:]))
	return fingerprint, nil
}
