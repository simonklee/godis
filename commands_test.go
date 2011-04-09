package godis

import (
    "os"
    "testing"
    "reflect"
    "time"
    "log"
)

func TestGeneric(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := Randomkey(c); res != "" {
        error(t, "randomkey", "", res, err)
    }

    Set(c, "foo", "foo")

    if res, err := Randomkey(c); res != "foo" {
        error(t, "randomkey", "foo", res, err)
    }

    ex, _ := Exists(c, "key")
    nr, _ := Del(c, "key")

    if (ex && nr != 1) || (!ex && nr != 0) {
        error(t, "del", "unknown", nr, nil)
    }

    Set(c, "foo", "foo")
    Set(c, "bar", "bar")
    Set(c, "baz", "baz")

    if nr, err := Del(c, "foo", "bar", "baz"); nr != 3 {
        error(t, "del", 3, nr, err)
    }

    Set(c, "foo", "foo")

    if res, err := Expire(c, "foo", 10); !res {
        error(t, "expire", true, res, err)
    }
    if res, err := Persist(c, "foo"); !res {
        error(t, "persist", true, res, err)
    }
    if res, err := Ttl(c, "foo"); res == 0 {
        error(t, "ttl", 0, res, err)
    }
    if res, err := Expireat(c, "foo", time.Seconds()+10); !res {
        error(t, "expireat", true, res, err)
    }
    if res, err := Ttl(c, "foo"); res <= 0 {
        error(t, "ttl", "> 0", res, err)
    }
    if err := Rename(c, "foo", "bar"); err != nil {
        error(t, "rename", nil, nil, err)
    }
    if err := Rename(c, "foo", "bar"); err == nil {
        error(t, "rename", "error", nil, err)
    }
    if res, err := Renamenx(c, "bar", "foo"); !res {
        error(t, "renamenx", true, res, err)
    }

    Set(c, "bar", "bar")

    if res, err := Renamenx(c, "foo", "bar"); res {
        error(t, "renamenx", false, res, err)
    }

    c2 := New("", 1, "")
    Del(c2, "foo")
    if res, err := Move(c, "foo", 1); res != true {
        error(t, "move", true, res, err)
    }
}

func TestKeys(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }
    SendStr(c, "MSET", "foo", "one", "bar", "two", "baz", "three")

    res, err := Keys(c, "foo")

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
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }
    SendStr(c, "RPUSH", "foo", "2")
    SendStr(c, "RPUSH", "foo", "3")
    SendStr(c, "RPUSH", "foo", "1")

    res, err := Sort(c, "foo")

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
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := Decr(c, "qux"); err != nil || res != -1 {
        error(t, "decr", -1, res, err)
    }

    if res, err := Decrby(c, "qux", 1); err != nil || res != -2 {
        error(t, "decrby", -2, res, err)
    }

    if res, err := Incrby(c, "qux", 1); err != nil || res != -1 {
        error(t, "incrby", -1, res, err)
    }

    if res, err := Incr(c, "qux"); err != nil || res != 0 {
        error(t, "incrby", 0, res, err)
    }

    if res, err := Setbit(c, "qux", 0, 1); err != nil || res != 0 {
        error(t, "setbit", 0, res, err)
    }

    if res, err := Getbit(c, "qux", 0); err != nil || res != 1 {
        error(t, "getbit", 1, res, err)
    }

    if err := Set(c, "foo", "foo"); err != nil {
        t.Errorf(err.String())
    }

    if res, err := Append(c, "foo", "bar"); err != nil || res != 6 {
        error(t, "append", 6, res, err)
    }

    if res, err := Get(c, "foo"); err != nil || res.String() != "foobar" {
        error(t, "get", "foobar", res, err)
    }

    if _, err := Get(c, "foobar"); err == nil {
        error(t, "get", "error", nil, err)
    }

    if res, err := Getrange(c, "foo", 0, 2); err != nil || res.String() != "foo" {
        error(t, "getrange", "foo", res, err)
    }

    if res, err := Setrange(c, "foo", 0, "qux"); err != nil || res != 6 {
        error(t, "setrange", 6, res, err)
    }

    if res, err := Getset(c, "foo", "foo"); err != nil || res.String() != "quxbar" {
        error(t, "getset", "quxbar", res, err)
    }

    if res, err := Setnx(c, "foo", "bar"); err != nil || res != false {
        error(t, "setnx", false, res, err)
    }

    if res, err := Strlen(c, "foo"); err != nil || res != 3 {
        error(t, "strlen", 3, res, err)
    }

    if err := Setex(c, "foo", 10, "bar"); err != nil {
        error(t, "setex", nil, nil, err)
    }

    out := []string{"foo", "bar", "qux"}
    in := map[string]string{"foo": "foo", "bar": "bar", "qux": "qux"}

    if err := Mset(c, in); err != nil {
        error(t, "mset", nil, nil, err)
    }

    if res, err := Msetnx(c, in); err != nil || res == true {
        error(t, "msetnx", false, res, err)
    }

    res, err := Mget(c, append([]string{"il"}, out...)...)

    if err != nil || len(res.Elems) != 3 {
        error(t, "mget", 3, len(res.Elems), err)
        t.FailNow()
    }

    for i, v := range res.StringArray() {
        if v != out[i] {
            error(t, "mget", out[i], v, nil)
        }
    }
}

