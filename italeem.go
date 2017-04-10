package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/0xAX/notificator"
	"github.com/bgentry/speakeasy"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/puerkitobio/goquery"
	spin "github.com/tj/go-spin"
)

var notify *notificator.Notificator
var credFileName = ".italeem"

func home() string {
	home, err := homedir.Dir()
	check(err)
	return home + "/"
}

type Announcement struct {
	url      string
	text     string
	author   string
	time_ago string
}

func scrapeAnnouncements(resp *http.Response) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	countAnnouncements := 0

	// Find announcement items
	doc.Find(".messagemenu .dropdown-menu li a").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if !ok {
			log.Fatal("Cannot find href")
		}

		text := s.Find("span.notification-text").Text()

		time := s.Find("span.msg-time").Text()
		split_text := strings.Split(text, ": Announcements:")
		author, announcementText := split_text[0], split_text[1]

		// clean-up author string
		split_author := strings.Split(author, " posted in ")
		// course_code := split_author[1]
		author_slice := strings.Split(split_author[0], " ")
		author_name := strings.Join(author_slice[0:len(author_slice)-1], " ")

		announcement := Announcement{
			url:      url,
			text:     announcementText,
			author:   author_name,
			time_ago: time,
		}

		countAnnouncements++
		fmt.Printf("%s. %s - %s\n", color.CyanString(strconv.Itoa(countAnnouncements)), color.GreenString(announcement.author), color.YellowString(announcement.time_ago))
		fmt.Printf("→%s\n→ %s\n\n", announcement.text, announcement.url)
	})

	icon := "/home/user/icon.png"
	notify = notificator.New(notificator.Options{
		DefaultIcon: icon,
		AppName:     "Italeem CLI",
	})

	msg := strconv.Itoa(countAnnouncements) + " announcements fetched."
	notify.Push(msg, "", icon, notificator.UR_CRITICAL)
}

func createClient() http.Client {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	return http.Client{Jar: jar}
}

func login(client http.Client, loginUrl string, username string, password string) *http.Response {
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

func getAuthInfoFromUser() (string, string) {
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
func getCredFile() string {
	return home() + credFileName
}
func getCred() (string, string) {
	credFile := getCredFile()
	if _, err := os.Stat(credFile); os.IsNotExist(err) {
		username, password := getAuthInfoFromUser()

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	loginUrl := "http://italeem.iium.edu.my/2016/login/index.php"
	username, password := getCred()

	client := createClient()
	resp := login(client, loginUrl, username, password)
	fmt.Println("\n")
	scrapeAnnouncements(resp)
}
