package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/schollz/cowyo/encrypt"
	"github.com/urfave/cli"
)

var debug bool
var version string

func main() {
	run()
}

func run() error {
	var passphrase, page, server string
	var encryptFlag, store, name, binary bool
	app := cli.NewApp()
	app.Version = version
	app.Compiled = time.Now()
	app.Name = "cowyodel"
	app.Usage = "upload/download encrypted/unencrypted text/binary to cowyo.com"
	app.UsageText = `Upload a file:
		cowyodel upload README.md
		cat README.md | cowyodel upload
   
	 Download a file:
		cowyodel download 2-adoring-thompson

	 Persist (and don't delete after first access):
		cowyodel upload --store FILE

   Specify filename:
		cowyodel upload --name README.md

   Client-side encryption:
		cowyodel upload --encrypt README.md

	 Binary-file uploading/downloading:
		cowyodel upload --binary --name image.jpg
		cowyodel download image.jpg`
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server",
			Value:       "https://cowyo.com",
			Usage:       "cowyo server to use",
			Destination: &server,
		},

		cli.BoolFlag{
			Name:        "debug",
			Usage:       "debug mode",
			Destination: &debug,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "upload",
			Aliases: []string{"u"},
			Usage:   "upload document",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "encrypt, e",
					Usage:       "encrypt using passphrase",
					Destination: &encryptFlag,
				},
				cli.BoolFlag{
					Name:        "store, s",
					Usage:       "store and persist after reading",
					Destination: &store,
				},
				cli.BoolFlag{
					Name:        "name, n",
					Usage:       "use name of file",
					Destination: &name,
				},
				cli.StringFlag{
					Name:        "passphrase, a",
					Usage:       "passphrase to use for encryption",
					Destination: &passphrase,
				},
				cli.BoolFlag{
					Name:        "binary, b",
					Usage:       "binary mode (Gzip + Base64 encoding)",
					Destination: &binary,
				},
			},
			Action: func(c *cli.Context) error {
				var data []byte
				var err error
				if c.NArg() == 0 {
					data, err = ioutil.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
					if debug {
						log.Printf("stdin data")
					}
				} else {
					data, err = ioutil.ReadFile(c.Args().Get(0))
					if err != nil {
						return err
					}
					if debug {
						log.Printf("file data")
					}
					if name {
						page = c.Args().Get(0)
					}
				}
				text := ""
				if binary {
					text, err = BytesToString(data)
					if err != nil {
						return err
					}
				} else {
					text = string(data)
				}
				exists, err := pageExists(server, page)
				if err != nil {
					return err
				}
				if exists {
					reader := bufio.NewReader(os.Stdin)
					fmt.Printf("Page '%s' exists, do you want to overwrite (y/n): ", page)
					answer, _ := reader.ReadString('\n')
					if !strings.Contains(strings.ToLower(answer), "y") {
						return nil
					}
				}

				if page == "" {
					// generate page name
					page = GetRandomName()
				}
				if encryptFlag || passphrase != "" {
					if debug {
						log.Println("Encryption activated")
					}
					if passphrase == "" {
						reader := bufio.NewReader(os.Stdin)
						fmt.Print("Enter passphrase: ")
						passphrase, _ = reader.ReadString('\n')
						passphrase = strings.TrimSpace(passphrase)
					}
					text, err = encrypt.EncryptString(text, passphrase)
					if err != nil {
						return err
					}
					encryptFlag = true
				}

				return uploadData(server, page, text, encryptFlag, store)
			},
		},
		{
			Name:    "download",
			Aliases: []string{"d"},
			Usage:   "download document",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "passphrase, a",
					Usage:       "passphrase to use for encryption",
					Destination: &passphrase,
				},
			},
			Action: func(c *cli.Context) error {
				page := ""
				if c.NArg() == 1 {
					page = c.Args().Get(0)
				} else {
					return errors.New("Must specify page")
				}
				return downloadData(server, page, passphrase)
			},
		},
	}

	errMain := app.Run(os.Args)
	if errMain != nil {
		log.Println(errMain)
	}
	return errMain
}
