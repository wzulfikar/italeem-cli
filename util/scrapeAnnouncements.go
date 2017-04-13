package util

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/0xAX/notificator"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

var notify *notificator.Notificator
var credFileName = ".italeem"

type Announcement struct {
	url      string
	text     string
	author   string
	time_ago string
}

func ScrapeAnnouncements(resp *http.Response) {
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

		msg := fmt.Sprintf("%s. %s - %s\n", color.CyanString(strconv.Itoa(countAnnouncements)), color.GreenString(announcement.author), color.YellowString(announcement.time_ago))
		msg += fmt.Sprintf("→%s\n→ %s\n\n", announcement.text, announcement.url)

		if runtime.GOOS == "windows" {
			fmt.Fprintf(color.Output, msg)
		} else {
			fmt.Printf(msg)
		}
	})

	icon := "/home/user/icon.png"
	notify = notificator.New(notificator.Options{
		DefaultIcon: icon,
		AppName:     "Italeem CLI",
	})

	msg := strconv.Itoa(countAnnouncements) + " announcements fetched."
	notify.Push(msg, "", icon, notificator.UR_CRITICAL)

	// exit when user press enter
	fmt.Print("Finished fetching announcements. Press enter to exit..")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
