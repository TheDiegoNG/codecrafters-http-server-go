package main

import (
    "bufio"
    "bytes"
    "compress/gzip"
    "fmt"
    "net"
    "os"
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
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    myConn := Connection{
        Conn: conn,
    }

    reader := bufio.NewReader(conn)
    request, err := parseRequest(reader)
    fmt.Println("Request: ", request)

    if err != nil {
        fmt.Println("Error parsing the request")
        myConn.writeResponse("HTTP/1.1 500 Internal Server Error\r\n\r\n")
    }

    pathParts := strings.Split(request.Path, "/")
    fmt.Println("Path Parts: ", pathParts)
    pathCommand := pathParts[1]
    switch {
    case request.Method == "POST":
        dir := os.Args[2]
        fileName := pathParts[2]
        perm := os.FileMode(0644)
        fmt.Println("Complete Path: ", dir + fileName)
        if _, err := os.Stat(dir); os.IsNotExist(err) {
            if err := os.MkdirAll(dir, 0755); err != nil {
                fmt.Println("Error creating directory: ", err)
                response := request.HttpVersion + " 501 Internal Server Error\r\n\r\n"
                myConn.writeResponse(response)
            }
        }
        err := os.WriteFile(dir + fileName, []byte(request.Body), perm)
        if err != nil {
            response := request.HttpVersion + " 404 Not Found\r\n\r\n"
            myConn.writeResponse(response)
        } else {
            response := request.HttpVersion + " 201 Created\r\n\r\n"
            myConn.writeResponse(response)
        }
    case strings.ToLower(pathCommand) == "echo":
        contEncoding := ""
        var b bytes.Buffer
        encoders := strings.Split(request.Headers["Accept-Encoding"], ", ")
        for _, x := range encoders {
            if x == "gzip" {
                contEncoding = "Content-Encoding: gzip\r\n"
                gz := gzip.NewWriter(&b)
                if _, err := gz.Write([]byte(pathParts[2])); err != nil {
                    fmt.Println("Error compressing the data: ", err)
                    response := request.HttpVersion + " 501 Internal Server Error\r\n\r\n"
                    myConn.writeResponse(response)
                }
                gz.Close()
                response := request.HttpVersion + " 200 OK\r\n" +
                "Content-Type: text/plain\r\n" +
                contEncoding +
                "Content-Length: " + strconv.Itoa(len(b.String())) + "\r\n\r\n" +
                b.String()
                myConn.writeResponse(response)
                return
            }
        }
        response := request.HttpVersion + " 200 OK\r\n" +
        "Content-Type: text/plain\r\n" +
        contEncoding +
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

func parseRequest(reader *bufio.Reader) (*HttpRequest, error) {
    var req HttpRequest
    req.Headers = make(map[string]string)
    line, err := reader.ReadString('\n')
    if err != nil {
        return nil, err
    }

    parts := strings.Fields(line)

    if len(parts) < 3 {
        return nil, fmt.Errorf("Invalid request line")
    }
    req.Method = parts[0]
    req.Path = parts[1]
    req.HttpVersion = parts[2]

    for {
        line, err = reader.ReadString('\n')
        if err != nil {
            return nil, err
        }
        line = strings.TrimSpace(line)
        if line == "" {
            break
        }
        headers := strings.SplitN(line, ": ", 2)
        if len(headers) < 2 {
            return nil, fmt.Errorf("Invalid header line")
        }
        req.Headers[headers[0]] = headers[1]
    }

    if contentLength, ok := req.Headers["Content-Length"]; ok {
        if length, err := strconv.Atoi(contentLength); err == nil {
            body := make([]byte, length)
            _, err := reader.Read(body)
            if err != nil {
                return nil, err
            }
            req.Body = string(body)
        }
    }
    return &req, nil
}

func (conn Connection) writeResponse(response string) {
    fmt.Println("Init writeResponse. Response: ", response)
            _, errConn := conn.Conn.Write([]byte(response))
            if errConn != nil {
                fmt.Println("Error accepting connection: ", errConn.Error())
            }
}
