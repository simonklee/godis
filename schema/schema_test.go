package schema

import (
    "testing"

    "insmo.com/godis/exp"
)

var db *redis.Client

func init() {
    redis.MaxConnections = 1
    db = redis.NewClient("tcp:localhost:6379", 9, "")
}

type S1 struct {
    Id int64 `redis:"id"`
}

type User struct {
    Id       int64  `redis:"id"`
    Username string `redis:"username,unique,index"`
    Email    string `redis:"email,unique"`
}

type putTest struct {
    s   interface{}
    k   *Key
    e   error
}

var putTests = [][]putTest{
    {{&S1{1}, NewKey("s1", 1), nil}},
    {{S1{1}, NewKey("s1", 1), TypeError}},
    {{&S1{1}, NewKey("s1", 1), nil}, {&S1{2}, NewKey("s1", 2), nil}},
    {{&S1{1}, NewKey("s1", 1), nil}, {&S1{1}, NewKey("s1", 1), nil}},
    {{&User{1, "foo", "foo@foo.com"}, NewKey("user", 1), nil}},
    {
        {&User{1, "foo", "foo@foo.com"}, NewKey("user", 1), nil},
        {&User{1, "foo", "foo@foo.com"}, NewKey("user", 1), nil},
    },
    {
        {&User{1, "foo", "foo@foo.com"}, NewKey("user", 1), nil},
        {&User{2, "foo", "foo@foo.com"}, NewKey("user", 2), newUniqueError("username", "foo")},
    },
    {
        {&User{1, "foo", "foo@foo.com"}, NewKey("user", 1), nil},
        {&User{2, "bar", "foo@foo.com"}, NewKey("user", 2), newUniqueError("email", "foo@foo.com")},
    },
    {
        {&User{1, "foo", "foo@foo.com"}, NewKey("user", 1), nil},
        {&User{2, "bar", "bar@foo.com"}, NewKey("user", 2), nil},
    },
    {
        {&User{0, "foo", "foo@foo.com"}, NewKey("user", 0), nil},
        {&User{0, "foo", "foo@foo.com"}, NewKey("user", 0), newUniqueError("username", "foo")},
    },
}

func TestPut(t *testing.T) {
    for _, tests := range putTests {
        _, e := db.Call("FLUSHDB")

        if e != nil {
            t.Fatalf(e.Error())
        }

        for _, o := range tests {
            _, e := Put(db, o.k, o.s)

            if e != o.e {
                if e == nil {
                    t.Fatalf("expected `%s`, got e: nil", o.e.Error())
                } else {
                    t.Fatal(e.Error())
                }
            }
        }
    }
}
