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

    c1.Send("SET", "foo", "foo")
    c2.Send("SET", "bar", "bar")

    start := time.Nanoseconds()

    for i := 0; i < 1000; i++ {
        r1 := c1.Send("GET", "foo")
        r2 := c2.Send("GET", "bar")
        
        if r1.Elem.String() != "foo" && r2.Elem.String() != "bar" {
            t.Error(r1, r2)
            t.FailNow()
        }
    }

    stop := time.Nanoseconds() - start
    t.Logf("time: %.3f\n", float32(stop / 1.0e+6) / 1000.0)

    if expected != ConnCtr {
        t.Errorf("ConnCtr: expected %d got %d ", expected, ConnCtr)
    }
}
