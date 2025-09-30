package chekwebsocket

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ChekWebSocket() (string, error) {
	var websocket string
	if len(os.Args) > 1 {
		websocket = strings.Join(os.Args[1:], "")
	} else {
		// Запрашиваем  у пользователя название продукта, если аргумент не  указан
		fmt.Print("Enter the city name: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("Error reading input: %v", err)
		}
		websocket = strings.TrimSpace(input)
	}

	return websocket, nil
}