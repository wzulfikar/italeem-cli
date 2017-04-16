package util

import (
	"io/ioutil"
	"os"
	"strings"
)

func getCredFile() string {
	return home() + credFileName
}

func GetCred() (string, string) {
	key_string := home()
	if length := len(key_string); length > 16 {
		key_string = key_string[length-16:]
	} else if length < 16 {
		// pad key_string with `#` to make up 16 chars string
		key_string = strings.Repeat("#", 16-length) + key_string
	}

	key := []byte(key_string)
	credFile := getCredFile()
	if _, err := os.Stat(credFile); os.IsNotExist(err) {
		username, password := GetAuthInfoFromUser()

		cred := Enc(key, username+password)
		err := ioutil.WriteFile(credFile, []byte(cred), 0644)
		checkError(err)

		// trims space to avoid trailing "\n" char
		return strings.TrimSpace(username), strings.TrimSpace(password)
	}

	cred_string, err := ioutil.ReadFile(credFile)
	checkError(err)

	cred := Dec(key, string(cred_string))
	strCred := strings.Split(cred, "\n")
	if len(strCred) != 2 {
		deleteFile(credFile)
		exitWithMessage("Invalid cred file. Please try again :)", 2)
		return "", ""
	}
	username, password := strCred[0], strCred[1]
	return username, password
}
