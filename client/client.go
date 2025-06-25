package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"chat/chatpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func parseID(s string) int64 {
	var id int64
	fmt.Sscanf(s, "%d", &id)
	return id
}

func RunClient(username string) {
	conn, err := grpc.NewClient("dns:///localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := chatpb.NewChatServiceClient(conn)
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("[%s] Введите команду или сообщение (напишите /help для списка команд):\n", username)
	for {
		fmt.Printf("[%s] > ", username)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case line == "/quit":
			return

		case line == "/help":
			fmt.Println("Доступные команды:")
			fmt.Println("  /show                  — Показать все сообщения чата")
			fmt.Println("  /edit ID НовыйТекст    — Редактировать Ваше сообщение по ID")
			fmt.Println("  /del ID                — Удалить Ваше сообщение по ID")
			fmt.Println("  /quit                  — Выйти из чата")
			fmt.Println("  /help                  — Показать эту справку")
			fmt.Println("Пример использования:")
			fmt.Printf("  /edit 2 Привет, мир!\n")

		case line == "/show":
			clearScreen()
			resp, err := client.GetMessages(context.Background(), &chatpb.GetMessagesRequest{})
			if err != nil {
				fmt.Println("Ошибка получения сообщений:", err)
				continue
			}
			for _, m := range resp.Messages {
				ts := m.Timestamp.AsTime().Local().Format("2006-01-02 15:04:05")
				fmt.Printf("#%d [%s] %s: %s\n", m.Id, ts, m.User, m.Content)
			}

		case strings.HasPrefix(line, "/edit "):
			parts := strings.SplitN(line, " ", 3)
			if len(parts) != 3 {
				fmt.Println("Неверный формат. Используйте: /edit ID НовыйТекст")
				continue
			}
			id := parseID(parts[1])
			_, err := client.EditMessage(context.Background(), &chatpb.EditMessageRequest{
				Id:         id,
				User:       username,
				NewContent: parts[2],
			})
			if err != nil {
				fmt.Println("Ошибка редактирования:", err)
			}

		case strings.HasPrefix(line, "/del "):
			parts := strings.SplitN(line, " ", 2)
			if len(parts) != 2 {
				fmt.Println("Неверный формат. Используйте: /del ID")
				continue
			}
			id := parseID(parts[1])
			_, err := client.DeleteMessage(context.Background(), &chatpb.DeleteMessageRequest{
				Id:   id,
				User: username,
			})
			if err != nil {
				fmt.Println("Ошибка удаления:", err)
			}

		default:
			_, err := client.CreateMessage(context.Background(), &chatpb.CreateMessageRequest{
				User:    username,
				Content: line,
			})
			if err != nil {
				fmt.Println("Ошибка отправки:", err)
			}
		}
	}
}

