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
