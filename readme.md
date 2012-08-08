# godis

Implements a few database clients for Redis.

There is a stable client and an experimental client, `redis`
and `exp`, respectively. To use any of them simply add.

    import "github.com/simonz05/godis/redis"

or 

    import "github.com/simonz05/godis/exp"

Both versions provide a `redis` package which is used to
create a client and talk to the database. For a quick start
check out either projects readme and example. Package
reference is also available.

1. [godis/redis](http://go.pkgdoc.org/github.com/simonz05/godis/redis)
2. [godis/exp](http://go.pkgdoc.org/github.com/simonz05/godis/exp)

**HINT**

If you installed godis with the go tool

    go get github.com/hoisie/web
    
you currently need to use

    import "github.com/simonz05/godis"
    
in order to use godis since the go tool does not fetch the current version,
which has undergone some stuctural changes.