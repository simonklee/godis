package main

import (
    "github.com/simonz05/godis"
    "fmt"
)

func main() {
    // new client on default IP/port, redis db to 0 and no password
    c := godis.New("127.0.0.1:6379", 0, "")

    // set a "foo" to "bar" 
    godis.Set(c, "foo", "bar")

    // retrieve the value of "foo"
    foo, _ := godis.Get(c, "foo")

    // convert return value back to string and print it
    fmt.Println("foo: ", foo.String())
}
