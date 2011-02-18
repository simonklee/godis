package main

import (
    "net"
    "fmt"
    "os"
    "bufio"
    "bytes"
    "strings"
)

type Client struct {
    host string
    port int
    db int
}

const (
    CRLF = "\n\r"
)

func read(reader *bufio.Reader) (string, os.Error) {
    var line string
    var err os.Error

    for {
        line, err = reader.ReadString('\n')
        if err != nil {
            return "", err
        }
        break
    }
    line = strings.TrimSpace(line)

    switch line[0] {
        case '+':
            fmt.Printf("single line\n")
        case '-':
            fmt.Printf("error\n")
        case ':':
            fmt.Printf("integer\n")
        case '$':
            fmt.Printf("bulk\n")
        case '*':
            fmt.Printf("multi-bulk\n")
    }

    fmt.Printf("%q\n", line);
    return line, nil
}

func (client *Client) send(cmd []byte) (string, os.Error) {
    var addrString string = fmt.Sprintf("%s:%d", client.host, client.port)
    addr, err := net.ResolveTCPAddr(addrString)
    if err != nil {
        fmt.Println("Error resolving TCP addr")
        os.Exit(1)
    }

    con, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        fmt.Println("Error connecting to Redis at", addr.String())
        os.Exit(1)
    }

    _, err = con.Write(cmd)
    if err != nil {
        fmt.Println("Error writing cmd", err.String())
        os.Exit(1)
    }
    
    reader := bufio.NewReader(con)
    return read(reader)
}

func byteCommand(cmd string, args ...string) []byte {
    buf := bytes.NewBufferString(fmt.Sprintf("*%d\r\n$%d\r\n%s\r\n", len(args) + 1, len(cmd), cmd))
    for _, arg := range args {
        buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
    }    
    return buf.Bytes()
}

func main() {
    var client Client = Client{"127.0.0.1", 6379, 0} 

    var enc_set []byte = byteCommand("SET", "key", "hello")
    fmt.Printf("%q\n", enc_set)

    var enc_get []byte = byteCommand("GET", "key")
    fmt.Printf("%q\n", enc_get)

    client.send(enc_set)
    client.send(enc_get)
    client.send(byteCommand("RPUSH", "keylist", "1"))
    client.send(byteCommand("GET", "keylist"))
    client.send(byteCommand("GET", "nonexistant"))
    client.send(byteCommand("LRANGE", "keylist", "0", "2"))
}
