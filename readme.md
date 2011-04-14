# godis

godis - a [Redis](http://redis.io) client for Go.

* Commands API exactly mimics the Redis commands. 
* Flexible design with support for pipelines, pubsub etc.

## Install

Simply use goinstall to get the client and dependencies.

    $ goinstall github.com/simonz05/godis

### Example

    package main

    import (
        "github.com/simonz05/godis"
        "fmt"
    )

    func main() {
        // new client on default IP/port, redis db to 0 and no password
        c := godis.New("127.0.0.1:6379", 0, "") 

        // set a "foo" to "bar" 
        c.Set("foo", "bar")

        // retrieve the value of "foo"
        foo, _ := c.Get("foo")

        // convert return value back to string and print it
        fmt.Println("foo: ", foo.String())
    }

### Docs

[godis docs](http://simonz05.github.com/godis/) is available on the web.

## todo

* Write documentation and add some examples.

* Add all tests for sorted set and some server stuff.

* Implement transactions.

* Pipeline need more testing.

## acknowledgment

The work on this client started as I was hacking around on Michael Hoisie's
original redis client for Go. Also the recent work done by Frank Müller on his
client gave me some pointers how to better handle return values. 
