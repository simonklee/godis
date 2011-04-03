package godis

import (
    "os"
    "strconv"
    "testing"
    "reflect"
    "time"
    "fmt"
)

func error(t *testing.T, name string, expected, got interface{}, err os.Error) {
    t.Errorf("`%s` expected `%v` got `%v`, err(%v)", name, expected, got, err)
}

func TestGeneric(t *testing.T) {
    c := New("", 0, "")
    Send(c, "FLUSHDB")

    if res, err := c.Randomkey(); res != "" {
        error(t, "randomkey", "", res, err)
    }

    c.Set("foo", "foo")

    if res, err := c.Randomkey(); res != "foo" {
        error(t, "randomkey", "foo", res, err)
    }

    ex, _ := c.Exists("key")
    nr, _ := c.Del("key")

    if (ex && nr != 1) || (!ex && nr != 0) {
        error(t, "del", "unknown", nr, nil)
    }

    c.Set("foo", "foo")
    c.Set("bar", "bar")
    c.Set("baz", "baz")

    if nr, err := c.Del("foo", "bar", "baz"); nr != 3 {
        error(t, "del", 3, nr, err)
    }

    c.Set("foo", "foo")

    if res, err := c.Expire("foo", 10); !res {
        error(t, "expire", true, res, err)
    }
    if res, err := c.Persist("foo"); !res {
        error(t, "persist", true, res, err)
    }
    if res, err := c.Ttl("foo"); res == 0 {
        error(t, "ttl", 0, res, err)
    }
    if res, err := c.Expireat("foo", time.Seconds()+10); !res {
        error(t, "expireat", true, res, err)
    }
    if res, err := c.Ttl("foo"); res <= 0 {
        error(t, "ttl", "> 0", res, err)
    }
    if err := c.Rename("foo", "bar"); err != nil {
        error(t, "rename", nil, nil, err)
    }
    if err := c.Rename("foo", "bar"); err == nil {
        error(t, "rename", "error", nil, err)
    }
    if res, err := c.Renamenx("bar", "foo"); !res {
        error(t, "renamenx", true, res, err)
    }

    c.Set("bar", "bar")

    if res, err := c.Renamenx("foo", "bar"); res {
        error(t, "renamenx", false, res, err)
    }

    c2 := New("", 1, "")
    c2.Del("foo")
    if res, err := c.Move("foo", 1); res != true {
        error(t, "move", true, res, err)
    }
}

func TestKeys(t *testing.T) {
    c := New("", 0, "")
    Send(c, "FLUSHDB")
    Send(c, "MSET", "foo", "one", "bar", "two", "baz", "three")

    res, err := c.Keys("foo")

    if err != nil {
        error(t, "keys", nil, nil, err)
    }

    expected := []string{"foo"}

    if len(res) != len(expected) {
        error(t, "keys", len(res), len(expected), nil)
    }

    for i, v := range res {
        if v != expected[i] {
            error(t, "keys", expected[i], v, nil)
        }
    }
}

func TestSort(t *testing.T) {
    c := New("", 0, "")
    Send(c, "FLUSHDB")
    Send(c, "RPUSH", "foo", "2")
    Send(c, "RPUSH", "foo", "3")
    Send(c, "RPUSH", "foo", "1")

    res, err := c.Sort("foo")

    if err != nil {
        error(t, "sort", nil, nil, err)
    }

    expected := []int{1, 2, 3}
    if len(res.Elems) != len(expected) {
        error(t, "sort", len(res.Elems), len(expected), nil)
    }

    for i, v := range res.Elems {
        r := int(v.Elem.Int64())
        if r != expected[i] {
            error(t, "sort", expected[i], v, nil)
        }
    }
}

func TestString(t *testing.T) {
    c := New("", 0, "")
    Send(c, "FLUSHDB")

    if res, err := c.Decr("qux"); err != nil || res != -1 {
        error(t, "decr", -1, res, err)
    }

    if res, err := c.Decrby("qux", 1); err != nil || res != -2 {
        error(t, "decrby", -2, res, err)
    }

    if res, err := c.Incrby("qux", 1); err != nil || res != -1 {
        error(t, "incrby", -1, res, err)
    }

    if res, err := c.Incr("qux"); err != nil || res != 0 {
        error(t, "incrby", 0, res, err)
    }

    if res, err := c.Setbit("qux", 0, 1); err != nil || res != 0 {
        error(t, "setbit", 0, res, err)
    }

    if res, err := c.Getbit("qux", 0); err != nil || res != 1 {
        error(t, "getbit", 1, res, err)
    }

    if err := c.Set("foo", "foo"); err != nil {
        t.Errorf(err.String())
    }

    if res, err := c.Append("foo", "bar"); err != nil || res != 6 {
        error(t, "append", 6, res, err)
    }

    if res, err := c.Get("foo"); err != nil || res != "foobar" {
        error(t, "get", "foobar", res, err)
    }

    if _, err := c.Get("foobar"); err == nil {
        error(t, "get", "error", nil, err)
    }

    if res, err := c.Getrange("foo", 0, 2); err != nil || res != "foo" {
        error(t, "getrange", "foo", res, err)
    }

    if res, err := c.Setrange("foo", 0, "qux"); err != nil || res != 6 {
        error(t, "setrange", 6, res, err)
    }

    if res, err := c.Getset("foo", "foo"); err != nil || res != "quxbar" {
        error(t, "getset", "quxbar", res, err)
    }

    if res, err := c.Setnx("foo", "bar"); err != nil || res != false {
        error(t, "setnx", false, res, err)
    }

    if res, err := c.Strlen("foo"); err != nil || res != 3 {
        error(t, "strlen", 3, res, err)
    }

    if err := c.Setex("foo", 10, "bar"); err != nil {
        error(t, "setex", nil, nil, err)
    }

    out := []string{"foo", "bar", "qux"}
    in := map[string]string{"foo": "foo", "bar": "bar", "qux": "qux"}

    if err := c.Mset(in); err != nil {
        error(t, "mset", nil, nil, err)
    }

    if res, err := c.Msetnx(in); err != nil || res == true {
        error(t, "msetnx", false, res, err)
    }

    res, err := c.Mget(append([]string{"il"}, out...)...)

    if err != nil || len(res) != 3 {
        error(t, "mget", 3, len(res), err)
        t.FailNow()
    }

    for i, v := range res {
        if v != out[i] {
            error(t, "mget", out[i], v, nil)
        }
    }
}

