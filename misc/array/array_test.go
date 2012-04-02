package main

import (
    "testing"
    "io"
)

type IOReader struct {
    buf []byte
}

func (b *IOReader) Read(p []byte) (int, error) {
    if len(b.buf) == 0 {
        return 0, io.EOF
    }

    return copy(p, b.buf), nil
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
        p := &IOReader{make([]byte, c.iolen)}
        r := NewReader(p)
        buf := make([]byte, c.rlen)

        if n, err := r.Read(buf); c.err != err || c.exp != n {
            t.Errorf("expected `%v` got `%v`, expected %v got %v", c.exp, n, c.err, err)
        }
    }
}

func TestManual(t *testing.T) {
    p := &IOReader{make([]byte, 0)}
    r := NewReader(p)
    buf := make([]byte, 16)

    n, e := r.Read(buf)
    println(n, e.Error())
}
