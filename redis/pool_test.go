package redis

import (
    "testing"
    "time"
)

func getConn(t *testing.T) *conn {
    c, err := newConn("tcp", "127.0.0.1:6379", 0, "")

    if err != nil {
        t.Errorf("err " + err.Error())
    }

    return c
}

func TestPoolSimple(t *testing.T) {
    p := newPool()

    for i := 0; i < 10; i++ {
        c := p.pop()
        if c == nil {
            c = getConn(t)
        }

        go func(c *conn) {
            p.push(c)
        }(c)
    }
}

func TestPoolSize(t *testing.T) {
    c1 := New("", 0, "")
    c2 := New("", 0, "")
    expected := MaxClientConn*2 + connCount

    if r := SendStr(c1.Rw, "SET", "foo", "foo"); r.Err != nil {
        t.Fatalf("'%s': %s", "SET", r.Err)
    }

    SendStr(c2.Rw, "SET", "bar", "bar")

    start := time.Now()

    for i := 0; i < 1000; i++ {
        r1 := SendStr(c1.Rw, "GET", "foo")
        r2 := SendStr(c2.Rw, "GET", "bar")

        if r1.Elem.String() != "foo" && r2.Elem.String() != "bar" {
            t.Error(r1, r2)
            t.FailNow()
        }
    }

    stop := time.Now().Sub(start)
    t.Logf("time: %.3f\n", float32(stop/1.0e+6)/1000.0)

    if expected != connCount {
        t.Errorf("connCount: expected %d got %d ", expected, connCount)
    }
}
