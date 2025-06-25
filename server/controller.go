package server

import (
	"context"

	"chat/chatpb"
)

type ChatController struct {
	chatpb.UnimplementedChatServiceServer
	svc *ChatService
}

func NewChatController(svc *ChatService) *ChatController {
	return &ChatController{svc: svc}
}

func (c *ChatController) CreateMessage(ctx context.Context, req *chatpb.CreateMessageRequest) (*chatpb.CreateMessageResponse, error) {
	msg := c.svc.Create(req.User, req.Content)
	return &chatpb.CreateMessageResponse{Message: msg}, nil
}

func (c *ChatController) EditMessage(ctx context.Context, req *chatpb.EditMessageRequest) (*chatpb.EditMessageResponse, error) {
	msg, err := c.svc.Edit(req.Id, req.User, req.NewContent)
	if err != nil {
		return nil, err
	}
	return &chatpb.EditMessageResponse{Message: msg}, nil
}

func (c *ChatController) DeleteMessage(ctx context.Context, req *chatpb.DeleteMessageRequest) (*chatpb.DeleteMessageResponse, error) {
	err := c.svc.Delete(req.Id, req.User)
	if err != nil {
		return nil, err
	}
	return &chatpb.DeleteMessageResponse{Success: true}, nil
}

func (c *ChatController) GetMessages(ctx context.Context, req *chatpb.GetMessagesRequest) (*chatpb.GetMessagesResponse, error) {
	msgs := c.svc.List()
	return &chatpb.GetMessagesResponse{Messages: msgs}, nil
}
