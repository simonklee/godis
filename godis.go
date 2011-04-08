package godis

import (
    "bytes"
    "fmt"
    "log"
    "net"
    "os"
    "strconv"
)

const (
    LOG_CMD = false
)

type ReaderWriter interface {
    write(b []byte) (*conn, os.Error)
    read(c *conn) *Reply
}

type Client struct {
    Addr     string
    Db       int
    Password string
    pool     *pool
}

// writes a command a and returns single the Reply object
func Send(rw ReaderWriter, args ...[]byte) *Reply {
    c, err := rw.write(buildCmd(args...))

    if err != nil {
        return &Reply{Err: err}
    }

    return rw.read(c)
}

// uses reflection to create a bytestring of the name and args parameters, 
// then calls Send()
func SendIface(rw ReaderWriter, name string, args ...interface{}) *Reply {
    buf := make([][]byte, len(args)+1)
    buf[0] = []byte(name)

    for i, arg := range args {
        switch v := arg.(type) {
        case []byte:
            buf[i+1] = v
        case string:
            buf[i+1] = []byte(v)
        default:
            buf[i+1] = []byte(fmt.Sprint(arg))
        }
    }

    return Send(rw, buf...)
}

// creates a bytestring of the name and args parameters, then calls Send()
func SendStr(rw ReaderWriter, name string, args ...string) *Reply {
    buf := make([][]byte, len(args)+1)
    buf[0] = []byte(name)

    for i, arg := range args {
        buf[i+1] = []byte(arg)
    }

    return Send(rw, buf...)
}

// takes a [][]byte and returns a redis command formatted using
// the unified request protocol
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

func New(addr string, db int, password string) *Client {
    if addr == "" {
        addr = "127.0.0.1:6379"
    }

    return &Client{Addr: addr, Db: db, Password: password, pool: newPool()}
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
        co := newConn(conn)
        _, err = co.rwc.Write(buildCmd([]byte("SELECT"), []byte(strconv.Itoa(c.Db))))

        if err != nil {
            return nil, err
        }

        r := co.readReply()
        if r.Err != nil {
            return nil, r.Err
        }
    }

    if c.Password != "" {
        co := newConn(conn)
        _, err := co.rwc.Write(buildCmd([]byte("AUTH"), []byte(c.Password)))

        if err != nil {
            return nil, err
        }

        r := co.readReply()
        if r.Err != nil {
            return nil, r.Err
        }
    }

    return conn, err
}

func (c *Client) read(conn *conn) *Reply {
    reply := conn.readReply()
    c.pool.push(conn)
    return reply
}

func (c *Client) write(cmd []byte) (conn *conn, err os.Error) {
    conn = c.pool.pop()

    defer func() {
        if err != nil {
            log.Printf("ERR (%v), conn: %q", err, conn)
            c.pool.push(nil)
        }
    }()

    if conn == nil {
        rwc, err := c.createConn()

        if err != nil {
            return nil, err
        }

        conn = newConn(rwc)
        connCount++
    }

    _, err = conn.buf.Write(cmd)
    conn.buf.Flush()
    return conn, err
}

type Pipe struct {
    *Client
    conn       *conn
    appendMode bool
}

// Pipe implements the ReaderWriter interface, can be used with all commands
func NewPipe(addr string, db int, password string) *Pipe {
    return &Pipe{New(addr, db, password), nil, true}
}

// will return the Reply object made in the order commands where made
func (p *Pipe) GetReply() *Reply {
    if p.appendMode {
        p.appendMode = false
    }

    return p.read(p.conn)
}

func (p *Pipe) read(conn *conn) *Reply {
    if p.appendMode {
        return &Reply{}
    }

    if p.conn.buf.Available() > 0 {
        p.conn.buf.Flush()
    }

    reply := p.conn.readReply()

    if reply.Err != nil {
        // TODO: find out when there are no more replies
        p.pool.push(p.conn)
        p.conn = nil
        p.appendMode = true
    }

    return reply
}

func (p *Pipe) write(cmd []byte) (*conn, os.Error) {
    var err os.Error

    if p.conn == nil {
        c := p.pool.pop()

        defer func() {
            if err != nil {
                p.pool.push(nil)
            }
        }()

        if c == nil {
            rwc, err := p.createConn()

            if err != nil {
                return nil, err
            }

            c = newConn(rwc)
            connCount++
        }

        p.conn = c
    }

    _, err = p.conn.buf.Write(cmd)
    return p.conn, err
}
