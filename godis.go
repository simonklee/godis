package godis

import (
    "net"
    "fmt"
    "os"
    "bufio"
    "bytes"
    "strings"
    "strconv"
)

type pool struct {
    free chan *net.TCPConn
}

type Client struct {
    Host string 
    Port int 
    Db int 
    Password string
    pool *pool
}

func log(args ...interface{}) {
    fmt.Printf("DEBUG: ")
    fmt.Println(args...)
}

func command(cmd string, args ...string) []byte {
    buf := bytes.NewBufferString(fmt.Sprintf("*%d\r\n$%d\r\n%s\r\n", len(args) + 1, len(cmd), cmd))
    for _, arg := range args {
        buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
    }    
    return buf.Bytes()
}

func read(head *bufio.Reader) (interface{}, os.Error) {
    var res string
    var err os.Error

    for {
        res, err = head.ReadString('\n')
        if err != nil {
            return nil, err
        }
        break
    }
    res_type := res[0]
    res = strings.TrimSpace(res[1:])

    fmt.Printf("%c\n", res_type)
    switch res_type {
        case '+':
            log(res)
            return res, nil
        case '-':
            log(res)
            return nil, os.NewError(res)
        case ':':
            n, err := strconv.Atoi64(res)
            log(n)
            return n, err
        case '$':
            l, _ := strconv.Atoi(res)
            if l == -1 {
                return nil, os.NewError("Key does not exist")
            }

            l += 2 
            data := make([]byte, l)

            n, err := head.Read(data)
            if n != l || err != nil {
                if n != l {
                    err = os.NewError("Len mismatch")
                }
                return nil, err
            }

            log("bulk-len: " + strconv.Itoa(l))
            log("bulk-value: " + string(data))
            fmt.Printf("%q\n", data)

            return data[:l - 2], nil
        case '*':
            l, _ := strconv.Atoi(string(res[0]))
            log("multi-bulk-len: " + strconv.Itoa(l))
            var data = make([][]byte, l)
            for i := 0; i < l; i++ {
                d, err := read(head)
                if err != nil {
                    log("returned with error")
                    return nil, err
                }
                data[i] = d.([]byte)
            }

            fmt.Printf("%q\n", data)
            return data, nil
    }
    return nil, os.NewError("Undefined redis response") 
}

func write(conn *net.TCPConn, cmd string, args ...string) (err os.Error) {
    _, err = conn.Write(command(cmd, args...))
    return
}

func (p *pool) pop() (*net.TCPConn, os.Error) { 
    return nil, nil
}

func (p *pool) push(*net.TCPConn) {
}


func (client *Client) connect() (*net.TCPConn, os.Error) {
    if client.pool == nil {
        //client.pool = pool{}
    }

    addrString := fmt.Sprintf("%s:%d", client.Host, client.Port)
    addr, err := net.ResolveTCPAddr(addrString)
    if err != nil {
        return nil, os.NewError("Error resolving Redis TCP addr")
    }

    conn, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        return nil, os.NewError("Error connection to Redis at " + addr.String())
    }
    return conn, err
}

func (client *Client) Send(cmd string, args...string) (interface{}, os.Error) {
    conn, err := client.connect()
    if err != nil {
        return nil, err
    }

    err = write(conn, cmd, args...)
    data, err := read(bufio.NewReader(conn)) 
    conn.Close()

    return data, err
}
