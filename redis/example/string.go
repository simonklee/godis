package main

import (
    "fmt"
    "github.com/simonz05/godis/redis"
    "os"
)

func main() {
    // new client on default port 6379, select db 0 and use no password
    c := redis.New("", 0, "")

    // set the key "foo" to "Hello Redis"
    if err := c.Set("foo", "Hello Redis"); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // retrieve the value of "foo". Returns an Elem obj
    elem, _ := c.Get("foo")

    // convert the obj to a string and print it 
    fmt.Println("foo:", elem.String())
}