func TestList(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := Lpush(c, "foobar", "foo"); err != nil || res != 1 {
        error(t, "LPUSH", 1, res, err)
    }

    if res, err := Linsert(c, "foobar", "AFTER", "foo", "bar"); err != nil || res != 2 {
        error(t, "Linsert", 2, res, err)
    }

    if res, err := Linsert(c, "foobar", "AFTER", "qux", "bar"); err != nil || res != -1 {
        error(t, "Linsert", -1, res, err)
    }

    if res, err := Llen(c, "foobar"); err != nil || res != 2 {
        error(t, "Llen", 2, res, err)
    }

    if res, err := Lindex(c, "foobar", 0); err != nil || res.String() != "foo" {
        error(t, "Lindex", "foo", res, err)
    }

    if res, err := Lpush(c, "foobar", "qux"); err != nil || res != 3 {
        error(t, "Lpush", 3, res, err)
    }

    if res, err := Lpop(c, "foobar"); err != nil || res.String() != "qux" {
        error(t, "Lpop", "qux", res, err)
    }

    want1 := []string{"foo", "bar"}

    if out, err := Lrange(c, "foobar", 0, 1); err != nil || !reflect.DeepEqual(want1, out.StringArray()) {
        error(t, "Lrange", nil, nil, err)
    }

    want := []string{"foo"}

    if res, err := Lrem(c, "foobar", 0, "bar"); err != nil || res != 1 {
        error(t, "Lrem", 1, res, err)
    }

    want = []string{"bar"}

    if err := Lset(c, "foobar", 0, "bar"); err != nil {
        error(t, "Lrem", nil, nil, err)
    }

    want = []string{}

    if err := Ltrim(c, "foobar", 1, 0); err != nil {
        error(t, "Ltrim", nil, nil, err)
    }

    want = []string{"foo", "bar", "qux"}
    var res int64
    var err os.Error

    for _, v := range want {
        res, err = Rpush(c, "foobar", v)
    }

    if err != nil || res != 3 {
        error(t, "Rpush", 3, res, err)
    }

    if res, err := Rpushx(c, "foobar", "baz"); err != nil || res != 4 {
        error(t, "Rpushx", 4, res, err)
    }

    if res, err := Rpop(c, "foobar"); err != nil || res.String() != "baz" {
        error(t, "Rpop", "baz", res, err)
    }

    if res, err := Rpoplpush(c, "foobar", "foobaz"); err != nil || res.String() != "qux" {
        error(t, "Rpop", "qux", res, err)
    }
}

func TestHash(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := Hset(c, "foobar", "foo", "foo"); err != nil || res != true {
        error(t, "Hset", true, res, err)
    }

    if res, err := Hset(c, "foobar", "foo", "foo"); err != nil || res != false {
        error(t, "Hset", false, res, err)
    }

    if res, err := Hget(c, "foobar", "foo"); err != nil || res.String() != "foo" {
        error(t, "Hget", "foo", res, err)
    }

    if res, err := Hdel(c, "foobar", "foo"); err != nil || res != true {
        error(t, "Hdel", true, res, err)
    }

    if res, err := Hexists(c, "foobar", "foo"); err != nil || res != false {
        error(t, "Hexists", false, res, err)
    }

    if res, err := Hsetnx(c, "foobar", "foo", 1); err != nil || res != true {
        error(t, "Hsetnx", true, res, err)
    }
    Hset(c, "foobar", "bar", 2)

    want := []string{"foo", "1", "bar", "2"}

    if res, err := Hgetall(c, "foobar"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Hgetall", want, res, err)
    }

    if res, err := Hincrby(c, "foobar", "foo", 1); err != nil || int64(2) != res {
        error(t, "Hincrby", int64(2), res, err)
    }

    want1 := []string{"foo", "bar"}

    if res, err := Hkeys(c, "foobar"); err != nil || !reflect.DeepEqual(want1, res) {
        error(t, "Hkeys", want1, res, err)
    }

    if res, err := Hlen(c, "foobar"); err != nil || int64(2) != res {
        error(t, "Hlen", int64(2), res, err)
    }

    if res, err := Hlen(c, "foobar"); err != nil || int64(2) != res {
        error(t, "Hlen", int64(2), res, err)
    }

    want = []string{"2"}

    if res, err := Hmget(c, "foobar", "bar"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Hgetall", want, res, err)
    }

    m := map[string]interface{}{
        "foo": 1,
        "bar": 2,
        "qux": 3,
    }

    if err := Hmset(c, "foobar", m); err != nil {
        error(t, "Hmset", nil, nil, err)
    }

    want2 := []int64{1, 2, 3}
    if res, err := Hvals(c, "foobar"); err != nil || !reflect.DeepEqual(want2, res.IntArray()) {
        error(t, "Hvals", want2, res, err)
    }
}

