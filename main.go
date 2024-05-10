package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hoglandets-it/go-bankgiro/shell"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "go-bankgiro",
		HelpName:    "go-bankgiro",
		Description: "A tool to seal and validate Bankgiro files with HMAC",
		Commands: []*cli.Command{
			{
				Name:      "seal",
				Aliases:   []string{"s"},
				Usage:     "seal a file with a given key",
				Args:      true,
				ArgsUsage: " [file-to-sign]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Required: true,
						Usage:    "key to seal the file with",
						EnvVars:  []string{"BG_SEAL_KEY"},
					},
					&cli.StringFlag{
						Name:     "kvv",
						Aliases:  []string{"v"},
						Required: false,
						Usage:    "kvv to check the seal with (optional)",
						EnvVars:  []string{"BG_SEAL_KVV"},
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Required: false,
						Usage:    "output file, default is [file-to-sign]-signed",
						EnvVars:  []string{"BG_SEAL_OUTPUT"},
					},
					&cli.StringFlag{
						Name:        "overwrite",
						Aliases:     []string{"f"},
						Required:    false,
						DefaultText: "false",
						Usage:       "overwrite the output file if it exists",
						EnvVars:     []string{"BG_SEAL_OVERWRITE"},
					},
				},
				Action: func(c *cli.Context) error {
					err := shell.ParseVars(c)
					if err != nil {
						return err
					}

					fmt.Printf("Parameters, valid, starting seal on file %s \r\n", c.Args().First())

					return shell.SealFile(c)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
