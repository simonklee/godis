package main

import (
    "fmt"
    "godis"
    "time"
)

func main() {
    c := godis.New("", 0, "")
    p := c.Pipeline(true)

    c.Set("foo", 1)
    c.Get("foo")

    replies := p.Exec()

    // By calling p.Exec() the following commands were
    // executed.
    //
    //    "MULTI"
    //    "SET" "foo" "1"
    //    "GET" "foo"
    //    "EXEC"
    // 
    // Only the result of the EXEC command is returned
    // in the form of a []*Reply 

    fmt.Println("GET foo:", replies[1].Elem.Int64())

    c.Sync()
    // calling c.Sync() will change the state of the client to 
    // regular non-buffered client again.

    c.Set("foo", 2)
    res, _ := c.Get("foo")
    fmt.Println("GET foo:", res.Int64())
}
