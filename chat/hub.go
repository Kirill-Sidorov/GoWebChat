package chat

import "log"

var (
	clients    = make(map[*Client]bool)
	register   = make(chan *Client)
	unregister = make(chan *Client)
	broadcast  = make(chan []byte)
)

func Run() {
	for {
		select {
		case client := <-register:
			clients[client] = true
		case client := <-unregister:
			if _, ok := clients[client]; ok {
				log.Printf("Unregister client - %s from chat", client.name)
				delete(clients, client)
				close(client.send)
			}
		case message := <-broadcast:
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}
