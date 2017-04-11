package main

import (
	"fmt"

	"github.com/wzulfikar/italeem-cli/util"
)

func main() {
	loginUrl := "http://italeem.iium.edu.my/2016/login/index.php"
	username, password := util.GetCred()

	client := util.CreateClient()
	resp := util.Login(client, loginUrl, username, password)
	fmt.Println("\n")
	util.ScrapeAnnouncements(resp)
}
