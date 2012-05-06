package bufin

import (
    "io"
    "testing"
)

type IOReader struct {
    buf []byte
    off int
}

func (b *IOReader) Read(p []byte) (int, error) {
    if len(b.buf[b.off:]) == 0 {
        return 0, io.EOF
    }

    n := copy(p, b.buf[b.off:])
    b.off += n
    return n, nil
}

type readerTest struct {
    iolen int
    rlen  int
    exp   int
    err   error
}

var readerTests = []readerTest{
    {1024, 16, 16, nil},
    {1024, 1024, 1024, nil},
    {1024, 2048, 1024, nil},
    {1024, 0, 0, nil},
    {0, 16, 0, io.EOF},
}

func TestReader(t *testing.T) {
    for _, c := range readerTests {
        p := &IOReader{make([]byte, c.iolen), 0}
        r := NewReader(p)
        buf := make([]byte, c.rlen)

        if n, err := r.Read(buf); c.err != err || c.exp != n {
            t.Errorf("expected `%v` got `%v`, expected %v got %v", c.exp, n, c.err, err)
        }
    }
}

func TestManual(t *testing.T) {
    p := &IOReader{make([]byte, 0), 0}
    r := NewReader(p)
    buf := make([]byte, 16)

    n, e := r.Read(buf)
    println(n, e.Error())
}

func cmpSlice(t *testing.T, a, b []byte) {
    if len(a) != len(b) {
        t.Errorf("expected `%v` got `%v`", len(a), len(b))
    }

    for i := 0; i < len(b); i++ {
        if b[i] != a[i] {
            t.Errorf("expected `%v` got `%v`", a[i], b[i])
        }
    }
}

type sliceTest struct {
    iodata   string
    delim    byte
    head     string
    data     string
    buffered int
    reset    bool
    err      error
}

// tests depend on each other
var sliceTests = []sliceTest{
    {"$3\r\ndata\r\n", '\n', "$3\r\n", "data\r\n", 0, true, nil},
    {"$3\r\ndatamore\r\n", '\n', "$3\r\n", "datamore\r\n", 0, true, nil},
    {"$3\r\ndata\r\n$3\r\ndata\r\n", '\n', "$3\r\n", "data\r\n", 10, false, nil},
    {"", '\n', "$3\r\n", "data\r\n", 0, true, nil},
}

func TestReadSlice(t *testing.T) {
    p := &IOReader{make([]byte, 0), 0}
    r := NewReader(p)

    for _, c := range sliceTests {
        p.buf = append(p.buf, c.iodata...)

        slice, err := r.ReadSlice(c.delim)

        if err != c.err {
            t.Errorf("read expected `%v` got `%v`", c.err, err)
        }

        cmpSlice(t, []byte(c.head), slice)
        data := make([]byte, len(c.data))

        if n, err := r.Copy(data); n != len(c.data) || err != nil {
            t.Errorf("copy expected `%v` got `%v`", len(c.data), n)
        } else {
            cmpSlice(t, data, []byte(c.data))
        }

        if r.Buffered() != c.buffered {
            t.Errorf("buffered expected `%v` got `%v`", c.buffered, r.Buffered())
        }

        if r.Reset() != c.reset {
            t.Errorf("reset expected `%v` got `%v`", c.reset, r.Reset())
        }
    }
}
