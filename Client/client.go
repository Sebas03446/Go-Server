package client

import (
	"bufio"
	"fmt"
	"net"
)

func ClientInit() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(netData)
	}
}
