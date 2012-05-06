package main

import (
    "fmt"
    "github.com/simonz05/godis"
)

func main() {
    c := godis.NewPipeClient("", 0, "")

    c.Multi()
    c.Set("foo", 1)
    c.Get("foo")

    replies := c.Exec()

    // By calling c.Exec() the following commands were
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
}
