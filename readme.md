# godis

godis - a [Redis](http://redis.io) client for Go. It
supports commands and features through a simple API which is
easy to use.

1. [Package docs](http://gopkgdoc.appspot.com/pkg/github.com/simonz05/godis)
2. [Source code](https://github.com/simonz05/godis)

## Install godis

godis is available at github.com. Get it by running.

    $ go get github.com/simonz05/godis

Importing godis to your code can be done with `import "github.com/simonz05/godis"`. Thats it!

## Use godis

Checking out the code include a few examples. Here is the code for
the `example/strings.go`.

    package main

    import (
        "github.com/simonz05/godis"
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

Build the examples. 

    $ make 

To run it we type.

    $ ./string
    GET foo: bar 

In case your redis server isn't running the output looks like this.

    $ ./string 
    Connection error 127.0.0.1:6379
