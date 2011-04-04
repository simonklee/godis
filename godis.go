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
    MaxClientConn = 1
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

func (e Elem) Bytes() []byte {
    return []byte(e)
}

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

func (r *Reply) BytesArray() [][]byte {
    buf := make([][]byte, len(r.Elems))

    for i, v := range r.Elems {
        buf[i] = v.Elem
    }

    return buf
}

func (r *Reply) StringArray() []string {
    buf := make([]string, len(r.Elems))

    for i, v := range r.Elems {
        buf[i] = v.Elem.String()
    }

    return buf
}

func (r *Reply) parseErr(res []byte) {
    r.Err = os.NewError(string(res))

    if LOG_CMD {
        log.Println("GODIS: " + string(res))
    }
}

func (r *Reply) parseStr(res []byte) {
    r.Elem = res

    if LOG_CMD {
        log.Println("GODIS: " + string(res))
    }
}

func (r *Reply) parseInt(res []byte) {
    r.Elem = res

    if LOG_CMD {
        log.Println("GODIS: " + string(res))
    }
}

func (r *Reply) parseBulk(reader *bufio.Reader, res []byte) {
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

    r.Elem = buf.Bytes()

    if LOG_CMD {
        log.Printf("G: %d %q %q\n", l, buf, buf.Bytes())
    }
}

func (r *Reply) parseMultiBulk(reader *bufio.Reader, res []byte) {
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

    // buffer is reduced to account for `nil` value returns
    r.Elems = r.Elems[:l]

    if LOG_CMD {
        log.Printf("GODIS: %d == %d %q\n", l, len(r.Elems), r.Elems)
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
    line := res[1 : len(res)-2]

    if LOG_CMD {
        log.Printf("GODIS: %c\n", typ)
    }

    switch typ {
    case minus:
        reply.parseErr(line)
    case plus:
        reply.parseStr(line)
    case colon:
        reply.parseInt(line)
    case dollar:
        reply.parseBulk(reader, line)
    case star:
        reply.parseMultiBulk(reader, line)
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

func write(w io.Writer, name string, args ...interface{}) os.Error {
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
        log.Printf("GODIS: %q", string(buf.Bytes()))
    }

    if _, err := w.Write(buf.Bytes()); err != nil {
        return err
    }

    return nil
}

type Writer interface {
    Write(name string, args ...interface{}) *Reply
}

type Reader interface {
    Read() *Reply
}

func Send(w Writer, name string, args ...interface{}) *Reply {
    return w.Write(name, args...)
}

func ReadReply(r Reader) *Reply {
    return r.Read()
}

type Client struct {
    Addr     string
    Db       int
    Password string
    pool     *Pool
}

func New(addr string, db int, password string) *Client {
    if addr == "" {
        addr = "127.0.0.1:6379"
    }

    return &Client{Addr: addr, Db: db, Password: password, pool: NewPool()}
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

func (c *Client) Write(name string, args ...interface{}) *Reply {
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

type PipeClient struct {
    *Client
    c  *net.TCPConn
    w  *bufio.Writer
    r  *bufio.Reader
}

func NewPipe(addr string, db int, password string) *PipeClient {
    return &PipeClient{New(addr, db, password), nil, nil, nil}
}

func (p *PipeClient) Write(name string, args ...interface{}) *Reply {
    if p.w == nil {
        conn := p.pool.Pop()

        if conn == nil {
            var err os.Error
            conn, err = p.newConn()

            if err != nil {
                return &Reply{Err: err}
            }
        }

        p.w = bufio.NewWriter(conn)
        p.c = conn
    }


    if err := write(p.w, name, args...); err != nil {
        return &Reply{Err: err}
    }

    return &Reply{}
}

func (p *PipeClient) Read() *Reply {
    if p.w != nil {
        if p.w.Available() > 0 { 
            p.w.Flush()
        }

        p.w = nil
    }

    if p.r == nil {
        p.r = bufio.NewReader(p.c)
        p.c.SetReadTimeout(1e8) // 100ms
    }

    reply := read(p.r)

    if reply.Err != nil {
        // check if timeout
        p.pool.Push(p.c)
        p.c = nil
        p.r = nil
    }

    return reply
}
