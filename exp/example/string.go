package main

import (
    "insmo.com/godis/exp"
)

func main() {
    c := redis.NewClient("tcp:127.0.0.1:6379")

    res, err := c.Call("SET", "foo", "bar")

    if err != nil {
        println(err.Error())
        return
    }

    res, _ = c.Call("GET", "foo")
    println("GET foo:", res.Elem.String())
}
