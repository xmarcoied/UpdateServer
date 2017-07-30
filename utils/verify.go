package utils

import (
	"log"
	"os"
	"strconv"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/model"
	"golang.org/x/crypto/openpgp"
)

func ProcessRelease(release model.Release) bool {

	ReleaseDir := "static/releases/" + strconv.Itoa(int(release.ID))
	SignatureDir := "static/signatures/" + strconv.Itoa(int(release.ID)) + ".asc"

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
