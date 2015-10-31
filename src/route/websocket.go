package route

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"mongo"
)

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}

func EchoHandler(ws *websocket.Conn) {
	ActiveClients := make(map[ClientConn]int)

	defer func() {
		if err := ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()
	client := ws.Request().RemoteAddr
	log.Println("Client connected:", client)
	sockCli := ClientConn{ws, client}
	fmt.Println("sockCli")
	fmt.Println(sockCli)
	ActiveClients[sockCli] = 0
	log.Println("Number of clients connected ...", len(ActiveClients))

	var thread mongo.Thread
	for {
		if err := websocket.JSON.Receive(ws, &thread); err != nil {
			log.Println("Websocket Disconnected waiting", err.Error())
			delete(ActiveClients, sockCli)
			log.Println("Number of clients still connected ...", len(ActiveClients))
			return
		}

		for cs, _ := range ActiveClients {
			//		if err = Message.Send(cs.websocket, clientMessage); err != nil {
			if err := websocket.JSON.Send(cs.websocket, thread); err != nil {
				// we could not send the message to a peer
				log.Println("Could not send message to ", cs.clientIP, err.Error())
			}
		}

		log.Printf("thread=%#v\n", thread)
	}
	// send JSON type T
}
