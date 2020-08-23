package decrypt

import (
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"os"
)

func remove(path string) {
	e := os.Remove(path)
	if e != nil {
	}
}

// Decrypt ...
func Decrypt(filePathIn string, filePathOut string, keyAes, keyHmac []byte, aes cipher.Block, iv []byte) error {
	const bufferSize int = 4096
	inFile, err := os.Open(filePathIn)
	if err != nil {
		return err
	}
	defer inFile.Close()
	outFile, err := os.Create(filePathOut)
	if err != nil {
		return err
	}
	defer outFile.Close()
	ctr := cipher.NewCTR(aes, iv)
	hmac := hmac.New(sha256.New, keyHmac)
	buf := make([]byte, bufferSize)
	for {
		n, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		outBuf := make([]byte, n)
		ctr.XORKeyStream(outBuf, buf[:n])
		hmac.Write(outBuf)
		outFile.Write(outBuf)
	}
	inFile.Close()
	remove(filePathIn)
	return nil
}