func TestList(t *testing.T) {
    c := New("", 0, "")
    Send(c, "FLUSHDB")

    if res, err := c.Lpush("foobar", "foo"); err != nil || res != 1 {
        error(t, "LPUSH", 1, res, err)
    }

    if res, err := c.Linsert("foobar", "AFTER", "foo", "bar"); err != nil || res != 2 {
        error(t, "Linsert", 2, res, err)
    }

    if res, err := c.Linsert("foobar", "AFTER", "qux", "bar"); err != nil || res != -1 {
        error(t, "Linsert", -1, res, err)
    }

    if res, err := c.Llen("foobar"); err != nil || res != 2 {
        error(t, "Llen", 2, res, err)
    }

    if res, err := c.Lindex("foobar", 0); err != nil || res.String() != "foo" {
        error(t, "Lindex", "foo", res, err)
    }

    if res, err := c.Lpush("foobar", "qux"); err != nil || res != 3 {
        error(t, "Lpush", 3, res, err)
    }

    if res, err := c.Lpop("foobar"); err != nil || string(res) != "qux" {
        error(t, "Lpop", "qux", res, err)
    }

    want1 := []*Reply{
        &Reply{Elem: []byte("foo")},
        &Reply{Elem: []byte("bar")},
    }

    if out, err := c.Lrange("foobar", 0, 1); err != nil || !reflect.DeepEqual(want1, out.Elems) {
        error(t, "Lrange", nil, nil, err)
    }

    want1 = []*Reply{}
    for i := 0; i < 600; i++ {
        want1 = append(want1, &Reply{Elem: []byte(strconv.Itoa(i))})
        c.Rpush("foobaz", i)
        j := 0

        if i > 50 {
            j = i - 50
        }

        if res, err := c.Lrange("foobaz", j, i); err != nil || !reflect.DeepEqual(want1[j:], res.Elems) {
            error(t, "Lranges", nil, res, err)
            t.FailNow()
        }
    }

    want := []string{"foo"}

    if res, err := c.Lrem("foobar", 0, "bar"); err != nil || res != 1 {
        error(t, "Lrem", 1, res, err)
    }

    want = []string{"bar"}

    if err := c.Lset("foobar", 0, "bar"); err != nil {
        error(t, "Lrem", nil, nil, err)
    }

    want = []string{}

    if err := c.Ltrim("foobar", 1, 0); err != nil {
        error(t, "Ltrim", nil, nil, err)
    }

    want = []string{"foo", "bar", "qux"}
    var res int64
    var err os.Error

    for _, v := range want {
        res, err = c.Rpush("foobar", v)
    }

    if err != nil || res != 3 {
        error(t, "Rpush", 3, res, err)
    }

    if res, err := c.Rpushx("foobar", "baz"); err != nil || res != 4 {
        error(t, "Rpushx", 4, res, err)
    }

    if res, err := c.Rpop("foobar"); err != nil || res.String() != "baz" {
        error(t, "Rpop", "baz", res, err)
    }

    if res, err := c.Rpoplpush("foobar", "foobaz"); err != nil || res.String() != "qux" {
        error(t, "Rpop", "qux", res, err)
    }
}

func TestHash(t *testing.T) {
    c := New("", 0, "")
    Send(c, "FLUSHDB")

    if res, err := c.Hset("foobar", "foo", "foo"); err != nil || res != true {
        error(t, "Hset", true, res, err)
    }

    if res, err := c.Hset("foobar", "foo", "foo"); err != nil || res != false {
        error(t, "Hset", false, res, err)
    }

    if res, err := c.Hget("foobar", "foo"); err != nil || res.String() != "foo" {
        error(t, "Hget", "foo", res, err)
    }

    if res, err := c.Hdel("foobar", "foo"); err != nil || res != true {
        error(t, "Hdel", true, res, err)
    }

    if res, err := c.Hexists("foobar", "foo"); err != nil || res != false {
        error(t, "Hexists", false, res, err)
    }

    c.Hset("foobar", "foo", 1)
    c.Hset("foobar", "bar", 2)
    want := []*Reply{
        &Reply{Elem: []byte("foo")},
        &Reply{Elem: []byte("1")},
        &Reply{Elem: []byte("bar")},
        &Reply{Elem: []byte("2")},
    }

    if res, err := c.Hgetall("foobar"); err != nil || !reflect.DeepEqual(want, res.Elems) {
        error(t, "Hexists", want, res, err)
    }
}

func BenchmarkRpush(b *testing.B) {
    c := New("", 0, "")
    start := time.Nanoseconds()
    for i := 0; i < b.N; i++ {
        c.Rpush("qux", "qux")
    }
    c.Del("qux")
    stop := time.Nanoseconds() - start
    fmt.Fprintf(os.Stdout, "time: %.2f\n", float32(stop/1.0e+6)/1000.0)
}
