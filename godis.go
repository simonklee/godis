package godis

import (
    "fmt"
    "os"
    "bufio"
    "bytes"
    "strconv"
    "net"
    "log"
    "io"
)

const (
    MaxClientConn = 5
    LOG_CMD       = false
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

type Elem []byte

func (e Elem) String() string {
    return string([]byte(e))
}

func (e Elem) Int64() int64 {
    v, _ := strconv.Atoi64(string([]byte(e)))
    return v
}

type Reply struct {
    Err   os.Error
    Elem  Elem
    Elems []*Reply
}

func (r *Reply) Len() int {
    return len(r.Elems)
}

func (r *Reply) Strings() []string {
    buf := make([]string, r.Len())

    for i, v := range r.Elems {
        buf[i] = v.Elem.String()
    }

    return buf
}

func (r *Reply) errorReply(res []byte) {
    r.Err = os.NewError(string(res))

    if LOG_CMD {
        log.Println("GODIS: " + string(res))
    }
}

func (r *Reply) singleReply(res []byte) {
    r.Elem = res

    if LOG_CMD {
        log.Println("GODIS: " + string(res))
    }
}

func (r *Reply) integerReply(res []byte) {
    r.Elem = res

    if LOG_CMD {
        log.Println("GODIS: " + string(res))
    }
}

func (r *Reply) bulkReply(reader *bufio.Reader, res []byte) {
    l, _ := strconv.Atoi(string(res))

    if l == -1 {
        return 
    }

    lr := io.LimitReader(reader, int64(l))
    buf := bytes.NewBuffer(make([]byte, 0, l))
    n, err := buf.ReadFrom(lr)

    if err == nil {
        _, err = reader.ReadBytes(ln)
    }

    if n != int64(l) {
        log.Println(n, l)
    }

    if LOG_CMD {
        log.Printf("GODIS: %d %q\n", l, buf)
    }

    r.Elem = buf.Bytes()
}

func (r *Reply) multiBulkReply(reader *bufio.Reader, res[]byte) {
    l, _ := strconv.Atoi(string(res))

    if l == -1 {
        r.Err = nil //os.NewError("nothing to read")
        return
    }

    r.Elems = make([]*Reply, l)

    for i := 0; i < l; i++ {
        rr := read(reader)

        if rr.Err != nil {
            r.Err = rr.Err
            return 
        }

        // key not found, ignore `nil` value
        if rr.Elem == nil {
            i -= 1
            l -= 1
            continue
        }

        r.Elems[i] = rr
    }

    if LOG_CMD {
        log.Printf("GODIS: %d %q\n", l, r.Elems)
    }
}

func read(reader *bufio.Reader) *Reply {
    reply := new(Reply)
    res, err := reader.ReadBytes(ln)

    if err != nil {
       reply.Err = err
       return reply
    }

    typ := res[0]
    line := res[1:len(res) - 2]

    if LOG_CMD {
        log.Printf("GODIS: %c\n", typ)
    }

    switch typ {
    case minus:
        reply.errorReply(line)
    case plus:
        reply.singleReply(line)
    case colon:
        reply.integerReply(line)
    case dollar:
        reply.bulkReply(reader, line)
    case star:
        reply.multiBulkReply(reader, line)
    default:
        reply.Err = os.NewError("Unknown response " + string(typ))
    }

    return reply 
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

func (c *Client) Send(name string, args ...interface{}) *Reply {
    conn := c.pool.Pop()

    if conn == nil {
        ConnCtr++
        var err os.Error
        conn, err = c.newConn()

        if err != nil {
            return &Reply{Err: err}
        }
    }

    if err := write(conn, name, args...); err != nil {
        return &Reply{Err: err}
    }

    reply := read(bufio.NewReader(conn))
    c.pool.Push(conn)
    return reply
}
