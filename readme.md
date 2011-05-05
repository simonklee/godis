# godis

godis - a [Redis](http://redis.io) client for Go.

## Description

The godis package implements a client for Redis. It supports all
redis commands and common features such
as pipelines and pubsub.

  - [Package docs](http://susr.org/godis/pkg/)
  - [Readme](http://susr.org/godis/)

## Install

Use either goinstall or git to make and install the package.

### goinstall

Use [goinstall](http://golang.org/cmd/goinstall/) to download and install the
client with one command.

    $ goinstall github.com/simonz05/godis

goinstall installs godis in the $GOROOT/src/pkg/github.com/simonz05
directory. You can now import godis with `import
"github.com/simonz05/godis"`

### git

The [godis source code](https://github.com/simonz05/godis) is available at
github.com and can be checked out using git.

    $ git clone git://github.com/simonz05/godis.git

To compile it we only need to run make. 

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

## TODO

  * Add tests for sorted set and server commands.
  * Implement transactions.
  * Pipeline need more testing.
  * Multi.
  * Add support for unix socket.

## Acknowledgment

The work on this client started as I was hacking around on Michael Hoisie's
original redis client for Go. Also the recent work done by Frank Müller on his
client gave me some pointers how to better handle return values. 
