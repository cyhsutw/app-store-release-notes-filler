package lib

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Channel struct {
	clients     map[*websocket.Conn]bool
	Subscribe   chan *websocket.Conn
	Unsubscribe chan *websocket.Conn
	Broadcast   chan []byte
	destroy     chan bool
}

func newChannel() *Channel {
	return &Channel{
		clients:   make(map[*websocket.Conn]bool),
		Subscribe: make(chan *websocket.Conn),
		Broadcast: make(chan []byte),
		destroy:   make(chan bool),
	}
}

func (ch *Channel) run() {
LOOP:
	for {
		select {
		case client := <-ch.Subscribe:
			ch.clients[client] = true
		case message := <-ch.Broadcast:
			for client := range ch.clients {
				client.WriteMessage(websocket.TextMessage, message)
			}
		case flag := <-ch.destroy:
			if flag {
				close(ch.destroy)
				close(ch.Subscribe)
				close(ch.Broadcast)

				for client := range ch.clients {
					client.Close()
				}

				ch.clients = map[*websocket.Conn]bool{}

				break LOOP
			}
		}

	}
}

type channels struct {
	mutex    sync.RWMutex
	channels map[uint]*Channel
}

var allChannels = channels{channels: map[uint]*Channel{}}

func CreateChannel(id uint) (*Channel, error) {
	allChannels.mutex.Lock()
	defer allChannels.mutex.Unlock()

	if findChannelWithoutLock(id) != nil {
		message := fmt.Sprintf("channel for `%d` already exists", id)
		return nil, errors.New(message)
	}

	newChannel := newChannel()
	allChannels.channels[id] = newChannel

	go newChannel.run()

	return newChannel, nil
}

func FindChannel(id uint) *Channel {
	allChannels.mutex.RLock()
	defer allChannels.mutex.RUnlock()

	return findChannelWithoutLock(id)
}

func DestroyChannel(id uint) {
	allChannels.mutex.Lock()
	defer allChannels.mutex.Unlock()

	channel := findChannelWithoutLock(id)
	if channel == nil {
		return
	}

	channel.destroy <- true

	delete(allChannels.channels, id)
}

// private
func findChannelWithoutLock(id uint) *Channel {
	if channel, ok := allChannels.channels[id]; ok {
		return channel
	}

	return nil
}
