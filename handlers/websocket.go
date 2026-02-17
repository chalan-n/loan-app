// handlers/websocket.go
package handlers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)

func init() {
	go handleMessages()
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("WebSocket error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func WsHandler(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func WsConnect(c *websocket.Conn) {
	defer c.Close()

	clients[c] = true
	log.Printf("เชื่อมต่อแล้ว รวม %d คน", len(clients))

	for {
		mt, message, err := c.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			delete(clients, c)
			log.Printf("หลุดการเชื่อมต่อ เหลือ %d คน", len(clients))
			break
		}

		// เพิ่มบรรทัดนี้: รับข้อความแล้วกระจายต่อทันที!
		var data map[string]string
		if json.Unmarshal(message, &data) == nil {
			log.Printf("ได้รับแจ้งเตือนจากลูกค้า: %s - %s", data["title"], data["message"])
		}
		broadcast <- message // ← สำคัญมาก! ส่งต่อให้ทุกคน
	}
}
