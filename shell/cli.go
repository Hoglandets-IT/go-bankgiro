package shell

import (
	"fmt"
	"os"

	"github.com/hoglandets-it/go-bankgiro/sign"
	"github.com/urfave/cli/v2"
)

func ParseVars(c *cli.Context) error {
	key := c.String("key")
	kvv := c.String("kvv")

	if key == "" {
		return cli.Exit("key is required", 1)
	}

	if kvv == "" {
		fmt.Println("No KVV provided, no validation will be done on key")
	}

	if c.Args().Len() == 0 {
		return cli.Exit("file-to-sign is required", 1)
	}

	file := c.Args().First()
	if file == "" {
		return cli.Exit("file-to-sign is required", 1)
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("%s does not exist", file), 1)
	}

	output := c.String("output")
	if output == "" {
		output = fmt.Sprintf("%s-signed", file)
	}

	if _, err := os.Stat(output); err == nil {
		overwrite := c.Bool("overwrite")
		if !overwrite {
			return cli.Exit(fmt.Sprintf("%s already exists, use -f to overwrite", output), 1)
		}
	}

	fmt.Println("Output set to ", output)

	return nil
}

func SealFile(c *cli.Context) error {
	file, err := os.ReadFile(c.Args().First())
	if err != nil {
		return err
	}

	bgFile, err := sign.CreateBankgiroFileBytes(file)
	if err != nil {
		return err
	}

	key := c.String("key")
	err = bgFile.SetSealKey(key)
	if err != nil {
		return err
	}

	kvv := c.String("kvv")
	if kvv != "" {
		err = bgFile.CheckKvv(kvv)
		if err != nil {
			return err
		}
	}

	err = bgFile.Sign()
	if err != nil {
		return err
	}

	fmt.Println("File signed successfully")

	content := bgFile.GetSignedData()
	output := c.String("output")
	if output == "" {
		output = fmt.Sprintf("%s-signed", c.Args().First())
	}

	fmt.Println("File saved to", output)

	return os.WriteFile(output, []byte(content), 0644)
}
