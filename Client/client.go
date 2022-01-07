package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

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

func main() {
	ClientInit()
}

func ClientInit() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	var operation string
	var name string
	var channel string
	for {
		go receive(conn)
		_, err := fmt.Scan(&operation)
		if err != nil {
			fmt.Println(err, "read operation")
		}
		newMessage := Message{Name: operation, Channel: "", SizeField: 0, TypeOfData: "", Data: []byte{}}
		err = gob.NewEncoder(conn).Encode(&newMessage)
		if err != nil {
			fmt.Println(err, "inicio")
			continue
		}
		if operation == "send" {
			_, err := fmt.Scan(&name)
			if err != nil {
				fmt.Println(err, "read name")
				continue
			}
			_, err = fmt.Scan(&channel)
			if err != nil {
				fmt.Println(err, "read channel")
				continue
			}
			send(name, channel, conn)

		} else if operation == "create" {
			_, err := fmt.Scan(&name)
			if err != nil {
				_, err := fmt.Scan(&name)
				if err != nil {
					fmt.Println(err, "read name")
					continue
				}
			}
			channel := Channel{name, []net.Conn{}}
			err = gob.NewEncoder(conn).Encode(&channel)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else if operation == "suscribe" {
			_, err := fmt.Scan(&channel)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Suscribing to channel: " + channel)
			suscribe(channel, conn)
		} else if operation == "receive" {
			go receive(conn)
		} else {
			continue
		}

	}
}
func send(name string, channel string, conn net.Conn) {

	fmt.Println(channel)
	data, err := ioutil.ReadFile(name)
	fmt.Println(string(data))
	if err != nil {
		fmt.Println(err)
	}
	res1 := strings.Split(name, "/")
	typeofData := "FileData"
	fileData := Message{res1[len(res1)-1], channel, len(data), typeofData, data}
	err = gob.NewEncoder(conn).Encode(&fileData)

	if err != nil {
		fmt.Println(err)
		return
	}
}
func suscribe(channel string, conn net.Conn) {
	err := gob.NewEncoder(conn).Encode(&channel)
	if err != nil {
		fmt.Println(err)
	}
}

func receive(conn net.Conn) {
	var message Message
	err := gob.NewDecoder(conn).Decode(&message)
	if err != nil {
		log.Fatal(err)
	}
	if message.TypeOfData == "FileData" {
		fmt.Println("Data downloaded!")
		ioutil.WriteFile(message.Name, message.Data, 0644)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
	} else {
		fmt.Println(string(message.Data))
	}
}
