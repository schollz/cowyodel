package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/schollz/cowyo/encrypt"
)

func pageExists(server string, page string) (exists bool, err error) {
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

	req, err := http.NewRequest("POST", server+"/exists", body)
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
		Exists  bool   `json:"exists"`
	}
	var target Response
	json.NewDecoder(resp.Body).Decode(&target)
	log.Trace("%v", target)
	if !target.Success {
		err = errors.New(target.Message)
	}
	exists = target.Exists
	return
}

func uploadData(server string, page string, text string, encrypt bool, store bool) (err error) {
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
	log.Trace("%v", target)
	if target.Message == "Saved" {
		fmt.Printf("uploaded to %s\n", page)
	} else {
		fmt.Println(target.Message)
	}
	return
}

func downloadData(server string, page string, passphrase string) (err error) {
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
	log.Trace("%+v", target)
	if target.Text == "" {
		fmt.Printf("'%s' not found", page)
		return nil
	}
	if target.Encrypted {
		log.Trace("Decryption activated")
		if passphrase == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Fprint(os.Stderr, "Enter passphrase: ")
			passphrase, _ = reader.ReadString('\n')
			passphrase = strings.TrimSpace(passphrase)
		}
		var decrypted string
		decrypted, err = encrypt.DecryptString(target.Text, passphrase)
		if err == nil {
			target.Text = decrypted
		} else {
			fmt.Println("Incorrect password, returning to server")
			uploadData(server, page, target.Text, true, false)
			return nil
		}
	}

	// assume its binary data
	binaryData, binaryAttemptError := StringToByte(target.Text)
	if binaryAttemptError == nil {
		err = ioutil.WriteFile(page, binaryData, 0644)
		if err != nil {
			return
		}
		fmt.Printf("Wrote binary data to '%s'\n", page)
		return
	}

	// its just a text file
	err = ioutil.WriteFile(page, []byte(target.Text), 0644)
	if err != nil {
		return
	}
	fmt.Printf("Wrote text to '%s'\n", page)
	return
}
