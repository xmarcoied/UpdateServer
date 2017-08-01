package utils

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp"
)

func Sign(PrivateKeyFileName string, ReleaseFileName string, SignatureFileName string) (err error) {

	pgpend := `-----END PGP PRIVATE KEY BLOCK-----`
	scanner := bufio.NewScanner(os.Stdin)
	var output []string
	for scanner.Scan() {
		output = append(output, scanner.Text())
		if strings.Contains(scanner.Text(), pgpend) {
			break
		}
	}
	privateKey := strings.Join(output, "\n")
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}

	ioutil.WriteFile(PrivateKeyFileName, []byte(privateKey), 0644)

	keyRingFile, err := os.Open(PrivateKeyFileName)

	defer os.Remove(PrivateKeyFileName)

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
		log.Println(err)
		return err
	}
	defer message.Close()

	var w *os.File
	if w, err = os.Create(SignatureFileName); err != nil {
		log.Println(err)
		return err
	}
	defer w.Close()

	if err = openpgp.ArmoredDetachSign(w, signer[0], message, nil); err != nil {
		return err
	}

	return nil
}
