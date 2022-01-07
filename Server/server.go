package main

import (
	"encoding/gob"
	"fmt"
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

type ChannelList struct {
	Channels []Channel
}

var (
	channelList ChannelList
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
			fmt.Println(err)
			return
		}
		defer server.Close()
		for {
			fmt.Println(channelList, "channelList")
			client, err := server.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go handleConnection(client, &channelList)
		}
	} else {
		fmt.Println("Please provide a correct command!")
		return
	}
}
func handleConnection(conn net.Conn, channelList *ChannelList) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	clientNumber := conn.RemoteAddr().String()
	message := writeMessage("The client " + clientNumber[10:len(clientNumber)-1] + " is connected!")
	err := gob.NewEncoder(conn).Encode(&message)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		opMessage := readMessage(conn)
		fmt.Println(opMessage.Name)
		if opMessage.Name == "send" {
			fmt.Println("the data is in send")
			send(conn, channelList)
		} else if opMessage.Name == "suscribe" {
			suscribe(conn, channelList)
		} else if opMessage.Name == "create" {
			fmt.Println("the data is in create")
			create(conn, channelList)

		} else {
			fmt.Println("the data is in else")
		}

	}
}
func create(client net.Conn, channelList *ChannelList) {
	var channel Channel
	err := gob.NewDecoder(client).Decode(&channel)
	if err != nil {
		return
	}
	fmt.Println(channel)
	for value := range channelList.Channels {
		if channelList.Channels[value].Name == channel.Name {
			data := []byte("error, channel already exists")
			message := Message{Name: "Server", Channel: channel.Name, SizeField: len(data), TypeOfData: "string", Data: data}
			gob.NewEncoder(client).Encode(&message)
			fmt.Println("Channel already exists")
			return
		}
	}
	channelList.Channels = append(channelList.Channels, channel)
	data := []byte("Channel created")
	message := Message{Name: "Server", Channel: channel.Name, SizeField: len(data), TypeOfData: "string", Data: data}
	gob.NewEncoder(client).Encode(message)
}
func suscribe(client net.Conn, channelList *ChannelList) {
	var channelName string
	err := gob.NewDecoder(client).Decode(&channelName)
	if err != nil {
		fmt.Println(err)
		return
	}
	for value := range channelList.Channels {
		if channelList.Channels[value].Name == channelName {
			if contains(channelList.Channels[value].Clients, client) {
				data := "error, you are already suscribed"
				message := writeMessage(data)
				gob.NewEncoder(client).Encode(&message)
			} else {
				channelList.Channels[value].Clients = append(channelList.Channels[value].Clients, client)
				data := "Client added to channel"
				message := writeMessage(data)
				gob.NewEncoder(client).Encode(&message)
				return
			}

		}
	}
	data := "The channel is not created"
	message := writeMessage(data)
	gob.NewEncoder(client).Encode(message)

}
func send(client net.Conn, channelList *ChannelList) {
	fileData := readMessage(client)
	fmt.Println(fileData, "fileData", fileData.TypeOfData)
	for value := range channelList.Channels {
		if channelList.Channels[value].Name == fileData.Channel {
			fmt.Println("the channel is found")
			for _, clientData := range channelList.Channels[value].Clients {
				fmt.Println(clientData, "clientData")
				fmt.Println(client, "client")
				err := gob.NewEncoder(clientData).Encode(&fileData)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("the data is sent")
			}
			return
		}
	}
	message := writeMessage("The channel does not exist!")
	gob.NewEncoder(client).Encode(message)
}
func contains(clients []net.Conn, client net.Conn) bool {
	for _, value := range clients {
		if value == client {
			return true
		}
	}
	return false
}
func readMessage(conn net.Conn) Message {
	var message Message
	err := gob.NewDecoder(conn).Decode(&message)
	if err != nil {
		log.Fatal(err)
	}
	return message
}
func writeMessage(msg string) Message {
	fmt.Println(msg)
	response := Message{Name: "message", Channel: "nil", SizeField: len([]byte(msg)), TypeOfData: "string", Data: []byte(msg)}
	return response
}
