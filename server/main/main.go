package main

import (
	"fmt"
	"net"
)

func main() {

	//server := serv.GnetKVStoreServer{}
	//server.Start(9000)
	//fmt.Println("server started")
	server, err := net.Listen("tcp", ":9001")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(conn.RemoteAddr())
		var buf = make([]byte, 1024)
		if _, err := conn.Read(buf); err != nil {
			fmt.Println(err)
		}
		fmt.Println(buf)

		var n int
		if n, err = conn.Write([]byte("hello")); err != nil {
			fmt.Println(err)
		}
		fmt.Println(n)
	}
}
