package main

import (
	"log"
	"net"

	"chat/chatpb"
	"chat/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}

	grpcServer := grpc.NewServer()
	service := server.NewChatService()
	controller := server.NewChatController(service)
	chatpb.RegisterChatServiceServer(grpcServer, controller)

	log.Println("gRPC-сервер запущен на :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("ошибка запуска gRPC-сервера: %v", err)
	}
}
