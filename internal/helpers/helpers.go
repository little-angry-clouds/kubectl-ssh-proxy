package helpers

import (
	"fmt"
	"os"
)

// CheckGenericError checks if there's an error, shows it and exits the program if it is
func CheckGenericError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// CheckActiveProcess checks if a process is active and exits the program if it is
func CheckActiveProcess(pidPath string) {
	if _, err := os.Stat(pidPath); err == nil {
		fmt.Println("# There's already an active process!")
		os.Exit(1)
	}
}
