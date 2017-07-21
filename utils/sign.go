package utils

import (
	"log"
	"os"

	"golang.org/x/crypto/openpgp"
)

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
