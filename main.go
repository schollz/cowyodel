package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/cowyo/encrypt"
	"github.com/schollz/lumber"
	"github.com/urfave/cli"
)

var debug bool
var log *lumber.ConsoleLogger
var version string

func main() {
	run()
}

func run() error {
	var passphrase, page, codename, server string
	var encryptFlag, store, name bool
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

		`
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
			},
			Action: func(c *cli.Context) error {
				var data []byte
				var err error
				if debug {
					log = lumber.NewConsoleLogger(lumber.TRACE)
				} else {
					log = lumber.NewConsoleLogger(lumber.WARN)
				}
				codename = GetRandomName()
				if c.NArg() == 0 {
					log.Trace("stdin")
					data, err = ioutil.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
				} else {
					log.Trace("file data")
					page = filepath.Base(c.Args().Get(0))
					data, err = ioutil.ReadFile(c.Args().Get(0))
					if err != nil {
						return err
					}
					if name {
						codename = page
					}
				}

				text := ""
				dataType := "textual"
				log.Trace("[%+v]\n", http.DetectContentType(data))
				if strings.Contains(http.DetectContentType(data), "text") {
					log.Trace("Assuming textual data")
					text = string(data)
				} else {
					log.Trace("Assuming binary data")
					dataType = "binary"
					text, err = BytesToString(data)
					if err != nil {
						return err
					}
				}
				exists, err := pageExists(server, codename)
				if err != nil {
					log.Trace("Could not check if exists")
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

				if encryptFlag || passphrase != "" {
					if debug {
						log = lumber.NewConsoleLogger(lumber.TRACE)
					} else {
						log = lumber.NewConsoleLogger(lumber.WARN)
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

				err = uploadData(server, page, codename, text, encryptFlag, store)
				if err == nil {
					fmt.Printf("Uploaded %s (%s data). Your codephrase: %s\n\n", page, dataType, codename)
					if dataType != "binary" {
						fmt.Printf("View/edit your data:\n\n\t%s/%s\n\n", server, codename)
					}
					if strings.Contains(server, "cowyo") {
						fmt.Printf("Download using cowyodel:\n\n\tcowyodel download %s\n\n", codename)
					} else {
						fmt.Printf("Download using cowyodel:\n\n\tcowyodel --server %s download %s\n\n", server, codename)
					}
				}
				return err
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
				if debug {
					log = lumber.NewConsoleLogger(lumber.TRACE)
				} else {
					log = lumber.NewConsoleLogger(lumber.WARN)
				}

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
		fmt.Print(errMain.Error())
	}
	return errMain
}
