package main

import (
    "fmt"
    "godis"
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

    fmt.Println("GET foo:", replies[1].Elem.String())
}
