package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	ServerInit()
}

type Message struct {
	Name       string
	Channel    string
	SizeField  int
	TypeOfData string
	Data       []byte
}
type Channel struct {
	Name    string
	Clients []net.Conn
}

var (
	channelList []Channel
)

func ServerInit() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a command!")
		fmt.Println("Usage: start for init the server")
		return
	}
	COMMAND := arguments[1]

	if COMMAND == "start" {

		server, err := net.Listen("tcp", ":8000")
		if err != nil {
			fmt.Println(err, "err")
			return
		}
		defer server.Close()
		for {
			client, err := server.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go handleConnection(client)
		}
	} else {
		fmt.Println("Please provide a correct command!")
		return
	}
}
func handleConnection(conn net.Conn) {
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	clientNumber := conn.RemoteAddr().String()
	message := writeMessage("The client " + clientNumber[10:] + " is connected!")
	err := enc.Encode(&message)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		opMessage := readMessage(conn, dec)
		fmt.Println(opMessage.Name)
		if opMessage.Name == "send" {
			fmt.Println("the data is in send")
			send(conn, dec, enc)
		} else if opMessage.Name == "suscribe" {
			suscribe(conn, dec, enc)
		} else if opMessage.Name == "create" {
			fmt.Println("the data is in create")
			create(conn, dec, enc)

		} else if opMessage.Name == "errorC" {
			fmt.Println(string(opMessage.Data))
			break
		} else if opMessage.Name == "receive" {
			message := writeMessage("File Send!")
			enc.Encode(&message)
		} else {
			message := writeMessage("Server: Unknown command " + opMessage.Name)
			enc.Encode(&message)
		}

	}
}
func create(client net.Conn, dec *gob.Decoder, enc *gob.Encoder) {
	var Channel Channel
	err := dec.Decode(&Channel)
	if err != nil {
		return
	}
	fmt.Println(Channel)
	for value := range channelList {
		if channelList[value].Name == Channel.Name {
			data := []byte("error, channel already exists")
			message := Message{Name: "Server", Channel: Channel.Name, SizeField: len(data), TypeOfData: "string", Data: data}
			enc.Encode(&message)
			fmt.Println("Channel already exists")
			return
		}
	}
	channelList = append(channelList, Channel)
	data := []byte("Channel created")
	fmt.Println(channelList)
	message := Message{Name: "Server", Channel: Channel.Name, SizeField: len(data), TypeOfData: "string", Data: data}
	enc.Encode(message)
}
func suscribe(client net.Conn, dec *gob.Decoder, enc *gob.Encoder) {
	var channelName string
	err := dec.Decode(&channelName)
	if err != nil {
		fmt.Println(err)
		return
	}
	for value := range channelList {
		if channelList[value].Name == channelName {
			if contains(channelList[value].Clients, client) {
				data := "error, you are already suscribed"
				message := writeMessage(data)
				enc.Encode(&message)
				return
			} else {
				channelList[value].Clients = append(channelList[value].Clients, client)
				data := "Client added to channel"
				message := writeMessage(data)
				enc.Encode(&message)
				return
			}

		}
	}
	data := "The channel is not created"
	message := writeMessage(data)
	enc.Encode(message)

}
func send(client net.Conn, dec *gob.Decoder, enc *gob.Encoder) {
	fileData := readMessage(client, dec)
	//fmt.Println(fileData, "fileData", fileData.TypeOfData)
	if fileData.Name == "Error" {
		messageForCLient := writeMessage("The file is not found")
		enc.Encode(&messageForCLient)
		return
	}
	for value := range channelList {
		if channelList[value].Name == fileData.Channel {
			fmt.Println("the channel is found")
			if len(channelList[value].Clients) == 0 {
				data := "error, no clients in channel"
				message := writeMessage(data)
				enc.Encode(&message)
				return
			} else {
				messageForCLient := writeMessage("The file is ready to be downloaded")
				enc.Encode(&messageForCLient)
				for _, clientData := range channelList[value].Clients {
					if clientData != client {
						err := gob.NewEncoder(clientData).Encode(&fileData)
						if err != nil {
							fmt.Println(err)
							return
						}
						fmt.Println("the data is sent")
					} else {
						continue
					}
				}
				return
			}
		}
	}
	message := writeMessage("The channel does not exist!")
	enc.Encode(message)
}
func contains(clients []net.Conn, client net.Conn) bool {
	for _, value := range clients {
		if value == client {
			return true
		}
	}
	return false
}
func readMessage(conn net.Conn, dec *gob.Decoder) Message {
	var message Message
	err := dec.Decode(&message)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Client disconnected")
			for i := range channelList {
				for j := range channelList[i].Clients {
					fmt.Println("paso2", j)
					if channelList[i].Clients[j] == conn {
						channelList[i].Clients = append(channelList[i].Clients[:j], channelList[i].Clients[j+1:]...)
						break
					}
				}
			}
			msg := "The client " + conn.RemoteAddr().String() + " get out"
			response := Message{Name: "errorC", Channel: "nil", SizeField: len([]byte(msg)), TypeOfData: "string", Data: []byte(msg)}
			return response
		}
		log.Fatal(err)
	}
	return message
}
func writeMessage(msg string) Message {
	response := Message{Name: "message", Channel: "nil", SizeField: len([]byte(msg)), TypeOfData: "string", Data: []byte(msg)}
	return response
}
