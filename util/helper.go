package util

import (
	"bufio"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
)

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
	fmt.Println(msg)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	exit(errorCode)
}
