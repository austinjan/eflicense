package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"github.com/adtac/go-lcns"
	"io/ioutil"
	"os"
)

type option struct {
	PrivateKeyFile string
	OptionFile     string
	Output         string
	Help           bool
}

var opt option

func init() {
	const (
		defaultPrivateFile = "private.key"
		defaultOptionFile  = "setting.json"
		defaultOutput      = "license"
	)
	flag.BoolVar(&opt.Help, "help", false, "Help")
	flag.StringVar(&opt.PrivateKeyFile, "prikey", defaultPrivateFile, "private key file name.")
	flag.StringVar(&opt.Output, "out", defaultOutput, " output license file name ")
	flag.StringVar(&opt.OptionFile, "payload", defaultOptionFile, " setting json file ")
}

func getPayload(publickey *rsa.PublicKey, license string) (string, error) {
	payload, err := lcns.VerifyAndExtractPayload(publickey, license)
	if err != nil {
		// The license was probably invalid or corrupt.
		return "", err
	}

	fmt.Println("payload is: ", string(payload.([]byte))) // "some payload"
	return string(payload.([]byte)), nil
}

func main() {
	flag.Parse()

	if opt.Help || len(os.Args) < 2 {
		flag.Usage()
	}

	//fmt.Println("Reading public key...", opt.PrivateKeyFile, opt.OptionFile)
	prikey, err := lcns.ReadPrivateKey(opt.PrivateKeyFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	payload, err := ioutil.ReadFile(opt.OptionFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	licenseKeyString, err := lcns.GenerateFromPayload(prikey, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	fd, err := os.Create(opt.Output)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := fd.WriteString(licenseKeyString)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("License file create success! ", l)
	err = fd.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	return

	// licenseContent, err := ioutil.ReadFile("./license")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// publicKey, err := lcns.ReadPublicKey("./test-public.key")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// _, err = getPayload(publicKey, string(licenseContent))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// return

}
