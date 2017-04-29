package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bgentry/speakeasy"
	"github.com/fatih/color"
	spin "github.com/tj/go-spin"
)

func Login(client http.Client, loginUrl string, username string, password string) *http.Response {
	msg := "Authenticating user:"

	if runtime.GOOS == "windows" {
		fmt.Fprintf(color.Output, "\r%s %s", color.CyanString(msg), username)
	} else {
		s := spin.New()
		for i := 0; i < 30; i++ {
			fmt.Printf("\r%s %s %s", color.CyanString(msg), username, s.Next())
			time.Sleep(100 * time.Millisecond)
		}
	}

	resp, err := client.PostForm(loginUrl, url.Values{
		"username":         {username},
		"password":         {password},
		"rememberusername": {"1"},
	})
	checkError(err)

	// read the response body to a variable
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	exitIfNotAuthenticated(bodyString)

	//reset the response body to the original unread state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	resp.Body.Close()

	fmt.Printf("\n")

	return resp
}

func exitIfNotAuthenticated(html string) {
	notLoggedInString := "Forgotten your username or password?"
	if strings.Contains(html, notLoggedInString) {
		auth_failed_msg := color.RedString("\nAuthentication failed. Please try again..")
		if runtime.GOOS == "windows" {
			fmt.Fprintf(color.Output, auth_failed_msg)
		} else {
			fmt.Println(auth_failed_msg)
		}

		err := os.Remove(getCredFile())
		checkError(err)
		os.Exit(2)
	}
}

func GetAuthInfoFromUser() (string, string) {
	ask_username := "Enter username/matric: "
	ask_password := "Password "

	reader := bufio.NewReader(os.Stdin)

	if runtime.GOOS == "windows" {
		fmt.Print(ask_username)
		username, _ := reader.ReadString('\n')

		password, err := speakeasy.Ask(ask_password + "(it won't be displayed): ")
		checkError(err)

		return username, password
	} else {
		fmt.Print(color.CyanString(ask_username))
		username, _ := reader.ReadString('\n')

		password, err := speakeasy.Ask(color.CyanString(ask_password) + "(it won't be displayed): ")
		checkError(err)

		return username, password
	}
}
