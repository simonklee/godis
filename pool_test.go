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
    expected := MaxClientConn*2 + ConnCtr

    if r := Send(c1, "SET", "foo", "foo"); r.Err != nil {
        t.Fatalf("'%s': %s", "SET", r.Err)
    }

    Send(c2, "SET", "bar", "bar")

    start := time.Nanoseconds()

    for i := 0; i < 1000; i++ {
        r1 := Send(c1, "GET", "foo")
        r2 := Send(c2, "GET", "bar")

        if r1.Elem.String() != "foo" && r2.Elem.String() != "bar" {
            t.Error(r1, r2)
            t.FailNow()
        }
    }

    stop := time.Nanoseconds() - start
    t.Logf("time: %.3f\n", float32(stop/1.0e+6)/1000.0)

    if expected != ConnCtr {
        t.Errorf("ConnCtr: expected %d got %d ", expected, ConnCtr)
    }
}
