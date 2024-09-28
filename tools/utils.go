package tools

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func HandleRecover() {
	if r := recover(); r != nil {
		Logger("error", fmt.Sprintf("%v", r))
		Logger("info", "Press Enter to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func RandomNumber(min, max int) int {
	return rand.Intn(max-min) + min
}

func GetTextAfterKey(urlData, key string) (string, error) {
	// Temukan lokasi key
	keyIndex := strings.Index(urlData, key)
	if keyIndex == -1 {
		return "", fmt.Errorf("key %s tidak ditemukan", key)
	}

	// Ambil substring setelah key
	startIndex := keyIndex + len(key)
	endIndex := strings.Index(urlData[startIndex:], "&")
	if endIndex == -1 {
		return urlData[startIndex:], nil
	}

	return urlData[startIndex : startIndex+endIndex], nil
}
