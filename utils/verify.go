package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/xmarcoied/go-updater/model"
	"golang.org/x/crypto/openpgp"
)

func ProcessRelease(release model.Release) bool {
	ReleaseJSON, _ := json.Marshal(release)
	ReleaseDir := "static/releases/" + strconv.Itoa(int(release.ID))
	SignatureDir := "static/signatures/" + strconv.Itoa(int(release.ID)) + ".asc"

	ioutil.WriteFile(ReleaseDir, ReleaseJSON, 0644)
	ioutil.WriteFile(SignatureDir, []byte(release.Signature), 0644)

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
