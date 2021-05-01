
package main

import (
	"log"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"bufio"
	"unicode/utf8"
)

func TakeMessage (c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		var k string
		for len(message) > 0 {
			r, size := utf8.DecodeRune(message)
			k+=string(r)
			message = message[size:]
		}
		
		fmt.Println("-", k)
		
	}
	
}

func SendMessage (c *websocket.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		if message != "" {
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					fmt.Println(err)
					return
				}
		}
		message = ""
	}
}

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/socket", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go TakeMessage(c)
	go SendMessage(c)
	for {
		
	}
	
}
