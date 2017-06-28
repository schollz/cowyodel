package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"
)

var debug bool

func main() {
	app := cli.NewApp()
	var passphrase, page, server string
	var encrypt, store, direct bool
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server",
			Value:       "https://cowyo.com",
			Usage:       "server to use",
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
					Destination: &encrypt,
				},
				cli.BoolFlag{
					Name:        "store, s",
					Usage:       "store and persist after reading",
					Destination: &store,
				},
				cli.StringFlag{
					Name:        "page, p",
					Usage:       "specific page to use",
					Destination: &page,
				},
				cli.StringFlag{
					Name:        "passphrase, a",
					Usage:       "passphrase to use for encryption",
					Destination: &passphrase,
				},
				cli.BoolFlag{
					Name:        "direct",
					Usage:       "direct mode (Gzip + Base64 encoding)",
					Destination: &direct,
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
				}
				dataString := ""
				if direct {
					dataString, err = BytesToString(data)
					if err != nil {
						return err
					}
				} else {
					dataString = string(data)
				}
				return uploadData(server, page, dataString, encrypt, passphrase, store)
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
				cli.BoolFlag{
					Name:        "direct",
					Usage:       "direct mode (Gzip + Base64 encoding)",
					Destination: &direct,
				},
			},
			Action: func(c *cli.Context) error {
				page := ""
				if c.NArg() == 1 {
					page = c.Args().Get(0)
				} else {
					return errors.New("Must specify page")
				}
				return downloadData(server, page, passphrase, direct)
			},
		},
	}

	errMain := app.Run(os.Args)
	if errMain != nil {
		log.Println(errMain)
	}

}

func uploadData(server string, page string, text string, encrypt bool, passphrase string, store bool) (err error) {
	if page == "" {
		// generate page name
		page = GetRandomName(1)
	}
	if encrypt || passphrase != "" {
		if debug {
			log.Println("Encryption activated")
		}
		if passphrase == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter passphrase: ")
			passphrase, _ = reader.ReadString('\n')
			passphrase = strings.TrimSpace(passphrase)
		}
		text, err = EncryptString(text, passphrase)
		if err != nil {
			return err
		}
		encrypt = true
	}

	type Payload struct {
		NewText     string `json:"new_text"`
		Page        string `json:"page"`
		IsEncrypted bool   `json:"is_encrypted"`
		IsPrimed    bool   `json:"is_primed"`
	}

	data := Payload{
		NewText:     text,
		Page:        page,
		IsEncrypted: encrypt,
		IsPrimed:    !store,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", server+"/update", body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	type Response struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}
	var target Response
	json.NewDecoder(resp.Body).Decode(&target)
	if debug {
		log.Printf("%v", target)
	}
	if target.Message == "Saved" {
		fmt.Printf("uploaded to %s\n", page)
	} else {
		fmt.Println(target.Message)
	}
	return
}

func downloadData(server string, page string, passphrase string, direct bool) (err error) {
	type Payload struct {
		Page string `json:"page"`
	}

	data := Payload{
		Page: page,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", server+"/relinquish", body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	type Response struct {
		Destroyed bool   `json:"destroyed"`
		Encrypted bool   `json:"encrypted"`
		Locked    bool   `json:"locked"`
		Message   string `json:"message"`
		Success   bool   `json:"success"`
		Text      string `json:"text"`
	}
	var target Response
	json.NewDecoder(resp.Body).Decode(&target)
	if debug {
		log.Printf("%v", target)
	}
	if target.Encrypted {
		if debug {
			log.Println("Decryption activated")
		}
		if passphrase == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Fprint(os.Stderr, "Enter passphrase: ")
			passphrase, _ = reader.ReadString('\n')
			passphrase = strings.TrimSpace(passphrase)
		}
		var decrypted string
		decrypted, err = DecryptString(target.Text, passphrase)
		if err == nil {
			target.Text = decrypted
		}
	}

	if direct {
		var data []byte
		data, err = StringToByte(target.Text)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(page, data, 0644)
		if err != nil {
			return
		}
		fmt.Printf("Wrote binary data to '%s'\n", page)
	} else {
		err = ioutil.WriteFile(page, []byte(target.Text), 0644)
		if err != nil {
			return
		}
		fmt.Printf("Wrote text to '%s'\n", page)
	}
	return
}
