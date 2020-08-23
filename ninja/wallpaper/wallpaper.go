package wallpaper

import (
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724947.aspx
const (
	spiGetDeskWallpaper = 0x0073
	spiSetDeskWallpaper = 0x0014

	uiParam = 0x0000

	spifUpdateINIFile = 0x01
	spifSendChange    = 0x02
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724947.aspx
var (
	user32               = syscall.NewLazyDLL("user32.dll")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

// Get returns the current wallpaper.
func Get() (string, error) {
	// the maximum length of a windows path is 256 utf16 characters
	var filename [256]uint16
	systemParametersInfo.Call(
		uintptr(spiGetDeskWallpaper),
		uintptr(cap(filename)),
		// the memory address of the first byte of the array
		uintptr(unsafe.Pointer(&filename[0])),
		uintptr(0),
	)
	return strings.Trim(string(utf16.Decode(filename[:])), "\x00"), nil
}

// SetFromFile sets the wallpaper for the current user.
func SetFromFile(filename string) error {
	filenameUTF16, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}

	systemParametersInfo.Call(
		uintptr(spiSetDeskWallpaper),
		uintptr(uiParam),
		uintptr(unsafe.Pointer(filenameUTF16)),
		uintptr(spifUpdateINIFile|spifSendChange),
	)
	return nil
}

func downloadImage(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
	}
	defer response.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func Change(url string) {
	downloadImage(url, "axe.jpg")
	cwd, err := os.Getwd()
	if err != nil {
	}
	SetFromFile(cwd + "\\axe.jpg")
}
