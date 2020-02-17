package main

import (
    "fmt"
    "insmo.com/godis/exp"
    "net"
    "os"
)

func init() {
    tests["set"] = setHandle
    tests["setpipe"] = setPipelineHandle
    tests["get"] = getHandle
    tests["rpush"] = rpushHandle
    tests["mock"] = mockHandle
}

func rpushHandle(c *redis.Client, ch chan bool) {
    for _ = range ch {
        c.Call("RPUSH", "foo", "bar")
    }
}

func setHandle(c *redis.Client, ch chan bool) {
    for _ = range ch {
        c.Call("SET", "foo", "bar")
    }
}

func getHandle(c *redis.Client, ch chan bool) {
    for _ = range ch {
        c.Call("GET", "0")
    }
}

func setPipelineHandle(c *redis.Client, ch chan bool) {
    p := c.AsyncClient()
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

func mockHandle(c *redis.Client, ch chan bool) {
    conn, err := net.Dial("tcp", "127.0.0.1:6381")

    if err != nil {
        fmt.Fprintln(os.Stderr, "dial error", err.Error())
        os.Exit(1)
    }

    cmd := []byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")
    buf := make([]byte, 16) //len([]byte("$3\r\nbar\r\n")))

    for _ = range ch {
        if _, err := conn.Write(cmd); err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
        }

        if _, err := conn.Read(buf); err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
        }
    }
}
