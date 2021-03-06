package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
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
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	defer conn.Close()
	var operation string
	reader := bufio.NewReader(os.Stdin)
	for {
		go receive(conn, dec)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		res1 := strings.Split(text, " ")
		operation = res1[0]
		if operation == "send" && len(res1) == 3 {
			var name string
			var channel string
			newMessage := Message{Name: operation, Channel: "", SizeField: 0, TypeOfData: "", Data: []byte{}}
			err = enc.Encode(&newMessage)
			if err != nil {
				fmt.Println(err)
				continue
			}
			name = res1[1]
			channel = res1[2]
			send(name, channel, conn, enc)

		} else if operation == "create" && len(res1) == 2 {
			var name string
			newMessage := Message{Name: operation, Channel: "", SizeField: 0, TypeOfData: "", Data: []byte{}}
			err = enc.Encode(&newMessage)
			if err != nil {
				fmt.Println(err, "inicio")
				continue
			}
			name = res1[1]

			channel := Channel{name, []net.Conn{}}
			err = enc.Encode(&channel)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else if operation == "suscribe" && len(res1) == 2 {
			var Channel string
			newMessage := Message{Name: operation, Channel: "", SizeField: 0, TypeOfData: "", Data: []byte{}}
			err = enc.Encode(&newMessage)
			if err != nil {
				fmt.Println(err, "inicio")
				continue
			}
			Channel = res1[1]
			fmt.Println("Suscribing to channel: " + Channel)
			suscribe(Channel, conn, enc)
		} else {
			continue
		}

	}
}
func send(name string, channel string, conn net.Conn, enc *gob.Encoder) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
		errorFile := []byte("Error")
		fileData := Message{"Error", channel, len(errorFile), "String", errorFile}
		err = enc.Encode(&fileData)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	res1 := strings.Split(name, "/")
	typeofData := "FileData"
	fileData := Message{res1[len(res1)-1], channel, len(data), typeofData, data}
	err = enc.Encode(&fileData)

	if err != nil {
		fmt.Println(err)
		return
	}
}
func suscribe(channel string, conn net.Conn, enc *gob.Encoder) {
	err := enc.Encode(&channel)
	if err != nil {
		fmt.Println(err)
	}
}

func receive(conn net.Conn, dec *gob.Decoder) {
	var message Message
	err := dec.Decode(&message)
	if err != nil {
		return
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
