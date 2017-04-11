package util

import homedir "github.com/mitchellh/go-homedir"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func home() string {
	home, err := homedir.Dir()
	check(err)
	return home + "/"
}
