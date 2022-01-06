package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	ServerInit()
}

type FileData struct {
	Name      string
	Channel   string
	SizeField int
	Data      []byte
}
type Channel struct {
	Name    string
	Clients []net.Conn
}

type ChannelList struct {
	Channels []Channel
}

func ServerInit() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a command!")
		fmt.Println("Usage: start for init the server")
		return
	}
	COMMAND := arguments[1]

	if COMMAND == "start" {
		var channelList ChannelList
		server, err := net.Listen("tcp4", ":8000")
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
func handleConnection(client net.Conn, channelList *ChannelList) {
	var wg sync.WaitGroup
	var lock sync.RWMutex
	fmt.Printf("Serving %s\n", client.RemoteAddr().String())
	clientNumber := client.RemoteAddr().String()
	err := gob.NewEncoder(client).Encode("The client " + clientNumber[10:len(clientNumber)-1] + " is connected!")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		var operation string
		err = gob.NewDecoder(client).Decode(&operation)
		if err != nil {
			continue
		}
		fmt.Println(operation + " Paso")
		if operation == "send" {
			wg.Add(1)
			go send(client, channelList, &wg, &lock)
		} else if operation == "suscribe" {
			wg.Add(1)
			go suscribe(client, channelList, &wg, &lock)
		} else if operation == "create" {
			fmt.Println("Creating channel")
			wg.Add(1)
			go create(client, channelList, &wg, &lock)

		}

	}
}
func send(client net.Conn, channelList *ChannelList, wg *sync.WaitGroup, lock *sync.RWMutex) {
	var fileData FileData
	err := gob.NewDecoder(client).Decode(&fileData)
	if err != nil {
		return
	}
	fmt.Println(fileData)
	for value := range channelList.Channels {
		if channelList.Channels[value].Name == fileData.Channel {
			for _, clientData := range channelList.Channels[value].Clients {
				fmt.Println("stay in the channel")
				gob.NewEncoder(clientData).Encode("FileData")
				err = gob.NewEncoder(clientData).Encode(fileData)
				if err != nil {
					fmt.Println(err)
					return
				}
				gob.NewEncoder(client).Encode("Response")
				gob.NewEncoder(client).Encode("Data sent")
			}
		}
		fmt.Println(value)
	}
	gob.NewEncoder(client).Encode("Response")
	gob.NewEncoder(client).Encode("The channel does not exist!")
}
func create(client net.Conn, channelList *ChannelList, wg *sync.WaitGroup, lock *sync.RWMutex) {
	var channel Channel
	err := gob.NewDecoder(client).Decode(&channel)
	if err != nil {
		return
	}
	gob.NewEncoder(client).Encode("response")
	for value := range channelList.Channels {
		if channelList.Channels[value].Name == channel.Name {
			gob.NewEncoder(client).Encode("Channel already exists")
			fmt.Println("Channel already exists")
			return
		}
	}
	channelList.Channels = append(channelList.Channels, channel)
	gob.NewEncoder(client).Encode("Channel created")
	//fmt.Println("add channel" + channel.Name)
	//continue
}
func suscribe(client net.Conn, channelList *ChannelList, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	lock.Lock()
	var channelName string
	err := gob.NewDecoder(client).Decode(&channelName)
	if err != nil {
		fmt.Println(err)
		return
	}
	gob.NewEncoder(client).Encode("response")
	for value := range channelList.Channels {
		if channelList.Channels[value].Name == channelName {
			if contains(channelList.Channels[value].Clients, client) {
				gob.NewEncoder(client).Encode("Client already exist")
			} else {
				channelList.Channels[value].Clients = append(channelList.Channels[value].Clients, client)
				gob.NewEncoder(client).Encode("Client added to channel")
				fmt.Println("Client added to channel")
				return
			}

		}
	}
	gob.NewEncoder(client).Encode("The channel is not created")
	lock.Unlock()
}
func contains(clients []net.Conn, client net.Conn) bool {
	for _, value := range clients {
		if value == client {
			return true
		}
	}
	return false
}
