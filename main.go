package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func invertHexBytes(hexStr string) string {
	bytes, _ := hex.DecodeString(hexStr)
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return hex.EncodeToString(bytes)
}

// AES+DES+UR License
const (
	PIX_515_UR = "39000000"
	PIX_525_UR = "38040000"
)

var licenseFutures = map[string]string{
	"PIX 501": "2A800100", // needs testing
	"PIX 506": "????????",
	"PIX 515": "39000000",
	"PIX 520": "????????",
	"PIX 525": "38040000",
	"PIX 535": "39000000", //needs testing
}

func main() {

	cliSerial := flag.Int64("serial", 0, "Specify a PIX serial number")
	cliModel := flag.String("model", "PIX 515", "Specify a PIX Model or set \"dev\" for development")
	cliListModels := flag.Bool("list", false, "List all available PIX models")

	devLicFuture := flag.String("dev-future", "", "(Test) Pass custom license future value")
	devPkey := flag.String("dev-pkey", "", "(Test) Pass a license of an existing UR box to match against")
	devBruteForce := flag.Bool("dev-bruteforce", false, "(Test) Bruteforce a future value based on license and SN")

	if len(os.Args) < 2 {
		flag.Usage()
	}

	if *cliListModels {
		println("Supported PIX Models:")
		println("  - PIX 515 (Includes 515e)")
		println("  - PIX 525")
		println("\nUntested Models:")
		println("  - PIX 501")
		println("  - PIX 535")
		println("\nUnsupported Models:")
		println("  - PIX 506 (Includes 506e)")
		println("  - PIX 520")
		os.Exit(1)
	}

	arg := *cliSerial
	if arg == 0 {
		println("Invalid serial number provided!")
		flag.Usage()
	}

	pixModel := *cliModel
	future, ok := licenseFutures[pixModel]
	if !ok && pixModel != "dev" {
		println("Invalid pix model provided!")
		flag.Usage()
	}

	if pixModel == "dev" && *devLicFuture != "" {
		future = *devLicFuture
	}

	serial := invertHexBytes(strconv.FormatInt(arg, 16))

	if *devBruteForce {
		bruteForceGenerate(pixModel, serial, *devPkey)
	} else {
		data, err := hex.DecodeString(future + serial)
		if err != nil {
			panic(err)
		}

		hash := fmt.Sprintf("%x", md5.Sum(data))
		generateOneTime(pixModel, future, hash, *devPkey)
	}

}

func generateOneTime(pixModel, future, hash, devKey string) {
	pkey := ""
	print("Here's your PIX (UR) License: ")
	for idx := 0; idx < 16; idx += 4 {
		part := invertHexBytes(hash[(idx * 2):((idx + 4) * 2)])
		pkey += "0x" + part
	}
	println(pkey)

	pkey = strings.TrimRight(pkey, " ")
	if devKey != "" && devKey == pkey {
		println("DEBUG: LICENSE KEY MATCH!")
		fmt.Printf("pkey: %s", pkey)
		fmt.Printf("devKey: %s", devKey)
		fmt.Printf("future: %s", future)
		fmt.Printf("model: %s", pixModel)
	}
}

func bruteForceGenerate(pixModel, serial, devKey string) {
	for idx := range 0x5F5E0FF {
		genFuture := invertHexBytes(fmt.Sprintf("%08X", idx))
		data, err := hex.DecodeString(genFuture + serial)
		if err != nil {
			continue
		}

		hash := fmt.Sprintf("%x", md5.Sum(data))
		pkey := ""
		for idx := 0; idx < 16; idx += 4 {
			part := invertHexBytes(hash[(idx * 2):((idx + 4) * 2)])
			pkey += "0x" + part
		}
		pkey = strings.TrimRight(pkey, " ")

		if devKey != "" && devKey == pkey {
			println("DEBUG: MAGIC FUTURE FOUND!")
			fmt.Printf("future: %s", genFuture)
			fmt.Printf("pkey: %s", pkey)
			fmt.Printf("devKey: %s", devKey)
			fmt.Printf("model: %s", pixModel)
		}
	}
}
