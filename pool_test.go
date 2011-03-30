package godis

import (
    "testing"
    "net"
    "time"
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
    c1 := New("", 0, "")
    c2 := New("", 0, "")
    expected := MaxClientConn * 2 + ConnCtr

    c1.Set("foo", "foo")
    c2.Set("bar", "bar")

    start := time.Nanoseconds()

    for i := 0; i < 1000; i++ {
        s1, _ := c1.Get("foo")
        s2, _ := c2.Get("foo")

        if string(s1) == "foo" && string(s2) == "bar" {
        }
    }

    stop := time.Nanoseconds() - start
    t.Log("time: %.2f", float32(stop / 1.0e+6) / 1000.0)

    if expected != ConnCtr {
        t.Errorf("ConnCtr: expected %d got %d ", expected, ConnCtr)
    }
}
