package encrypt

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

// Encrypt ...
func Encrypt(filePathIn string, filePathOut string, keyAes, keyHmac []byte, aes cipher.Block, iv []byte) error {
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

	//fmt.Println(reflect.TypeOf(iv))
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
		if err == io.EOF {
			break
		}
	}
	inFile.Close()
	remove(filePathIn)
	return nil
}
