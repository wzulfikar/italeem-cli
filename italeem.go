package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/0xAX/notificator"
	"github.com/fatih/color"
	"github.com/puerkitobio/goquery"
	"github.com/urfave/cli"
)

var notify *notificator.Notificator

func scrapeAnnouncements(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	// text: notification-text
	// time: msg-time
	// link: li a

	// Find announcement items
	doc.Find(".messagemenu .dropdown-menu li a").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if !ok {
			log.Fatal("Cannot find href")
		}

		text := s.Find("span.notification-text").Text()

		time := s.Find("span.msg-time").Text()
		split_text := strings.Split(text, ": Announcements:")
		author, announcement := split_text[0], split_text[1]

		// clean-up author string
		split_author := strings.Split(author, " posted in ")
		// course_code := split_author[1]
		author_slice := strings.Split(split_author[0], " ")
		author_name := strings.Join(author_slice[0:len(author_slice)-1], " ")

		fmt.Printf("%d. %s - %s\n", i+1, color.GreenString(author_name), color.YellowString(time))
		fmt.Printf("→%s\n→ %s\n\n", announcement, url)
	})

	icon := "/home/user/icon.png"
	notify = notificator.New(notificator.Options{
		DefaultIcon: icon,
		AppName:     "Italeem CLI",
	})

	notify.Push("Finished fetching recent announcements", "", icon, notificator.UR_CRITICAL)
}

func main() {
	app := cli.NewApp()

	urlAnnouncements := "http://php.dev:8888/announcements.html"

	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:  "lang, l",
	// 		Value: "english",
	// 		Usage: "Language for the greeting",
	// 	},
	// 	cli.StringFlag{
	// 		Name:  "config, c",
	// 		Usage: "Load configuration from `FILE`",
	// 	},
	// }

	app.Commands = []cli.Command{
		{
			Name:    "messages",
			Aliases: []string{"m"},
			Usage:   "display recent messages from italeem",
			Action: func(c *cli.Context) error {
				fmt.Println("displaying messages")
				return nil
			},
		},
		{
			Name:    "announcements",
			Aliases: []string{"a"},
			Usage:   "display recent announcements from italeem",
			Action: func(c *cli.Context) error {
				fmt.Println("Displaying announcements:\n")
				scrapeAnnouncements(urlAnnouncements)
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}
