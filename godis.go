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
    "io"
)

const (
    MaxClientConn = 5
    LOG_CMD = false
)

// protocol bytes
const (
    cr     byte = 13
    ln     byte = 10
    dollar byte = 36
    colon  byte = 58
    minus  byte = 45
    plus   byte = 43
    star   byte = 42
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

func errorReply(line string) (interface{}, os.Error) {
    if LOG_CMD {
        log.Println("GODIS: " + line)
    }

    if strings.HasPrefix(line, "ERR") {
        line = line[3:]
    }

    return nil, newError(line)
}

func singleReply(line string) (interface{}, os.Error) {
    if LOG_CMD {
        log.Println("GODIS: " + line)
    }

    return line, nil
}

func integerReply(line string) (interface{}, os.Error) {
    if LOG_CMD {
        log.Println("GODIS: " + line)
    }

    return strconv.Atoi64(line)
}

func bulkReply(reader *bufio.Reader, line string) (interface{}, os.Error) {
    l, _ := strconv.Atoi(line)

    if l == -1 {
        return nil, nil
    }

    r := io.LimitReader(reader, int64(l))
    buf := bytes.NewBuffer(make([]byte, 0, l))
    n, err := buf.ReadFrom(r)

    if err == nil {
        _, err = reader.ReadBytes(ln)
    }

    if n != int64(l) {
        log.Println(n, l)
    }

    if LOG_CMD {
        log.Printf("GODIS: %d %q\n", l, buf)
    }

    return buf.Bytes(), err
}

func multiBulkReply(reader *bufio.Reader, line string) (interface{}, os.Error) {
    l, _ := strconv.Atoi(line)

    if l == -1 {
        return nil, nil
    }

    var data = make([][]byte, l)

    for i := 0; i < l; i++ {
        v, err := read(reader)

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

func read(reader *bufio.Reader) (interface{}, os.Error) {
    res, err := reader.ReadBytes(ln)

    if err != nil {
        return nil, err
    }

    typ := res[0]
    line := strings.TrimSpace(string(res[1:]))

    if LOG_CMD {
        log.Printf("GODIS: %c\n", typ)
    }

    switch typ {
    case minus:
        return errorReply(line)
    case plus:
        return singleReply(line)
    case colon:
        return integerReply(line)
    case dollar:
        o, e := bulkReply(reader, line)
        if e != nil && LOG_CMD {
            l, _ := strconv.Atoi(line)
            log.Printf("typ: %c, res: %q, line: %q line-len(%d)\n", typ, res, line, l)
        }
        return o, e

    case star:
        o, e := multiBulkReply(reader, line)
        if e != nil && LOG_CMD {
            l, _ := strconv.Atoi(line)
            log.Printf("typ: %c, res: %q, line: %q line-len(%d)\n", typ, res, line, l)
        }
        return o, e
    }

    return nil, newError("Unknown response ", string(typ))
}

func appendCmd(buf *bytes.Buffer, a []byte) {
    buf.WriteByte(dollar)
    buf.WriteString(strconv.Itoa(len(a)))
    buf.WriteByte(cr)
    buf.WriteByte(ln)
    buf.Write(a)
    buf.WriteByte(cr)
    buf.WriteByte(ln)
}

func write(conn *net.TCPConn, name string, args ...interface{}) os.Error {
    n := len(args)
    buf := bytes.NewBuffer(nil)

    buf.WriteByte(star)
    buf.WriteString(strconv.Itoa(n + 1))
    buf.WriteByte(cr)
    buf.WriteByte(ln)

    appendCmd(buf, []byte(name))

    for i := 0; i < n; i++ {
        appendCmd(buf, []byte(fmt.Sprint(args[i])))
    }

    if LOG_CMD {
        log.Println("GODIS: " + string(buf.Bytes()))
    }

    if _, err := conn.Write(buf.Bytes()); err != nil {
        return err
    }

    return nil
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
        err = write(conn, "SELECT", c.Db)
        read(bufio.NewReader(conn))
        //defer rw.read()

        if err != nil {
            return nil, err
        }
    }

    if c.Password != "" {
        err = write(conn, "AUTH", c.Password)
        read(bufio.NewReader(conn))
        //defer rw.read()

        if err != nil {
            return nil, err
        }
    }
    return conn, err
}

func (c *Client) Send(name string, args ...interface{}) (interface{}, os.Error) {
    conn := c.pool.Pop()
    
    if conn == nil {
        ConnCtr++
        var err os.Error

        conn, err = c.newConn()
        if err != nil {
            return nil, err
        }
    }

    if err := write(conn, name, args...); err != nil {
        return nil, err
    }

    d, e := read(bufio.NewReader(conn))
    c.pool.Push(conn)
    return d, e
}
