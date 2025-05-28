package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cast"
)

func invertHexBytes(hexStr string) string {
	bytes, _ := hex.DecodeString(hexStr)
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return hex.EncodeToString(bytes)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./pixlic <serial number>")
		os.Exit(1)
	}

	arg, err := cast.ToInt64E(os.Args[1])
	if err != nil {
		fmt.Println("Invalid serial number provided!")
		fmt.Println("Usage: ./pixlic <serial number>")
		os.Exit(1)
	}

	serial := invertHexBytes(strconv.FormatInt(arg, 16))
	license := "39000000" // AES+DES+UR License
	data, err := hex.DecodeString(license + serial)
	if err != nil {
		panic(err)
	}
	hash := fmt.Sprintf("%x", md5.Sum(data))

	print("Here's your PIX UR License: ")
	for idx := 0; idx < 16; idx += 4 {
		part := hash[(idx * 2):((idx + 4) * 2)]
		fmt.Printf("0x%s ", invertHexBytes(part))
	}
	println()
}
