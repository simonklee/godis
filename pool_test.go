package godis

import (
    "testing"
    "net"
    "time"
    "log"
)

func getConn(t *testing.T) (conn *net.TCPConn) {
    var defaultAddr string = "127.0.0.1:6379"

    addr, err := net.ResolveTCPAddr(defaultAddr)
    if err != nil {
        t.Errorf("ResolveAddr error for " + defaultAddr)
    }

    conn, err = net.DialTCP("tcp", nil, addr)
    if err != nil {
        t.Errorf("Connection error " + addr.String())
    }
    return
}

func TestPool(t *testing.T) {
    p := NewPool()

    for i := 0; i < 10; i++ {
        c := p.Pop()
        if c == nil {
            c = getConn(t)
        }

        go func(c *net.TCPConn) {
            p.Push(c)
        }(c)
    }
}

func TestPoolSize(t *testing.T) {
    c := New("", 0, "")

    c.Send("SET", "key", "foo")
    start := time.Nanoseconds()
    for i := 0; i < 1000; i++ {
        in, _ := c.Send("GET", "key")
        s, _ := in.([]byte)
        if string(s) == "foo" {
        }
    }
    stop := time.Nanoseconds() - start

    log.Printf("time: %.2f", float32(stop / 1.0e+6) / 1000.0)
    in, _ := c.Send("GET", "key")
    s, _ := in.([]byte)
    l(string(s))

    if MaxClientConn * 2 != ConnCtr {
        t.Errorf("ConnCtr: expected %d got %d ", MaxClientConn * 2, ConnCtr)
    }

    log.Printf("%f", 1e+6)
}