func TestSet(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := Sadd(c, "foobar", "foo"); err != nil || res != true {
        error(t, "Sadd", true, res, err)
    }

    if res, err := Sadd(c, "foobar", "foo"); err != nil || res != false {
        error(t, "Sadd", false, res, err)
    }

    if res, err := Scard(c, "foobar"); err != nil || res != 1 {
        error(t, "Scard", 1, res, err)
    }

    Sadd(c, "foobar", "bar")
    Sadd(c, "foobaz", "foo")

    want := []string{"foo", "bar"}

    if res, err := Sunion(c, "foobar", "foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Sunion", want, res, err)
    }

    if res, err := Sunionstore(c, "fooqux", "foobar", "foobaz"); err != nil || res != 2 {
        error(t, "Sunionstore", 2, res, err)
    }

    want = []string{"bar"}

    if res, err := Sdiff(c, "foobar", "foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Sdiff", want, res, err)
    }

    if res, err := Sdiffstore(c, "foobar", "foobaz"); err != nil || res != 1 {
        error(t, "Sdiffstore", 1, res, err)
    }

    want = []string{"foo"}

    if res, err := Sinter(c, "foobar", "foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Sinter", want, res, err)
    }

    if res, err := Sinterstore(c, "foobar", "foobaz"); err != nil || res != 1 {
        error(t, "Sinterstore", 1, res, err)
    }

    if res, err := Sismember(c, "foobar", "qux"); err != nil || res != false {
        error(t, "Sismember", false, res, err)
    }

    if res, err := Smembers(c, "foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "smembers", want, res, err)
    }

    if res, err := Smove(c, "foobar", "foobaz", "foo"); err != nil || res != true {
        error(t, "smove", true, res, err)
    }

    if res, err := Spop(c, "foobaz"); err != nil || res.String() != "foo" {
        error(t, "spop", "foo", res, err)
    }

    if res, err := Srandmember(c, "foobaz"); err != nil || res != nil {
        error(t, "srandmember", nil, res, err)
    }

    Sadd(c, "foobar", "foo")
    Sadd(c, "foobar", "bar")
    Sadd(c, "foobar", "baz")

    if res, err := Srem(c, "foobar", "baz"); err != nil || res != true {
        error(t, "srem", nil, res, err)
    }
}

func TestSortedSet(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    m := map[string]float64{
        "foo": 1.0,
        "bar": 2.0,
        "baz": 3.0,
        "qux": 4.0,
    }

    for k, v := range m {
        if res, err := Zadd(c, "foobar", v, k); err != nil || res != true {
            error(t, "Zadd", true, res, err)
        }
    }

    if res, err := Zcard(c, "foobar"); err != nil || res != 4 {
        error(t, "Zcard", 4, res, err)
    }

    if res, err := Zcount(c, "foobar", 1, 2); err != nil || res != 2 {
        error(t, "Zcount", 2, res, err)
    }
}

func TestConnection(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := Echo(c, "foo"); err != nil || res.String() != "foo" {
        error(t, "Echo", "foo", res, err)
    }

    if res, err := Ping(c); err != nil || res.String() != "PONG" {
        error(t, "Ping", "PONG", res, err)
    }

    Set(c, "foo", "foo")

    if err := Select(c, 2); err != nil {
        error(t, "select", nil, nil, err)
    }

    if _, err := Get(c, "foo"); err == nil {
        error(t, "select", nil, nil, err)
    }

    // know bug will return EOF, but connection will not be restared
    //for i := 0; i < MaxClientConn; i++ {
    //    if err := Quit(c); err != nil {
    //        error(t, "quite", nil, nil, err)
    //    }
    //}

    //if err := Set(c, "foo", "foo"); err != nil {
    //    error(t, "quit", nil, nil, err)
    //}
}

func BenchmarkRpush(b *testing.B) {
    c := New("", 0, "")
    start := time.Nanoseconds()
    for i := 0; i < b.N; i++ {
        if _, err := Rpush(c, "qux", "qux"); err != nil {
            log.Println("RPUSH", err)
            return
        }
    }
    Del(c, "qux")
    stop := time.Nanoseconds() - start
    log.Printf("time: %.2f\n", float32(stop/1.0e+6)/1000.0)
}
