package SSEHub

import (
	"ChoHanJi/useCases/AdminWaitingRoomUseCase"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type Message struct {
	MessageType string `json:"MessageType"`
	Message     string `json:"Message"`
}

type SSEHub struct {
	mu      sync.Mutex
	clients map[string][]chan []byte
}

var _ AdminWaitingRoomUseCase.IHub = (*SSEHub)(nil)

func New() *SSEHub {
	return &SSEHub{
		clients: make(map[string][]chan []byte),
	}
}

func (h *SSEHub) Subscribe(roomId string) (<-chan []byte, int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan []byte, 16)
	list := h.clients[roomId]
	index := len(list)
	h.clients[roomId] = append(list, ch)

	return ch, index
}

func (h *SSEHub) Unsubscribe(roomId string, index int) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, found := h.clients[roomId]
	if !found {
		return errors.New("roomId is not registered")
	}

	if index < 0 || index >= len(clients) {
		return errors.New("index out of range")
	}

	if clients[index] == nil {
		return nil // already unsubscribed
	}

	clients[index] = nil

	isAllChannelDeleted := true
	for _, channel := range clients {
		if channel != nil {
			isAllChannelDeleted = false
		}
	}

	if isAllChannelDeleted {
		delete(h.clients, roomId)
	}

	return nil
}

func (h *SSEHub) Publish(roomId, messageType, messageBody string) error {
	message := Message{messageType, messageBody}

	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("SSEHub.Publish: Could not marshal the message: %w", err)
	}

	h.mu.Lock()
	list, ok := h.clients[roomId]
	if !ok {
		h.mu.Unlock()
		return errors.New("roomId is not registered")
	}

	snapshot := make([]chan []byte, len(list))
	copy(snapshot, list)
	h.mu.Unlock()

	for _, ch := range snapshot {
		if ch == nil {
			continue
		}
		select {
		case ch <- msg:
		default:
		}
	}
	return nil
}
