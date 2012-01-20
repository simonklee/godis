package godis

import (
    "bufio"
    "bytes"
    "errors"
    "log"
    "net"

    "strconv"
    "strings"
)

const (
    // protocol bytes
    cr     byte = 13
    lf     byte = 10
    dollar byte = 36
    colon  byte = 58
    minus  byte = 45
    plus   byte = 43
    star   byte = 42

    debug = false
)

var (
    // Max connection pool size
    MaxClientConn = 2

    // protocol bytes
    delim = []byte{cr, lf}

    // misc
    connCount int
    cmdCount  = map[byte]int{dollar: 0, colon: 0, minus: 0, plus: 0, star: 0}
)

type conn struct {
    rwc net.Conn
    r   *bufio.Reader
    w   *bufio.Writer
}

type pool struct {
    free chan *conn
}

type Elem []byte

type Message struct {
    Channel string
    Elem    Elem
}

type Reply struct {
    conn  *conn
    Err   error
    Elem  Elem
    Elems []*Reply
}

// takes a [][]byte and returns a redis command formatted using
// the unified request protocol
func buildCmd(args [][]byte) []byte {
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

    if debug {
        log.Printf("GODIS: %q", string(buf.Bytes()))
    }

    return buf.Bytes()
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
    v, _ := strconv.ParseInt(string([]byte(e)), 10, 64)
    return v
}

func (e Elem) Float64() float64 {
    v, _ := strconv.ParseFloat(string([]byte(e)), 64)
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
        v, _ := strconv.ParseInt(v.Elem.String(), 10, 64)
        buf[i] = v
    }

    return buf
}

func (r *Reply) StringMap() map[string]string {
    arr := r.StringArray()
    n := len(arr)
    buf := make(map[string]string, n/2)

    if n%2 == 1 {
        return buf
    }

    for i := 0; i < n; i += 2 {
        buf[arr[i]] = arr[i+1]
    }

    return buf
}

func (r *Reply) Message() *Message {
    if len(r.Elems) < 3 {
        return nil
    }

    typ := r.Elems[0].Elem.String()

    switch typ {
    case "message":
        return &Message{r.Elems[1].Elem.String(), r.Elems[2].Elem}
    case "pmessage":
        return &Message{r.Elems[2].Elem.String(), r.Elems[3].Elem}
    }

    if strings.HasSuffix(typ, "subscribe") {
        return nil
    }

    return nil
}

func (r *Reply) parseErr(res []byte) {
    r.Err = errors.New(string(res))

    if debug {
        log.Println("GODIS-ERR: " + string(res))
    }
}

func (r *Reply) parseStr(res []byte) {
    r.Elem = res

    if debug {
        log.Println("GODIS-STR: " + string(res))
    }
}

func (r *Reply) parseInt(res []byte) {
    r.Elem = res

    if debug {
        log.Println("GODIS-INT: " + string(res))
    }
}

func (r *Reply) parseBulk(res []byte) {
    l, _ := strconv.Atoi(string(res))

    if l == -1 {
        if debug {
            log.Println("GODIS-BULK: Key does not exist")
        }

        r.Err = errors.New("Nonexisting key")
        return
    }

    l += 2 // make sure to read \r\n
    data := make([]byte, l)

    n, err := r.conn.r.Read(data)

    // if we were unable to read all date from socket, try again
    if n != l && err == nil {
        more := make([]byte, l-n)

        if _, err := r.conn.r.Read(more); err != nil {
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

    if debug {
        log.Printf("GODIS-BULK: read %d byte, bulk-data %q\n", l, data)
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
        }

        // key not found, ignore `nil` value
        //if rr.Elem == nil {
        //    i -= 1
        //    l -= 1

        //    if debug {
        //        log.Printf("KEY NOT FOUND")
        //    }

        //    continue
        //}

        r.Elems[i] = rr
    }

    // buffer is reduced to account for `nil` value returns
    r.Elems = r.Elems[:l]

    if debug {
        //log.Printf("GODIS: %d == %d %q\n", l, len(r.Elems), r.Elems)
    }
}

func (c *conn) readReply() *Reply {
    r := new(Reply)
    r.conn = c
    res, err := c.r.ReadBytes(lf)

    if err != nil {
        r.Err = err
        return r
    }

    typ := res[0]
    line := res[1 : len(res)-2]

    if debug {
        cmdCount[typ]++
        //log.Printf("CONN: alloc new Reply for `%c` %s\n", typ, string(line))
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
        r.Err = errors.New("Unknown response " + string(typ))
    }

    return r
}

func newConn(netTyp, addr string, db int, password string) (*conn, error) {
    rwc, err := net.Dial(netTyp, addr)

    if err != nil {
        return nil, errors.New("Connection error " + addr)
    }

    connCount++
    cc := &conn{
        rwc: rwc,
        r:   bufio.NewReader(rwc),
        w:   bufio.NewWriter(rwc),
    }

    err = cc.configConn(db, password)
    return cc, err
}

func (cc *conn) configConn(db int, password string) error {
    if password != "" {
        buf := [][]byte{[]byte("AUTH"), []byte(password)}
        _, err := cc.rwc.Write(buildCmd(buf))

        if err != nil {
            return err
        }

        r := cc.readReply()
        if r.Err != nil {
            return r.Err
        }
    }

    if db != 0 {
        buf := [][]byte{[]byte("SELECT"), []byte(strconv.Itoa(db))}
        _, err := cc.rwc.Write(buildCmd(buf))

        if err != nil {
            return err
        }

        r := cc.readReply()
        if r.Err != nil {
            return r.Err
        }
    }
    return nil
}
