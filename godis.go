package godis

import (
    "fmt"
    "os"
    "bufio"
    "bytes"
    "strings"
    "strconv"
    "net"
    "log"
)

const (
    MaxClientConn = 5
    LOG_CMD = false
)

var (
    ConnCtr int
)

func newError(format string, args ...interface{}) os.Error {
    return os.NewError(fmt.Sprintf(format, args...))
}

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

type redisReadWriter struct {
    writer *bufio.Writer
    reader *bufio.Reader
}

func newRedisReadWriter(c *net.TCPConn) *redisReadWriter {
    return &redisReadWriter{bufio.NewWriter(c), bufio.NewReader(c)}
}

func (rw *redisReadWriter) errorReply(line string) (interface{}, os.Error) {
    if LOG_CMD {
        log.Println("GODIS: " + line)
    }

    if strings.HasPrefix(line, "ERR") {
        line = line[3:]
    }

    return nil, newError(line)
}

func (rw *redisReadWriter) singleReply(line string) (interface{}, os.Error) {
    if LOG_CMD {
        log.Println("GODIS: " + line)
    }

    return line, nil
}

func (rw *redisReadWriter) integerReply(line string) (interface{}, os.Error) {
    if LOG_CMD {
        log.Println("GODIS: " + line)
    }

    return strconv.Atoi64(line)
}

func (rw *redisReadWriter) bulkReply(line string) (interface{}, os.Error) {
    l, _ := strconv.Atoi(line)

    if l == -1 {
        return nil, nil
    }

    l += 2 // make sure to read \r\n
    data := make([]byte, l)

    n, err := rw.reader.Read(data)
    if n != l || err != nil {
        if n != l {
            err = newError("expected %d bytes got %d bytes", l, n)
        }
        return nil, err
    }
    l -= 2

    if LOG_CMD {
        log.Printf("GODIS: %d %q\n", l, data)
    }
    
    return data[:l], nil
}

func (rw *redisReadWriter) multiBulkReply(line string) (interface{}, os.Error) {
    l, _ := strconv.Atoi(line)

    if l == -1 {
        return nil, nil
    }

    var data = make([][]byte, l)

    for i := 0; i < l; i++ {
        v, err := rw.read()

        if err != nil {
            return nil, err
        }

        // key not found, ignore `nil` value
        if v == nil {
            i -= 1
            l -= 1
            continue
        }

        data[i] = v.([]byte)
    }

    if LOG_CMD {
        log.Printf("GODIS: %d %q\n", l, data)
    }

    return data[:l], nil
}

func (rw *redisReadWriter) read() (interface{}, os.Error) {
    res, err := rw.reader.ReadString('\n')

    if err != nil {
        return nil, err
    }

    typ := res[0]
    line := strings.TrimSpace(res[1:])

    if LOG_CMD {
        log.Printf("GODIS: %c\n", typ)
    }

    switch typ {
    case '-':
        return rw.errorReply(line)
    case '+':
        return rw.singleReply(line)
    case ':':
        return rw.integerReply(line)
    case '$':
        return rw.bulkReply(line)
    case '*':
        return rw.multiBulkReply(line)
    }

    return nil, newError("Unknown response ", string(typ))
}

func (rw *redisReadWriter) write(name []byte, args ...[]byte) os.Error {
    cmds := append([][]byte{name}, args...)
    buf := bytes.NewBuffer(nil)
    fmt.Fprintf(buf, "*%d\r\n", len(cmds))
    
    for _, v := range cmds {
        fmt.Fprintf(buf, "$%d\r\n%s\r\n", len(v), v)
    }

    if LOG_CMD {
        log.Println(buf)
    }

    if _, err := rw.writer.Write(buf.Bytes()); err != nil {
        return err
    }

    return rw.writer.Flush()
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

    rw := newRedisReadWriter(conn)
    if c.Db != 0 {
        err = rw.write([]byte("SELECT"), []byte(strconv.Itoa(c.Db)))
        defer rw.read()

        if err != nil {
            return nil, err
        }
    }

    if c.Password != "" {
        err = rw.write([]byte("AUTH"), []byte(c.Password))
        defer rw.read()

        if err != nil {
            return nil, err
        }
    }
    return conn, err
}

func (c *Client) SendStr(name string, args ...string) (interface{}, os.Error) {
    buf := make([][]byte, len(args))

    for i, v := range args {
        buf[i] = []byte(v)
    }

    return c.Send(name, buf...)
}

func (c *Client) Send(name string, args ...[]byte) (interface{}, os.Error) {
    conn := c.pool.Pop()
    
    if conn == nil {
        ConnCtr++
        var err os.Error

        conn, err = c.newConn()
        if err != nil {
            return nil, err
        }
    }
    
    rw := newRedisReadWriter(conn)

    if err := rw.write([]byte(name), args...); err != nil {
        return nil, err
    }

    d, e := rw.read()
    c.pool.Push(conn)
    return d, e
}
