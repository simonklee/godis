package main

import (
    "net"
    "fmt"
    "os"
    "bufio"
    "bytes"
    "strings"
    "strconv"
)

type Client struct {
    host string
    port int
    db int
}

func bytesCommand(cmd string, args ...string) []byte {
    buf := bytes.NewBufferString(fmt.Sprintf("*%d\r\n$%d\r\n%s\r\n", len(args) + 1, len(cmd), cmd))
    for _, arg := range args {
        buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
    }    
    return buf.Bytes()
}

func read(reader *bufio.Reader) ([]byte, os.Error) {
    var res string
    var err os.Error

    for {
        res, err = reader.ReadString('\n')
        if err != nil {
            return nil, err
        }
        break
    }
    res_type := res[0]
    res = strings.TrimSpace(res[1:])

    switch res_type {
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
            l, _ := strconv.Atoi(string(res[0]))
            fmt.Println(l)
    }

    fmt.Printf("%q\n", res);
    return []byte(res), nil
}

func write(con net.Conn, cmd []byte) (*bufio.Reader, os.Error) {
    _, err := con.Write(cmd)
    if err != nil {
        return nil, os.NewError("Error writing cmd " + err.String())
    }
    
    return bufio.NewReader(con), nil
}

func (client *Client) send(cmd string, args...string) (data []byte, err os.Error) {
    var addrString string = fmt.Sprintf("%s:%d", client.host, client.port)

    addr, err := net.ResolveTCPAddr(addrString)
    if err != nil {
        return nil, os.NewError("Error resolving Redis TCP addr")
    }

    con, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        return nil, os.NewError("Error connection to Redis at " + addr.String())
    }

    reader, err := write(con, bytesCommand(cmd, args...))
    if err != nil {
        return nil, err
    }

    data, err = read(reader) 
    con.Close()

    return
}

func main() {
    var client Client = Client{"127.0.0.1", 6379, 0} 

    // var enc_set []byte = bytesCommand("SET", "key", "hello")
    // fmt.Printf("%q\n", enc_set)

    // var enc_get []byte = bytesCommand("GET", "key")
    // fmt.Printf("%q\n", enc_get)

    // client.write(enc_set)
    // client.write(enc_get)
    // client.write(bytesCommand("RPUSH", "keylist", "1"))
    // client.write(bytesCommand("GET", "keylist"))
    // client.write(bytesCommand("GET", "nonexistant"))
    client.send("LRANGE", "keylist", "0", "2")
    client.send("KEYS", "*")
}
