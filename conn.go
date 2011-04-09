package godis

import (
    "net"
    "bufio"
    "os"
    "log"
    "strconv"
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

const (
    MaxClientConn = 1
    LOG_CMD       = false
)

var (
    delim     = []byte{cr, ln}
    connCount int
    cmdCount  = map[byte]int{dollar: 0, colon: 0, minus: 0, plus: 0, star: 0}
)

type conn struct {
    rwc net.Conn
    buf *bufio.ReadWriter
}

type pool struct {
    free chan *conn
}

type Elem []byte

type Reply struct {
    conn  *conn
    Err   os.Error
    Elem  Elem
    Elems []*Reply
}

func newPool() *pool {
    p := pool{make(chan *conn, MaxClientConn)}

    for i := 0; i < MaxClientConn; i++ {
        p.free <- nil
    }

    return &p
}

func (p *pool) pop() *conn {
    return <-p.free
}

func (p *pool) push(c *conn) {
    p.free <- c
}

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

func (r *Reply) parseBulk(res []byte) {
    l, _ := strconv.Atoi(string(res))

    if l == -1 {
        return
    }

    l += 2 // make sure to read \r\n
    data := make([]byte, l)

    n, err := r.conn.buf.Read(data)

    // if we were unable to read all date from socket, try again
    if n != l && err == nil {
        more := make([]byte, l - n)

        if _, err := r.conn.buf.Read(more); err != nil {
            r.Err = err
            return
        }

        data = append(data[:n], more...)
    }

    if err != nil {
        r.Err = err
        return
    }

    l -= 2
    r.Elem = data[:l]

    if LOG_CMD {
        //log.Printf("CONN: read %d byte, bulk-data %q\n", l, data)
    }
}

func (r *Reply) parseMultiBulk(res []byte) {
    l, _ := strconv.Atoi(string(res))

    if l == -1 {
        r.Err = nil //os.NewError("nothing to read")
        return
    }

    r.Elems = make([]*Reply, l)

    for i := 0; i < l; i++ {
        rr := r.conn.readReply()

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
        //log.Printf("GODIS: %d == %d %q\n", l, len(r.Elems), r.Elems)
    }
}

func (c *conn) readReply() *Reply {
    r := new(Reply)
    r.conn = c
    res, err := c.buf.ReadBytes(ln)

    if err != nil {
        r.Err = err
        return r
    }

    typ := res[0]
    line := res[1 : len(res)-2]

    if LOG_CMD {
        cmdCount[typ]++
        //log.Printf("CONN: alloc new Reply for `%c`: %q\n", typ, r.conn)
    }

    switch typ {
    case minus:
        r.parseErr(line)
    case plus:
        r.parseStr(line)
    case colon:
        r.parseInt(line)
    case dollar:
        r.parseBulk(line)
    case star:
        r.parseMultiBulk(line)
    default:
        r.Err = os.NewError("Unknown response " + string(typ))
    }

    return r
}

func newConn(rwc net.Conn) *conn {
    br := bufio.NewReader(rwc)
    bw := bufio.NewWriter(rwc)

    return &conn{
        rwc: rwc,
        buf: bufio.NewReadWriter(br, bw),
    }
}
