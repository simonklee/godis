package schema

import (
    "testing"

    "github.com/simonz05/godis/exp"
)

var db *redis.Client

func init() {
    db = redis.NewClient("tcp:localhost:6379", 9, "")
}

func TestPut(t *testing.T) {
    _, e := db.Call("FLUSHDB")
    if e != nil {
        t.Fatalf(e.Error())
    }
}
