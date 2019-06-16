package util

import (
	"fmt"
	"os"
)

func HandleError(err error, message string) {
	if err != nil {
		fmt.Println(message)
		os.Exit(2)
	}
}
