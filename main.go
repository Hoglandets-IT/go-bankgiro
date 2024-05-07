package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "seal",
		Usage: "apply a seal to a given file",
		Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	// content, err := os.ReadFile("tests/sealFile/00-basic.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// signedContent, err := os.ReadFile("tests/sealFile/00-basic-signed.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// bgf, err := sign.CreateBankgiroFileBytes(content)
	// if err != nil {
	// 	panic(err)
	// }

	// bgf.SetSealKey("1234567890ABCDEF1234567890ABCDEF")
	// if bgf.ReadyToSign() {
	// 	err := bgf.Sign()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	signed := bgf.GetSignedData()
	// 	fmt.Println(signed)
	// 	fmt.Println(string(signedContent))
	// } else {
	// 	panic("Not ready to sign")
	// }
}
