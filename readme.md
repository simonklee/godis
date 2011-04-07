# godis

A simple client for [Redis](http://redis.io).

* Commands API exactly mimics the Redis commands. It is extremly consistent with
  the real deal.
* Flexible design to with the goal of bringing support for pipelining to the client.

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
        godis.Set(c, "foo", "bar")

        // retrieve the value of "foo"
        foo, _ := godis.Get(c, "foo")

        // convert return value back to string and print it
        fmt.Println("foo: ", foo.String())
    }

## todo

* Write documentation and add some examples.

* Add all tests for sorted set and some server stuff.

* Implement pub-sub and transactions.

* PipeClient logic is not safe at all. If an error occurs and the user
continues to call ReadReply after that, a new connection will be poped and we
will try to continue reading, making the client hang/timeout ... repeat.

## acknowledgment

The work on this client started as I was hacking around on Michael Hoisie's
original redis client for Go. Also the recent work done by Frank Müller on his
client gave me some pointers how to better handle return values. 
