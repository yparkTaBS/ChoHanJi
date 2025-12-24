package SSEHub

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

type Message struct {
	MessageType string `json:"MessageType"`
	Message     string `json:"Message"`
}

type SSEHub struct {
	mu      sync.RWMutex
	clients map[string]map[string]chan []byte
}

func New() *SSEHub {
	return &SSEHub{
		clients: make(map[string]map[string]chan []byte),
	}
}

func (h *SSEHub) Subscribe(roomId, subscriberId string) <-chan []byte {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan []byte, 16)

	clients, found := h.clients[roomId]
	if !found {
		h.clients[roomId] = make(map[string]chan []byte)
		clients = h.clients[roomId]
	}

	client, found := clients[subscriberId]
	if found {
		close(client)
	}
	clients[subscriberId] = ch

	return ch
}

func (h *SSEHub) Unsubscribe(roomId, subscriberId string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, found := h.clients[roomId]
	if !found {
		return errors.New("roomId is not registered")
	}

	ch, found := clients[subscriberId]
	if !found {
		return nil
	}

	delete(clients, subscriberId)
	close(ch)

	if len(clients) == 0 {
		delete(h.clients, roomId)
	}

	return nil
}

func (h *SSEHub) Publish(roomId, subscriberId, messageType, messageBody string) error {
	message := Message{messageType, messageBody}

	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("SSEHub.Publish: Could not marshal the message: %w", err)
	}

	h.mu.RLock()
	clients, found := h.clients[roomId]
	if !found {
		h.mu.RUnlock()
		return errors.New("roomId is not registered")
	}

	logger := slog.Default()
	logger.Error("---")
	logger.Error(subscriberId)
	logger.Error("---")
	for subId := range clients {
		logger.Error(subId)
	}

	client, found := clients[subscriberId]
	if !found {
		h.mu.RUnlock()
		return fmt.Errorf("subscriber, %s, is not found", subscriberId)
	}

	select {
	case client <- msg:
	default:
	}

	h.mu.RUnlock()

	return nil
}

func (h *SSEHub) PublishToAll(roomId, messageType string, messageBody string) error {
	message := Message{messageType, messageBody}

	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("SSEHub.Publish: Could not marshal the message: %w", err)
	}

	h.mu.RLock()
	clients, ok := h.clients[roomId]
	if !ok {
		h.mu.RUnlock()
		return errors.New("roomId is not registered")
	}

	chans := make([]chan []byte, 0, len(clients))
	for _, ch := range clients {
		if ch != nil {
			chans = append(chans, ch)
		}
	}
	h.mu.RUnlock()

	for _, ch := range chans {
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
