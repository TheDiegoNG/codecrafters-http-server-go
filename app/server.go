package main

import (
	"fmt"
	 "net"
	 "os"
    "bufio"
    "strings"
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

    handleConnection(conn)
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)

    requestLine, err := reader.ReadString('\n')

    if err != nil {
        fmt.Println("Error reading the requestLine")
        os.Exit(1)
    }

    requestLine = strings.TrimSpace(requestLine)
    fmt.Println("Request Line: ", requestLine)

    parts := strings.Split(requestLine, " ")
    if len(parts) != 3 {
        fmt.Println("Invalid request line")
            os.Exit(1)
    }

    // method := parts[0]
    path := parts[1]
    // httpVersion := parts[2]
    fmt.Println("Path: ", path)
    if path == "/" {
        _, errConn := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
        if errConn != nil {
            fmt.Println("Error accepting connection: ", errConn.Error())
            os.Exit(1)
        }
    } else {
        _, errConn := conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
        if errConn != nil {
            fmt.Println("Error accepting connection: ", errConn.Error())
            os.Exit(1)
        }
    }


}
