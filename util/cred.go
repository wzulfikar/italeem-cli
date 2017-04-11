package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func getCredFile() string {
	return home() + credFileName
}

func GetCred() (string, string) {
	credFile := getCredFile()
	if _, err := os.Stat(credFile); os.IsNotExist(err) {
		username, password := GetAuthInfoFromUser()

		cred := []byte(username + password)
		err := ioutil.WriteFile(credFile, cred, 0644)
		check(err)

		// trims space to avoid trailing "\n" char
		return strings.TrimSpace(username), strings.TrimSpace(password)
	}

	cred, err := ioutil.ReadFile(credFile)
	check(err)

	strCred := strings.Split(string(cred), "\n")
	if len(strCred) != 2 {
		fmt.Println("Invalid cred file")
		return "", ""
	}
	username, password := strCred[0], strCred[1]
	return username, password
}
