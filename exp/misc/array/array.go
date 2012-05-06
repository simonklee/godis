package main

import (
    "fmt"
    "io"
    "log"
)

const (
    MAX_IOBUFLEN = uint16(1024)
    MIN_IOBUFLEN = uint16(8)
)

type Reader struct {
    data [MAX_IOBUFLEN]byte
    buf  []byte
    rd   io.Reader
    r, w int
}

func NewReader(rd io.Reader) (r *Reader) {
    r = new(Reader)
    r.buf = r.data[:MAX_IOBUFLEN/2]
    r.rd = rd
    return r
}

func (b *Reader) String() string {
    return fmt.Sprintf("len: %d, cap: %d, read: %d, width: %d", len(b.buf), cap(b.buf), b.r, b.w)
}

// will overwrite any existing data in the buffer
func (b *Reader) fill() error {
    b.r = 0
    n, e := b.rd.Read(b.buf)

    if e != nil {
        return e
    }

    b.w = n
    return nil
}

func (b *Reader) AdjustBuflen(n uint16) {
    switch {
    case n > MAX_IOBUFLEN:
        n = MAX_IOBUFLEN
    case n < MIN_IOBUFLEN:
        n = MIN_IOBUFLEN
    default:
        n--
        n |= n >> 1
        n |= n >> 2
        n |= n >> 4
        n |= n >> 8
        n++
    }

    b.buf = b.data[:n]
}

func (b *Reader) Buffered() int {
    return b.w - b.r
}

func (b *Reader) ReadFull(p []byte) (int, error) {
    return b.ReadAtLeast(p, len(p))
}

func (b *Reader) ReadAtLeast(p []byte, min int) (n int, e error) {
    for n < min && e == nil {
        var nn int
        nn, e = b.Read(p[n:])
        n += nn
    }

    if e == io.EOF {
        if n >= min {
            e = nil
        }
    }

    return
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
        if n >= len(b.buf) {
            log.Println("Read directly from IO")
            return b.rd.Read(p)
        }

        if e = b.fill(); e != nil {
            return 0, e
        }
    }

    // drain buffer
    if n > b.w-b.r {
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

    if n > b.w-b.r {
        n = b.w - b.r
    }

    copy(p[0:n], b.buf[b.r:])
    b.r += n
    return n, nil
}

func main() {
    r := NewReader(nil)
    println(r.String())
}
