package util

import (
	"fmt"
	"os"
)

func HandleError(err error, message string) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(" >", err.Error())
		os.Exit(2)
	}
}
