# exp-godis

godis/exp - an experimental [Redis](http://redis.io) client
for Go. It supports commands and features through a simple
API which aims to be is easy to use.

1. [Package docs](http://go.pkgdoc.org/github.com/simonz05/godis/exp)
2. [Source code](https://github.com/simonz05/godis/exp)

## Install godis

godis/exp is available at github.com. Get it by running.

    $ go get github.com/simonz05/godis/exp

Importing godis/exp will add a `redis` context to program.

## Use exp-godis

A few examples are included. The following demonstrates SET
and GET. See `example/` for more.

    package main

    import (
        "github.com/simonz05/godis/exp"
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

Build and run the example. 

    $ make string; ./string

You should see the following printed in the terminal.

    GET foo: bar 

In case your redis server isn't running, you'll get an
error.

    ERR 127.0.0.1:6379
