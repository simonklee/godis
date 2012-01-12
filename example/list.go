package main

import (
    "fmt"
    "godis"
    "os"
)

func main() {
    // new client on default port 6379, select db 0 and use no password.
    c := godis.New("", 0, "")

    // values we want to store
    values := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}

    // push the values to the redis list 
    for _, v := range values {
        if _, err := c.Rpush("bar", v); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
    }

    // retrieve the items in the list. Returns a Reply object.
    res, _ := c.Lrange("bar", 0, 9)

    // convert the list to an array of ints and print it.
    fmt.Println("bar:", res.IntArray())
}
