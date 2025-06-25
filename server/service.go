package server

import (
	"errors"
	"sort"
	"sync"
	"time"

	"chat/chatpb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatService struct {
	mu       sync.Mutex
	messages map[int64]*chatpb.Message
	nextID   int64
}

func NewChatService() *ChatService {
	return &ChatService{
		messages: make(map[int64]*chatpb.Message),
		nextID:   1,
	}
}

func (s *ChatService) Create(user, content string) *chatpb.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := &chatpb.Message{
		Id:        s.nextID,
		User:      user,
		Content:   content,
		Timestamp: timestamppb.New(time.Now()),
	}
	s.messages[s.nextID] = msg
	s.nextID++
	return msg
}

func (s *ChatService) Edit(id int64, user, newContent string) (*chatpb.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.messages[id]
	if !ok {
		return nil, errors.New("сообщение не найдено")
	}
	if msg.User != user {
		return nil, errors.New("доступ запрещен: это не ваше сообщение")
	}

	msg.Content = newContent
	return msg, nil
}

func (s *ChatService) Delete(id int64, user string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.messages[id]
	if !ok {
		return errors.New("сообщение не найдено")
	}
	if msg.User != user {
		return errors.New("доступ запрещен: это не ваше сообщение")
	}

	delete(s.messages, id)
	return nil
}

func (s *ChatService) List() []*chatpb.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	msgs := make([]*chatpb.Message, 0, len(s.messages))
	for _, msg := range s.messages {
		msgs = append(msgs, msg)
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Timestamp.AsTime().Before(msgs[j].Timestamp.AsTime())
	})

	return msgs
}
