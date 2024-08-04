package utils

import "fmt"

func ErrorLog(message string) {
	fmt.Printf("\033[31m- %s\033[0m\n", message)
}

func SuccessLog(message string) {
	fmt.Printf("\033[32m- %s\033[0m\n", message)
}

func InfoLog(message string) {
	fmt.Printf("\033[93m- %s\033[0m\n", message)
}
