package helpers

import (
	"fmt"
	"os"
)

func CheckGenericError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CheckActiveProcess(pidPath string) {
	if _, err := os.Stat(pidPath); err == nil {
		fmt.Println("There's already an active process!")
		os.Exit(1)
	}
}
