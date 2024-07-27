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

type Connection struct {
    Conn net.Conn
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

    var myConn Connection

    myConn.Conn = conn
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
        myConn.writeResponse(response)
    case strings.ToLower(pathCommand) == "user-agent":
        response := request.HttpVersion + " 200 OK\r\n" +
        "Content-Type: text/plain\r\n" +
        "Content-Length: " + strconv.Itoa(len(request.Headers["User-Agent"])) + "\r\n\r\n" +
        request.Headers["User-Agent"]
        myConn.writeResponse(response)
    case strings.ToLower(pathCommand) == "files":
        dir := os.Args[2]
        fileName := pathParts[2]
        fmt.Println("Dir: ", dir)
        fmt.Println("FileName: ", fileName)
        data, err := os.ReadFile(dir + fileName)
        fmt.Println(data)
        if err != nil {
            response := request.HttpVersion + " 404 Not Found\r\n\r\n"
            myConn.writeResponse(response)
        } else {
            response := request.HttpVersion + " 200 OK\r\n" +
            "Content-Type: application/octet-stream\r\n" +
            "Content-Length: " + strconv.Itoa(len(data)) + "\r\n\r\n" +
            string(data)
            myConn.writeResponse(response)
        }
    case request.Path == "/":
        response := request.HttpVersion + " 200 OK\r\n\r\n"
        myConn.writeResponse(response)
    default:
        response := request.HttpVersion + " 404 Not Found\r\n\r\n"
        myConn.writeResponse(response)
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

func (conn Connection) writeResponse(response string) {
            _, errConn := conn.Conn.Write([]byte(response))
            if errConn != nil {
                fmt.Println("Error accepting connection: ", errConn.Error())
                os.Exit(1)
            }
}
