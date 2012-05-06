package main

import (
    "github.com/simonz05/exp-godis"
)

func main() {
    c := godis.NewClient("tcp:127.0.0.1:6379")

    res, err := c.Call("SET", "foo", "bar")

    if err != nil {
        println(err.Error())
        return
    }

    res, _ = c.Call("GET", "foo")
    println("GET foo:", res.Elem.String())
}
