// A fixed size buffer implenation of the bufio library
// not all methods are ported
package godis

import (
    "fmt"
    "io"
    "log"
    "bytes"
    "errors"
)

const IOBUFLEN = 1024

var (
    ErrFullBuf = errors.New("Full buffer")
    ErrNotFound = errors.New("Not found")
)

type Reader struct {
    data        [IOBUFLEN*8]byte
    buf         []byte
    rd          io.Reader
    r, w        int
    reads, movs int64
}

func NewReader(rd io.Reader) (r *Reader) {
    r = new(Reader)
    r.buf = r.data[:]
    r.rd = rd 
    return r
}

func (b *Reader) String() string {
    return fmt.Sprintf("len: %d, cap: %d, read: %d, width: %d, buffered: %d, sycall: %d, move: %d", len(b.buf), cap(b.buf), b.r, b.w, b.Buffered(), b.reads, b.movs)
}

// reset to recover space if buf is empty
func (b *Reader) Reset() bool {
    if b.w == b.r {
        b.w = 0
        b.r = 0
        return true
    }

    return false
}

func (b *Reader) fill() error {
    b.Reset()

    if b.r > 0 {
        // move existing data to beginning of buffer
        //println("move")
        copy(b.buf, b.buf[b.r:b.w])
        b.w -= b.r
        b.r = 0

        // statistics
        b.movs++
    }

    if len(b.buf[b.w:]) < IOBUFLEN {
        return ErrFullBuf
    }

    slice := b.buf[b.w:IOBUFLEN+b.w]
    n, e := b.rd.Read(slice)
    b.w += n

    // statistics
    b.reads++

    if e != nil {
        return e
    }

    return nil
}

func (b *Reader) Buffered() int {
    return b.w - b.r
}

func (b *Reader) Incr(n int) int {
    if n > b.Buffered() {
        return 0
    }

    b.r += n
    return n
}

// either reads from the static buffer or if len(p) > len(buf), 
// read len(p) bytes from socket directly into p
func (b *Reader) Read(p []byte) (n int, e error) {
    n = len(p)

    if n == 0 {
        return 0, nil
    }

    if b.w == b.r {
        // read request is larger then current window size
        if n >= len(b.buf[b.w:IOBUFLEN]) {
            //log.Println("Read directly from IO")
            n, e = b.rd.Read(p)
            b.reads++
            return n, e
        }

        log.Println("End of buffer")
        if e = b.fill(); e != nil {
            return 0, e
        }
    }

    // drain buffer
    if n > b.w - b.r {
        n = b.w - b.r
    }

    copy(p[0:n], b.buf[b.r:])
    b.r += n
    return n, nil
}

// copies len(p) bytes from r.buf[r:] to p
// if len(p) > r.buf[r:]
func (b *Reader) Copy(p []byte) (n int, e error) {
    n = len(p)

    if b.w == b.r || n == 0 {
        return 0, nil
    }

    if n > b.w - b.r {
        n = b.w - b.r
    }

    copy(p[0:n], b.buf[b.r:])
    b.r += n
    return n, nil
}

func (b *Reader) IndexSlice(delim byte) (line []byte, err error) {
    if i := bytes.IndexByte(b.buf[b.r:b.w], delim); i >= 0 {
        line = b.buf[b.r:b.r+i+1]
        b.r += i + 1

        return line, nil
    }

    return nil, ErrNotFound
}

func (b *Reader) ReadSlice(delim byte) (line []byte, err error) {
    for {
        off := b.r
        i := bytes.IndexByte(b.buf[off:b.w], delim)

        if i >= 0 {
            line = b.buf[b.r:off+i+1]
            b.r = off+i+1
            return line, nil
        }

        err = b.fill()

        if err != nil {
            line = b.buf[b.r:b.w]
            b.r = b.w
            return line, err
        }
    }

    panic("never")
}
