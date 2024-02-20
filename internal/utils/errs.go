package utils

import "fmt"

func Must(action string, err error) {
	if err != nil {
		fmt.Println("couldn't ", action, ": ", err.Error())
		panic(err)
	}
}
