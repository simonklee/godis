# godis

godis - a [Redis](http://redis.io) client for Go. It supports all
commands and features such as transactions and pubsub.

1. [Package docs](http://simonklee.org/pkg/godis/)
2. [Source code](https://github.com/simonz05/godis)

## Install godis

godis is available at github.com. Get it by running.

    $ git clone git://github.com/simonz05/godis.git

And now compile and install it with one command.

    $ make install

Importing godis to your code can be done with `import "godis"`. Thats it!

## Use godis

Checking out the code include a few examples. Here is the code for
the `example/strings.go`.

    package main

    import (
        "godis"
        "fmt"
    )

    func main() {
        // new client on default port 6379, select db 0 and use no password
        c := godis.New("tcp:127.0.0.1:6379", 0, "") 

        // set the key "foo" to "Hello Redis"
        c.Set("foo", "Hello Redis")

        // retrieve the value of "foo". Returns an Elem obj
        elem, _ := c.Get("foo")

        // convert the obj to a string and print it 
        fmt.Println("foo: ", elem.String())
    }

Build the examples. 

    $ make 

To run it we type.

    $ ./string
    foo: Hello Redis

In case your redis server isn't running the output looks like this.

    $ ./string 
    Connection error 127.0.0.1:6379

## Transactions

Pipelines include support for MULTI/EXEC operations.

    c := godis.NewPipeClient("tcp:127.0.0.1:6379", 0, "")

Create a PipeClient. Subsequent commands will be buffered. PipeClient
acts as a regular client, but implements a few extra commands;
`Multi`, `Exec`, `Unwatch`, `Watch`.

    c.Multi()

Calling Multi() wraps subsequent commands inside MULTI .. EXEC.

    c.Set("foo", "bar")
    c.Get("foo")

Commands are still issued as usual, but will return an empty Reply.

    replies := c.Exec()

To execute the buffered commands we call c.Exec(). Exec handles both
MULTI/EXEC pipelines and simply buffered piplines. It returns a slice
of all the *Reply objects for every command we executed.

See `example/transaction.go` for a full example.

## TODO

  * Add tests server commands.

## Acknowledgment

The work on this client started as I was hacking around on Michael Hoisie's
original redis client for Go. Also the recent work done by Frank Müller on his
client gave me some pointers how to better handle return values. 
