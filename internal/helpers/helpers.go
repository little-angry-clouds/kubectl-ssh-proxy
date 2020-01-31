package helpers

import (
	"fmt"
	"os"
)

// CheckGenericError checks if there's an error, shows it and exits the program if it is
func CheckGenericError(err error) {
	if err != nil {
		message := fmt.Sprintf("An error was detected, exiting: %s", err)
		fmt.Println(message)
		os.Exit(1)
	}
}
