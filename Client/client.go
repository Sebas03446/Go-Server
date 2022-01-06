package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"sync"
)

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

func main() {
	ClientInit()
}

func ClientInit() {
	var wg sync.WaitGroup
	var lock sync.RWMutex
	conn, err := net.Dial("tcp4", "localhost:8000")
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}
	//go receive(conn)
	defer conn.Close()
	var message string
	err = gob.NewDecoder(conn).Decode(&message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(message)
	msg := "The client is connected!"
	err = gob.NewEncoder(conn).Encode(msg)
	if err != nil {
		fmt.Println(err)
	}
	var operation string
	var name string
	var channel string
	for {
		operation := readFile(&operation)
		err = gob.NewEncoder(conn).Encode(&operation)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if operation == "send" {
			fmt.Println("Sending")
			wg.Add(1)
			fmt.Scan(&name)
			fmt.Scan(&channel)
			go send(name, channel, conn, &wg, &lock)

		} else if operation == "create" {
			wg.Add(1)
			fmt.Scan(&name)
			fmt.Println("Creating channel: paso")
			channel := Channel{name, []net.Conn{}}
			go create(channel, conn, &wg, &lock)
		} else if operation == "suscribe" {
			wg.Add(1)
			fmt.Scan(&channel)
			fmt.Println("Suscribing to channel: " + channel)
			go suscribe(channel, conn, &wg, &lock)

		}
		/* var fileData FileData
		err = gob.NewDecoder(conn).Decode(&fileData)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ioutil.WriteFile(fileData.Name, fileData.Data, 0644)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		} */
		wg.Add(1)
		go receive(conn, &wg, &lock)

	}
}
func readFile(operation *string) string {
	fmt.Scan(operation)
	return *operation
}
func send(name string, channel string, conn net.Conn, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	lock.Lock()
	fmt.Println(channel)
	data, err := ioutil.ReadFile(name)
	fmt.Println(string(data))
	if err != nil {
		fmt.Println(err)
	}
	res1 := strings.Split(name, "/")
	fileData := FileData{res1[len(res1)-1], channel, len(data), data}
	err = gob.NewEncoder(conn).Encode(fileData)

	if err != nil {
		fmt.Println(err)
		return
	}
	lock.Unlock()
}
func suscribe(channel string, conn net.Conn, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	lock.Lock()
	err := gob.NewEncoder(conn).Encode(&channel)
	if err != nil {
		fmt.Println(err)
	}
	lock.Unlock()
}
func create(channel Channel, conn net.Conn, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	lock.Lock()
	err := gob.NewEncoder(conn).Encode(channel)
	//conn.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	lock.Unlock()
}

func receive(conn net.Conn, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	lock.RLock()
	fmt.Println("Receiving")
	info, _ := ioutil.ReadAll(conn)
	fmt.Println(string(info))
	/*var typeData string
	//conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	err := gob.NewDecoder(conn).Decode(&typeData)
	if err != nil {
		fmt.Println(err)
		lock.RUnlock()
		return
	}
	if typeData == "FileData" {
		fmt.Println("Receiving FileData")
		var fileData FileData
		err := gob.NewDecoder(conn).Decode(&fileData)
		if err != nil {
			fmt.Println(err)
		}
		ioutil.WriteFile(fileData.Name, fileData.Data, 0644)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
		fmt.Println("Received file: " + fileData.Name)
	} else {
		var response string
		err := gob.NewDecoder(conn).Decode(&response)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(response)
	}*/

	lock.RUnlock()
}
