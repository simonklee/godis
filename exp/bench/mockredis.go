package main

import (
    "fmt"
    "net"
    "os"
    "time"
)

func MockRedis(ln net.Listener) {
    cnt := 0

    for {
        conn, err := ln.Accept()
        cnt++

        if err != nil {
            //fmt.Fprintln(os.Stderr, "accept err", err.Error())
            return
        }

        go handle(conn, cnt)
    }
}

func handle(c net.Conn, nr int) {
    buf := make([]byte, 16)

    for {
        start := time.Now()
        _, err := c.Read(buf)
        fmt.Printf("%.6fs\n", time.Now().Sub(start).Seconds())

        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            return
        }

        _, err = c.Write([]byte("$3\r\nbar\r\n"))

        if err != nil {
            fmt.Fprintf(os.Stderr, "write err: %s\n", err.Error())
            return
        }

        //s := string(buf)
        //fmt.Fprintf(os.Stdout, "#%d nread: %d, nwrite: %d\n", nr, nread, nwrite)
    }
}
