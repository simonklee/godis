package godis

import (
    "errors"
    "io"
    "log"
    "strconv"
)

var (
    debug       = false
    ErrProtocol = errors.New("godis: protocol error")
)

func (r *Reply) parseErr(res []byte) {
    r.Err = errors.New(string(res))

    if debug {
        log.Println("-ERR: " + string(res))
    }
}

func (r *Reply) parseStr(res []byte) {
    b := make([]byte, len(res))
    copy(b, res)
    r.Elem = b

    if debug {
        log.Println("-STR: " + string(res))
    }
}

func (r *Reply) parseInt(res []byte) {
    b := make([]byte, len(res))
    copy(b, res)
    r.Elem = b

    if debug {
        log.Println("-INT: " + string(res))
    }
}

func (r *Reply) parseBulk(buf *Reader, res []byte) {
    l, e := strconv.Atoi(string(res))

    if e != nil {
        // TODO: handle error
    }

    if l == -1 {
        if debug {
            log.Println("-BULK: Key does not exist")
        }

        return
    }

    l += 2 // make sure to read \r\n
    data := make([]byte, l)
    n, err := io.ReadFull(buf, data)

    // if we were unable to read all data from socket
    if n != l && err == nil {
        // TODO: handle error
    }

    if err != nil {
        r.Err = err
        return
    }

    l -= 2
    r.Elem = data[:l]

    if debug {
        log.Printf("-BULK: read %d byte, bulk-data %q\n", l, data)
    }
}

func (r *Reply) parseMultiBulk(buf *Reader, res []byte) {
    l, _ := strconv.Atoi(string(res))

    if l == -1  {
        r.Err = errors.New("EOF")
        return
    }

    r.Elems = make([]*Reply, l)

    for i := 0; i < l; i++ {
        rr := Parse(buf)

        if rr.Err != nil {
            r.Err = rr.Err
        }

        r.Elems[i] = rr
    }

    // buffer is reduced to account for `nil` value returns
    r.Elems = r.Elems[:l]

    if debug {
        //log.Printf(": %d == %d %q\n", l, len(r.Elems), r.Elems)
    }
}

func Parse(buf *Reader) *Reply {
    r := new(Reply)
    res, err := buf.ReadSlice(lf)

    if err != nil {
        r.Err = err
        return r
    }

    typ := res[0]
    line := res[1 : len(res)-2]

    switch typ {
    case minus:
        r.parseErr(line)
    case plus:
        r.parseStr(line)
    case colon:
        r.parseInt(line)
    case dollar:
        r.parseBulk(buf, line)
    case star:
        r.parseMultiBulk(buf, line)
    default:
        r.Err = ErrProtocol
    }

    return r
}
