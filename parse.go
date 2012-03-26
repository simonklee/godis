package godis

import (
    "bufio"
    "errors"
    "log"
    "strconv"
)

var (
    debug = false
)

func (r *Reply) parseErr(res []byte) {
    r.Err = errors.New(string(res))

    if debug {
        log.Println("-ERR: " + string(res))
    }
}

func (r *Reply) parseStr(res []byte) {
    r.Elem = res

    if debug {
        log.Println("-STR: " + string(res))
    }
}

func (r *Reply) parseInt(res []byte) {
    r.Elem = res

    if debug {
        log.Println("-INT: " + string(res))
    }
}

func (r *Reply) parseBulk(buf *bufio.Reader, res []byte) {
    l, _ := strconv.Atoi(string(res))

    if l == -1 {
        if debug {
            log.Println("-BULK: Key does not exist")
        }

        r.Err = errors.New("Nonexisting key")
        return
    }

    l += 2 // make sure to read \r\n
    data := make([]byte, l)

    n, err := buf.Read(data)

    // if we were unable to read all data from socket
    if n != l && err == nil {
        more := make([]byte, l-n)

        if _, err := buf.Read(more); err != nil {
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
        log.Printf("-BULK: read %d byte, bulk-data %q\n", l, data)
    }
}

func (r *Reply) parseMultiBulk(buf *bufio.Reader, res []byte) {
    l, _ := strconv.Atoi(string(res))

    if l == -1 {
        r.Err = nil //os.NewError("nothing to read")
        return
    }

    r.Elems = make([]*Reply, l)

    for i := 0; i < l; i++ {
        rr := Parse(buf)

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
        //log.Printf(": %d == %d %q\n", l, len(r.Elems), r.Elems)
    }
}

func Parse(buf *bufio.Reader) *Reply {
    r := new(Reply)
    res, err := buf.ReadBytes(lf)

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
        r.Err = errors.New("Unknown response " + string(typ))
    }

    return r
}
