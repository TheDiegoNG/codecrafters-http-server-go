package main

import (
	"fmt"
	 "net"
	 "os"
)

func main() {

	 l, err := net.Listen("tcp", "0.0.0.0:4221")
	 if err != nil {
	 	fmt.Println("Failed to bind to port 4221")
	 	os.Exit(1)
	 }

    conn, errAcc := l.Accept()
	 if errAcc != nil {
	 	fmt.Println("Error accepting connection: ", errAcc.Error())
	 	os.Exit(1)
	 }
    _, errConn := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
    if errConn != nil {
	 	fmt.Println("Error accepting connection: ", errAcc.Error())
	 	os.Exit(1)
    }
}
