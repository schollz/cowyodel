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

func pageExists(serverString string, page string) (exists bool, err error) {
	server, user, password, err2 := parseURL(serverString)
	if err2 != nil {
		err = err2
		log.Trace("Problem parsing '%s'", serverString)
		return
	}
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
	if user != "" {
		req.SetBasicAuth(user, password)
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
		log.Trace("Problem parsing response")
		err = errors.New(target.Message)
	}
	exists = target.Exists
	return
}

func uploadData(serverString string, name string, codename string, text string, encrypt bool, store bool) (err error) {
	server, user, password, err2 := parseURL(serverString)
	if err2 != nil {
		err = err2
		log.Trace("Problem parsing '%s'", serverString)
		return
	}
	log.Trace("server: '%s'", server)
	log.Trace("uploading page: '%s' to '%s'", name, codename)

	type Payload struct {
		Page        string `json:"page"`
		Meta        string `json:"meta"`
		NewText     string `json:"new_text"`
		IsEncrypted bool   `json:"is_encrypted"`
		IsPrimed    bool   `json:"is_primed"`
	}

	data := Payload{
		Page:        codename,
		Meta:        name,
		NewText:     text,
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
	if user != "" {
		req.SetBasicAuth(user, password)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Trace("Problem with request")
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
		fmt.Printf("uploaded %s to %s\n", name, codename)
	} else {
		fmt.Println(target.Message)
	}
	return
}

func downloadData(serverString string, codename string, passphrase string) (err error) {
	server, user, password, err2 := parseURL(serverString)
	if err2 != nil {
		err = err2
		return
	}
	log.Trace("server: '%s'", server)
	log.Trace("page: '%s'", codename)

	type Payload struct {
		Page string `json:"page"`
	}

	data := Payload{
		Page: codename,
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
	if user != "" {
		req.SetBasicAuth(user, password)
		log.Trace("user: '%s'", user)
		log.Trace("password: '%s'", password)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Trace("Problem with request")
		return
	}
	defer resp.Body.Close()

	type Response struct {
		Name      string `json:"name"`
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
		fmt.Printf("'%s' not found", codename)
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
			uploadData(server, target.Name, codename, target.Text, true, false)
			return nil
		}
	}

	// assume its binary data
	binaryData, binaryAttemptError := StringToByte(target.Text)
	if binaryAttemptError == nil {
		err = ioutil.WriteFile(target.Name, binaryData, 0644)
		if err != nil {
			return
		}
		fmt.Printf("Wrote binary data to '%s'\n", target.Name)
		return
	}

	// its just a text file
	err = ioutil.WriteFile(target.Name, []byte(target.Text), 0644)
	if err != nil {
		return
	}
	fmt.Printf("Wrote text of '%s' to '%s'\n", codename, target.Name)
	return
}
