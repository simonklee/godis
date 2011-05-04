# godis

godis - a [Redis](http://redis.io) client for Go.

## Description

The godis package implements a client for Redis. It supports all redis
commands and common features such as pipelines and pubsub.

## Install

Use either goinstall or git and make to install the pkg.

### goinstall

[goinstall](http://golang.org/cmd/goinstall/) download and install the
client with one command.

    $ goinstall github.com/simonz05/godis

goinstall puts godis in the $GOROOT/src/pkg/github.com/simonz05
directory. You can now import godis with `import
"github.com/simonz05/godis" 

### git

[godis source code](https://github.com/simonz05/godis) is available at
github.com and can be checked out with git.

    $ git clone git://github.com/simonz05/godis.git

To compile we only need to run make in the godis directory.

    $ make install

You can now import godis with `import "godis"`.

## Examples

To get and run the examples use the git install method explained in
the section above.

    package main

    import (
        "godis"
        "fmt"
    )

    func main() {
        // new client on default port 6379, select db 0 and use no password
        c := godis.New("", 0, "") 

        // set the key "foo" to "Hello Redis"
        c.Set("foo", "Hello Redis")

        // retrieve the value of "foo". Returns an Elem obj
        elem, _ := c.Get("foo")

        // convert the obj to a string and print it 
        fmt.Println("foo: ", elem.String())
    }

To test this example go to the example/ directory.

    $ cd example/
    $ make string

By running make we got an executable called `string`.

    $ ./string
    foo: Hello Redis

If your redis-server isn't running the output looks like this.

    $ ./string 
    Connection error 127.0.0.1:6379

## Docs

This readme file as well as all the [package
docs](http://susr.org/godis/pkg/) for all the commands is available at
is available on [susr.org/godis](http://susr.org/godis/)

## todo

    * Write documentation and add some examples.

    * Add all tests for sorted set and some server stuff.

    * Implement transactions.

    * Pipeline need more testing.

## acknowledgment

The work on this client started as I was hacking around on Michael Hoisie's
original redis client for Go. Also the recent work done by Frank Müller on his
client gave me some pointers how to better handle return values. 
