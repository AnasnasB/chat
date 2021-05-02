package main

import (
	"net/http"
	"io"
	"github.com/gorilla/websocket"
	"fmt"
	//"os"
	//"bufio"
	"unicode/utf8"
)

var upgrader = websocket.Upgrader{}
type ConnectionUsers [] struct {
	Username string
	Conn *websocket.Conn
}
var Userlist ConnectionUsers

func Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	
	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {return }
		
		fmt.Printf("%s\n", data)
		io.WriteString(w, "successful post")
	} else if req.Method == "OPTIONS" {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(405)
	}
	
}

func TakeMessage (w http.ResponseWriter, req *http.Request, conn *websocket.Conn, name string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		var k string
		for len(message) > 0 {
			r, size := utf8.DecodeRune(message)
			k+=string(r)
			message = message[size:]
		}
		var messagetoothers string
		messagetoothers = name+" : "+k
		for i:= range Userlist {
			if Userlist[i].Conn != conn {
				if err := Userlist[i].Conn.WriteMessage(websocket.TextMessage, []byte(messagetoothers)); err != nil {
					fmt.Println("error here",err)
					return
				}
				
			}
		}
		
		
		
		
	}
}

//func SendMessage (w http.ResponseWriter, req *http.Request, conn *websocket.Conn) {
//	for {
//		message:= "Proverka 123"
//		if message != "" {
//				for i:= range Userlist {
//					if Userlist[i].Conn != conn {
//						if err := Userlist[i].Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
//							fmt.Println("error here",err)
//							return
//						}
//				
//					}
//				}
//		}
//	}
//}

func Socket(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return
	}
	var name string
	for len(message) > 0 {
		r, size := utf8.DecodeRune(message)
		name+=string(r)
		message = message[size:]
	}
	
	
	var user struct {
		Username string
		Conn *websocket.Conn
	}
	user.Username = name 
	user.Conn = conn
	Userlist = append(Userlist,user)
	go TakeMessage (w, req, conn, name)	
	for {
		
	}
	
}

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/socket", Socket)
	
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

