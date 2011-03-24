package godis

import (
    "fmt"
    "os"
    "bufio"
    "bytes"
    "strings"
    "strconv"
    "net"
//    "log"
)

const (
    MaxClientConn = 5
)

var (
    ConnCtr int
)

type Pool struct {
    pool chan *net.TCPConn
}

func NewPool() *Pool {
    p := Pool{make(chan *net.TCPConn, MaxClientConn)}
    for i := 0; i < MaxClientConn; i++ {
        p.pool <- nil
    }
    return &p
}

func (p *Pool) Pop() *net.TCPConn {
    return <-p.pool
}

func (p *Pool) Push(c *net.TCPConn) {
    p.pool <- c
}

func newError(format string, args ...interface{}) os.Error {
    return os.NewError(fmt.Sprintf(format, args...))
}

func errorReply(line string) (interface{}, os.Error) {
    // log.Println("GODIS: " + res)
    if strings.HasPrefix(line, "ERR") {
        line = line[3:]
    }
    return nil, newError(line)
}

func singleReply(line string) (string, os.Error) {
    // log.Println("GODIS: " + res)
    return line, nil
}

func integerReply(line string) (int64, os.Error) {
    // log.Println("GODIS: " + res)
    return strconv.Atoi64(line)
}

func bulkReply(line string, head *bufio.Reader) ([]byte, os.Error) {
    l, _ := strconv.Atoi(line)
    if l == -1 {
        return nil, nil
    }

    l += 2 // make sure to read \r\n
    data := make([]byte, l)

    n, err := head.Read(data)
    if n != l || err != nil {
        if n != l {
            err = newError("expected %d bytes got %d bytes", l, n)
        }
        return nil, err
    }
    l -= 2
    //log.Println("GODIS: bulk-len: " + strconv.Itoa(l))
    //log.Println("GODIS: bulk-value: " + string(data))
    //log.Printf("GODIS: %q\n", data)
    return data[:l], nil
}

func multiBulkReply(line string, head *bufio.Reader) ([][]byte, os.Error) {
    l, _ := strconv.Atoi(line)

    if l == -1 {
        return nil, nil
    }

    var data = make([][]byte, l)
    for i := 0; i < l; i++ {
        d, err := readReply(head)
        if err != nil {
            return nil, err
        }
        data[i] = d.([]byte)
    }

    return data, nil
}

func readReply(head *bufio.Reader) (interface{}, os.Error) {
    res, err := head.ReadString('\n')
    if err != nil {
        return nil, err
    }
    typ := res[0]
    line := strings.TrimSpace(res[1:])

    switch typ {
    case '-':
        return errorReply(line)
    case '+':
        return singleReply(line)
    case ':':
        return integerReply(line)
    case '$':
        return bulkReply(line, head)
    case '*':
        return multiBulkReply(line, head)
    }
    return nil, newError("Unknown response " + string(typ))
}

func buildCommand(args ...string) []byte {
    cmd := bytes.NewBuffer(nil)
    fmt.Fprintf(cmd, "*%d\r\n", len(args))
    for _, arg := range args {
        fmt.Fprintf(cmd, "$%d\r\n%s\r\n", len(arg), arg)
    }
    return cmd.Bytes()
}

func write(conn *net.TCPConn, cmd string, args ...string) os.Error {
    cmds := append([]string{cmd}, args...)
    _, err := conn.Write(buildCommand(cmds...))
    return err
}

func read(conn *net.TCPConn) (interface{}, os.Error) {
    return readReply(bufio.NewReader(conn))
}

type Client struct {
    Addr     string
    Db       int
    Password string
    pool     *Pool
}

func New(addr string, db int, password string) *Client {
    var c Client
    c.Addr = addr
    if c.Addr == "" {
        c.Addr = "127.0.0.1:6379"
    }

    c.Db = db
    c.Password = password
    c.pool = NewPool()
    return &c
}

func (c *Client) newConn() (conn *net.TCPConn, err os.Error) {
    addr, err := net.ResolveTCPAddr(c.Addr)
    if err != nil {
        return nil, os.NewError("ResolveAddr error for " + c.Addr)
    }

    conn, err = net.DialTCP("tcp", nil, addr)
    if err != nil {
        err = os.NewError("Connection error " + addr.String())
    }

    if c.Db != 0 {
        if err = write(conn, "SELECT", strconv.Itoa(c.Db)); err != nil {
            return nil, err
        }
    }

    if c.Password != "" {
        if err = write(conn, "AUTH", c.Password); err != nil {
            return nil, err
        }
    }
    return conn, err
}

func (c *Client) Send(cmd string, args ...string) (interface{}, os.Error) {
    conn := c.pool.Pop()
    
    if conn == nil {
        ConnCtr++
        // log.Printf("creating conn %d", ConnCtr)
        var err os.Error
        conn, err = c.newConn()
        if err != nil {
            return nil, err
        }
    }

    if err := write(conn, cmd, args...); err != nil {
        return nil, err
    }

    d, e := read(conn)
    c.pool.Push(conn)
    return d, e
}
