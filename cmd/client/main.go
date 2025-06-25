package main

import (
	"chat/client"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Укажите имя пользователя как аргумент.")
		fmt.Println("Пример: go run main.go Bob")
		return
	}

	username := os.Args[1]
	client.RunClient(username)
}