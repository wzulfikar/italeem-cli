package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bgentry/speakeasy"
	"github.com/fatih/color"
	spin "github.com/tj/go-spin"
)

func Login(client http.Client, loginUrl string, username string, password string) *http.Response {
	msg := "Authenticating user"

	s := spin.New()
	for i := 0; i < 30; i++ {
		fmt.Printf("\r\033[36m%s\033[m: %s %s", msg, username, s.Next())
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := client.PostForm(loginUrl, url.Values{
		"username":         {username},
		"password":         {password},
		"rememberusername": {"1"},
	})
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// read the response body to a variable
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	exitIfNotAuthenticated(bodyString)

	//reset the response body to the original unread state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return resp
}

func exitIfNotAuthenticated(html string) {
	notLoggedInString := "Forgotten your username or password?"
	if strings.Contains(html, notLoggedInString) {
		fmt.Println(color.RedString("\nAuthentication failed. Please try again :)"))

		err := os.Remove(getCredFile())
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(2)
	}
}

func GetAuthInfoFromUser() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.CyanString("Enter username: "))
	username, _ := reader.ReadString('\n')

	password, err := speakeasy.Ask(color.CyanString("Password ") + "(it won't be displayed): ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return username, password
}
