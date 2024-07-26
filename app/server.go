package main

import (
	"fmt"
	 "net"
	 "os"
    "bufio"
    "strconv"
    "strings"
)

type HttpRequest struct {
    Method string
    Path string
    HttpVersion string
    Headers map[string]string
    Body string
}
func main() {

	 l, err := net.Listen("tcp", "0.0.0.0:4221")
	 if err != nil {
	 	fmt.Println("Failed to bind to port 4221")
	 	os.Exit(1)
	 }

    for {
        conn, errAcc := l.Accept()
        if errAcc != nil {
            fmt.Println("Error accepting connection: ", errAcc.Error())
            os.Exit(1)
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    scanner := bufio.NewScanner(conn)
    request, err := parseRequest(scanner)

    if err != nil {
        fmt.Println("Error parsing the request")
        os.Exit(1)
    }

    pathParts := strings.Split(request.Path, "/")
    pathCommand := pathParts[1]
    fmt.Println(request)
    fmt.Println(pathCommand)
    switch {
    case strings.ToLower(pathCommand) == "echo":
        response := request.HttpVersion + " 200 OK\r\n" +
        "Content-Type: text/plain\r\n" +
        "Content-Length: " + strconv.Itoa(len(pathParts[2])) + "\r\n\r\n" +
        pathParts[2]
        _, errConn := conn.Write([]byte(response))
        if errConn != nil {
            fmt.Println("Error accepting connection: ", errConn.Error())
            os.Exit(1)
        }
    case strings.ToLower(pathCommand) == "user-agent":
        response := request.HttpVersion + " 200 OK\r\n" +
        "Content-Type: text/plain\r\n" +
        "Content-Length: " + strconv.Itoa(len(request.Headers["User-Agent"])) + "\r\n\r\n" +
        request.Headers["User-Agent"]
        _, errConn := conn.Write([]byte(response))
        if errConn != nil {
            fmt.Println("Error accepting connection: ", errConn.Error())
            os.Exit(1)
        }
    case request.Path == "/":

        _, errConn := conn.Write([]byte(request.HttpVersion + " 200 OK\r\n\r\n"))
        if errConn != nil {
            fmt.Println("Error accepting connection: ", errConn.Error())
            os.Exit(1)
        }
    default:

        _, errConn := conn.Write([]byte(request.HttpVersion + " 404 Not Found\r\n\r\n"))
        if errConn != nil {
            fmt.Println("Error accepting connection: ", errConn.Error())
            os.Exit(1)
        }
    }
}

func parseRequest(scanner *bufio.Scanner) (*HttpRequest, error) {
    var req HttpRequest
    req.Headers = make(map[string]string)
    scanner.Scan()
    parts := strings.Split(scanner.Text(), " ")
    req.Method = parts[0]
    req.Path = parts[1]
    req.HttpVersion = parts[2]
    for i := 0; scanner.Scan(); i++ {
        headers := strings.Split(scanner.Text(), ": ")
        //If headers return < 2 then it's the body
        if len(headers) < 2 {
            req.Body = headers[0]
            break
        }
        req.Headers[headers[0]] = headers[1]
    }
    return &req, nil
}
