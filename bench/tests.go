package main

import (
    "github.com/simonz05/exp-godis"
    "net"
    "os"
    "fmt"
)

func init() {
    tests["set"] = setHandle
    tests["setpipe"] = setPipelineHandle
    tests["get"] = getHandle
    tests["rpush"] = rpushHandle
    tests["calla"] = callaHandle
    tests["callb"] = callbHandle
    tests["mock"] = mockHandle
}

func rpushHandle(c *godis.Client, ch chan bool) {
    for _ = range ch {
        c.Call("RPUSH", "foo", "bar")
    }
}

func setHandle(c *godis.Client, ch chan bool) {
    for _ = range ch {
        c.Call("SET", "foo", "bar")
    }
}

func getHandle(c *godis.Client, ch chan bool) {
    for _ = range ch {
        c.Call("GET", "foo")
    }
}

func setPipelineHandle(c *godis.Client, ch chan bool) {
    p := c.Pipeline()
    send := 0

    for _ = range ch {
        p.Call("SET", "foo", "bar")
        send++

        if send == *P {
            for i := 0; i < *P; i++ {
                p.Read()
            }
            send = 0
        }
    }
}

func mockHandle(c *godis.Client, ch chan bool) {
    conn, err := net.Dial("tcp", "127.0.0.1:6381")

    if err != nil {
        fmt.Fprintln(os.Stderr, "dial error", err.Error())
        os.Exit(1)
    }

    cmd := []byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")
    buf := make([]byte, 16)//len([]byte("$3\r\nbar\r\n")))

    for _ = range ch {
        if _, err := conn.Write(cmd); err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
        }

        if _, err := conn.Read(buf); err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
        }
    }
}

func callaHandle(c *godis.Client, ch chan bool) {
    buf := make([]byte, 1024*16)
    var conn *godis.Conn 
    get := []byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")
    conn, _ = c.CallA("GET", "foo")

    for _ = range ch {
        conn.Conn.Write(get)
        conn.Conn.Read(buf)
    }
}

func callbHandle(c *godis.Client, ch chan bool) {
    buf := make([]byte, 1024*4)
    var conn *godis.Conn 
    get := []byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")
    conn, _ = c.CallA("GET", "foo")

    if tcp, ok := conn.Conn.(*net.IPConn); ok {
        tcp.SetWriteBuffer(16)
        tcp.SetReadBuffer(16)
    }

    for _ = range ch {
        conn.Conn.Write(get)
        conn.Conn.Read(buf)
    }
    c.CallADone(conn)
}
