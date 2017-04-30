package util

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
)

func CheckInternet() (bool, string) {
	url := "http://google.com/favicon.ico"

	resp, err := http.Head(url)
	if err != nil {
		return false, err.Error()
	}

	return true, strconv.Itoa(resp.StatusCode)
}

func checkError(err error) {
	if err != nil {
		exitWithMessage(err.Error(), 1)
	}
}

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	checkError(err)
}

func home() string {
	home, err := homedir.Dir()
	checkError(err)
	return home + "/"
}

func exit(errorCode int) {
	os.Exit(errorCode)
}

func exitWithMessage(msg string, errorCode int) {
	if errorCode > 0 {
		oops := fmt.Sprintf("%s\n", color.RedString("Oops! Something went wrong :("))

		if runtime.GOOS == "windows" {
			fmt.Fprintf(color.Output, oops)
		} else {
			fmt.Printf(oops)
		}
	}

	fmt.Println(msg)
	fmt.Println("\nPress enter to exit..")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	exit(errorCode)
}

func ExitIfNoInternet() {
	hasInternet, _ := CheckInternet()
	if !hasInternet {
		exitWithMessage("Seems like you've no internet connection.", 4)
	}
}
