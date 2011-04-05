package godis

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "strconv"
)

const (
    MaxClientConn = 1
    LOG_CMD       = true
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
    delim   = []byte{cr, ln}
    connCount int
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

func (e Elem) Float64() float64 {
    v, _ := strconv.Atof64(string([]byte(e)))
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

func (r *Reply) IntArray() []int64 {
    buf := make([]int64, len(r.Elems))

    for i, v := range r.Elems {
        v, _ := strconv.Atoi64(v.Elem.String())
        buf[i] = v
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
        rr := parseResponse(reader)

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

func parseResponse(reader *bufio.Reader) *Reply {
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

func bufferWrite(w io.Writer, cmd []byte) os.Error {
    _, err := w.Write(cmd)
    return err
}

func buildCmd(args ...[]byte) []byte {
    buf := bytes.NewBuffer(nil)

    buf.WriteByte(star)
    buf.WriteString(strconv.Itoa(len(args)))
    buf.Write(delim)

    for _, arg := range args {
        buf.WriteByte(dollar)
        buf.WriteString(strconv.Itoa(len(arg)))
        buf.Write(delim)
        buf.Write(arg)
        buf.Write(delim)
    }

    if LOG_CMD {
        log.Printf("GODIS: %q", string(buf.Bytes()))
    }

    return buf.Bytes()
}

type Reader interface {
    read(c *net.TCPConn) *Reply
}

type Writer interface {
    write(b []byte) (*net.TCPConn, os.Error)
}

type ReaderWriter interface {
    Reader
    Writer
}

// send writes a message and returns the reply
func Send(rw ReaderWriter, args ...[]byte) *Reply {
    c, err := rw.write(buildCmd(args...))

    if err != nil {
        return &Reply{Err: err}
    }

    return rw.read(c)
}

// uses reflection to create a bytestring of args, then calls Send()
func SendIface(rw ReaderWriter, name string, args ...interface{}) *Reply {
    buf := make([][]byte, len(args)+1)
    buf[0] = []byte(name)

    for i, arg := range args {
        buf[i+1] = []byte(fmt.Sprint(arg))
    }

    return Send(rw, buf...)
}

// creates a bytestring of the string parameters, then calls Send()
func SendStr(rw ReaderWriter, name string, args ...string) *Reply {
    buf := make([][]byte, len(args)+1)
    buf[0] = []byte(name)

    for i, arg := range args {
        buf[i+1] = []byte(arg)
    }

    return Send(rw, buf...)
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

func (c *Client) createConn() (conn *net.TCPConn, err os.Error) {
    addr, err := net.ResolveTCPAddr(c.Addr)

    if err != nil {
        return nil, os.NewError("ResolveAddr error for " + c.Addr)
    }

    conn, err = net.DialTCP("tcp", nil, addr)
    if err != nil {
        err = os.NewError("Connection error " + addr.String())
    }

    if c.Db != 0 {
        err = bufferWrite(conn, buildCmd([]byte("SELECT"), []byte(strconv.Itoa(c.Db))))

        if err != nil {
            return nil, err
        }

        r := parseResponse(bufio.NewReader(conn))
        if r.Err != nil {
            return nil, r.Err
        }
    }

    if c.Password != "" {
        err := bufferWrite(conn, buildCmd([]byte("AUTH"), []byte(c.Password)))

        if err != nil {
            return nil, err
        }

        r := parseResponse(bufio.NewReader(conn))
        if r.Err != nil {
            return nil, r.Err
        }
    }

    return conn, err
}

func (c *Client) read(conn *net.TCPConn) *Reply {
    reply := parseResponse(bufio.NewReader(conn))
    c.pool.Push(conn)
    return reply
}

func (c *Client) write(cmd []byte) (conn *net.TCPConn, err os.Error) {
    conn = c.pool.Pop()

    if conn == nil {
        if conn, err = c.createConn(); err != nil {
            return nil, err
        }
        connCount++
    }

    err = bufferWrite(conn, cmd)
    return conn, err
}

//type PipeClient struct {
//    *Client
//    writeOnly bool
//    c         *net.TCPConn
//    w         *bufio.Writer
//    r         *bufio.Reader
//}
//
//func NewPipe(addr string, db int, password string) *PipeClient {
//    return &PipeClient{New(addr, db, password), true, nil, nil, nil}
//}
//
//func (p *PipeClient) read(conn *net.TCPConn) *Reply {
//    if p.c == nil {
//        p.c = conn
//    }
//
//    if p.writeOnly {
//        return &Reply{}
//    }
//
//    if p.w != nil {
//        if p.w.Available() > 0 {
//            p.w.Flush()
//        }
//
//        p.w = nil
//    }
//
//    if p.r == nil {
//        p.r = bufio.NewReader(p.c)
//    }
//
//    reply := parseResponse(p.r)
//
//    if reply.Err != nil {
//        // check if timeout
//        p.pool.Push(p.c)
//        p.r = nil
//        p.c = nil
//    }
//
//    return reply
//}
//
//func (p *PipeClient) write(cmd []byte) (conn *net.TCPConn, err os.Error) {
//    if p.w == nil {
//        conn = p.pool.Pop()
//
//        if conn == nil {
//            conn, err = p.createConn(); if err != nil {
//                return nil, err
//            }
//        }
//
//        p.w = bufio.NewWriter(conn)
//    }
//
//    err = bufferWrite(conn, cmd)
//    return conn, err
//}
//
//func (p *PipeClient) GetReply() *Reply {
//    if p.writeOnly {
//        p.writeOnly = false
//    }
//    return p.read(p.c)
//}
//
