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
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }
    SendStr(c, "MSET", "foo", "one", "bar", "two", "baz", "three")

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
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }
    SendStr(c, "RPUSH", "foo", "2")
    SendStr(c, "RPUSH", "foo", "3")
    SendStr(c, "RPUSH", "foo", "1")

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
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

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

    if res, err := c.Get("foo"); err != nil || res.String() != "foobar" {
        error(t, "get", "foobar", res, err)
    }

    if _, err := c.Get("foobar"); err == nil {
        error(t, "get", "error", nil, err)
    }

    if res, err := c.Getrange("foo", 0, 2); err != nil || res.String() != "foo" {
        error(t, "getrange", "foo", res, err)
    }

    if res, err := c.Setrange("foo", 0, "qux"); err != nil || res != 6 {
        error(t, "setrange", 6, res, err)
    }

    if res, err := c.Getset("foo", "foo"); err != nil || res.String() != "quxbar" {
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

    if res, err := c.Lpop("foobar"); err != nil || res.String() != "qux" {
        error(t, "Lpop", "qux", res, err)
    }

    want1 := []string{"foo", "bar"}

    if out, err := c.Lrange("foobar", 0, 1); err != nil || !reflect.DeepEqual(want1, out.StringArray()) {
        error(t, "Lrange", nil, nil, err)
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
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

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

    if res, err := c.Hsetnx("foobar", "foo", 1); err != nil || res != true {
        error(t, "Hsetnx", true, res, err)
    }
    c.Hset("foobar", "bar", 2)

    want := []string{"foo", "1", "bar", "2"}

    if res, err := c.Hgetall("foobar"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Hgetall", want, res, err)
    }

    if res, err := c.Hincrby("foobar", "foo", 1); err != nil || int64(2) != res {
        error(t, "Hincrby", int64(2), res, err)
    }

    want1 := []string{"foo", "bar"}

    if res, err := c.Hkeys("foobar"); err != nil || !reflect.DeepEqual(want1, res) {
        error(t, "Hkeys", want1, res, err)
    }

    if res, err := c.Hlen("foobar"); err != nil || int64(2) != res {
        error(t, "Hlen", int64(2), res, err)
    }

    if res, err := c.Hlen("foobar"); err != nil || int64(2) != res {
        error(t, "Hlen", int64(2), res, err)
    }

    want = []string{"2"}

    if res, err := c.Hmget("foobar", "bar"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Hgetall", want, res, err)
    }

    m := map[string]interface{}{
        "foo": 1,
        "bar": 2,
        "qux": 3,
    }

    if err := c.Hmset("foobar", m); err != nil {
        error(t, "Hmset", nil, nil, err)
    }

    want2 := []int64{1, 2, 3}
    if res, err := c.Hvals("foobar"); err != nil || !reflect.DeepEqual(want2, res.IntArray()) {
        error(t, "Hvals", want2, res, err)
    }
}

func TestSet(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := c.Sadd("foobar", "foo"); err != nil || res != true {
        error(t, "Sadd", true, res, err)
    }

    if res, err := c.Sadd("foobar", "foo"); err != nil || res != false {
        error(t, "Sadd", false, res, err)
    }

    if res, err := c.Scard("foobar"); err != nil || res != 1 {
        error(t, "Scard", 1, res, err)
    }

    c.Sadd("foobar", "bar")
    c.Sadd("foobaz", "foo")

    want := []string{"foo", "bar"}

    if res, err := c.Sunion("foobar", "foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Sunion", want, res, err)
    }

    if res, err := c.Sunionstore("fooqux", "foobar", "foobaz"); err != nil || res != 2 {
        error(t, "Sunionstore", 2, res, err)
    }

    want = []string{"bar"}

    if res, err := c.Sdiff("foobar", "foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Sdiff", want, res, err)
    }

    if res, err := c.Sdiffstore("foobar", "foobaz"); err != nil || res != 1 {
        error(t, "Sdiffstore", 1, res, err)
    }

    want = []string{"foo"}

    if res, err := c.Sinter("foobar", "foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Sinter", want, res, err)
    }

    if res, err := c.Sinterstore("foobar", "foobaz"); err != nil || res != 1 {
        error(t, "Sinterstore", 1, res, err)
    }

    if res, err := c.Sismember("foobar", "qux"); err != nil || res != false {
        error(t, "Sismember", false, res, err)
    }

    if res, err := c.Smembers("foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "smembers", want, res, err)
    }

    if res, err := c.Smove("foobar", "foobaz", "foo"); err != nil || res != true {
        error(t, "smove", true, res, err)
    }

    if res, err := c.Spop("foobaz"); err != nil || res.String() != "foo" {
        error(t, "spop", "foo", res, err)
    }

    if res, err := c.Srandmember("foobaz"); err != nil || res != nil {
        error(t, "srandmember", nil, res, err)
    }

    c.Sadd("foobar", "foo")
    c.Sadd("foobar", "bar")
    c.Sadd("foobar", "baz")

    if res, err := c.Srem("foobar", "baz"); err != nil || res != true {
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
        if res, err := c.Zadd("foobar", v, k); err != nil || res != true {
            error(t, "Zadd", true, res, err)
        }
    }

    if res, err := c.Zcard("foobar"); err != nil || res != 4 {
        error(t, "Zcard", 4, res, err)
    }

    if res, err := c.Zcount("foobar", 1, 2); err != nil || res != 2 {
        error(t, "Zcount", 2, res, err)
    }

    if res, err := c.Zincrby("foobar", 0.5, "foo"); err != nil || res != 1.5 {
        error(t, "Zincrby", 1.5, res, err)
    }

    if res, err := c.Zinterstore("barbaz", []string{"foobar"}); err != nil || res != 4 {
        error(t, "Zinterstore", 4, res, err)
    }

    want := []string{"foo", "bar", "baz", "qux"}

    if res, err := c.Zrange("foobar", 0, 4); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Zrange", want, res, err)
    }

    if res, err := c.Zrangebyscore("foobar", "0", "+inf"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Zrangebyscore", want, res, err)
    }

    if res, err := c.Zrank("foobar", "baz"); err != nil || res != 2 {
        error(t, "Zrank", 2, res, err)
    }

    // if res, err := c.Zrank("foobar", "nil"); err == nil || res != 0  {
    //    error(t, "Zrank", 0, res, err)
    // }

    if res, err := c.Zrem("foobar", "bar"); err != nil || res != true {
        error(t, "Zrem", true, res, err)
    }
 
    if res, err := c.Zrem("foobar", "bar"); err != nil || res != false {
        error(t, "Zrem", false, res, err)
    }

    if res, err := c.Zremrangebyrank("foobar", 0, 0); err != nil || res != 1 {
        error(t, "Zremrangebyrank", 1, res, err)
    }

    if res, err := c.Zremrangebyscore("foobar", 0, 3); err != nil || res != 1 {
        error(t, "zremrangebyscore", 1, res, err)
    }

    c.Zinterstore("foobar", []string{"barbaz"})
    want = []string{"qux", "baz", "bar", "foo"}
    
    if res, err := c.Zrevrange("foobar", 0, 4); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error(t, "Zrevrange", want, res.StringArray(), err)
    }

    want2 := map[string]string{"qux": "4", "baz": "3", "bar": "2", "foo": "1.5"}
    if res, err := c.Zrevrangebyscore("foobar", 4, 0, "WITHSCORES"); err != nil || !reflect.DeepEqual(want2, res.StringMap()) {
        error(t, "Zrevrangebyscore", want, res.StringMap(), err)
    }
 
    if res, err := c.Zrevrank("foobar", "baz"); err != nil || res != 1 {
        error(t, "Zrevrank", 1, res, err)
    }

    if res, err := c.Zscore("foobar", "foo"); err != nil || res != 1.5 {
        error(t, "Zscore", 1.5, res, err)
    }

    if res, err := c.Zunionstore("foobar", []string{"nil"}); err == nil || res != -1 {
        error(t, "Zscore", -1, res, err)
    }
}

func TestConnection(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c, "flushall"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := c.Echo("foo"); err != nil || res.String() != "foo" {
        error(t, "Echo", "foo", res, err)
    }

    if res, err := c.Ping(); err != nil || res.String() != "PONG" {
        error(t, "Ping", "PONG", res, err)
    }

    c.Set("foo", "foo")

    if err := c.Select(2); err != nil {
        error(t, "select", nil, nil, err)
    }

    if res, err := c.Get("foo"); err == nil || res != nil {
        error(t, "get-select", nil, res, err)
    }

    // know bug will return EOF, but connection will not be restared
    //for i := 0; i < MaxClientConn; i++ {
    //    if err := Quit(c); err != nil {
    //        error(t, "quite", nil, nil, err)
    //    }
    //}

    //if err := c.Set("foo", "foo"); err != nil {
    //    error(t, "quit", nil, nil, err)
    //}
}

func TestPubSub(t *testing.T) {
    c := New("", 0, "")

    if r := SendStr(c, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    sub, err := c.Subscribe("foochan", "barchan")

    if err != nil {
        t.Fatalf("subscribe", nil, nil, err)
    }

    if res, err := c.Publish("foochan", "foo"); err != nil || res != 1 {
        error(t, "publish", 1, res, err)
    }

    go func() {
        m := <-sub.Messages
        if m.Elem.String() != "foo" || m.Channel != "foochan" {
            error(t, "res-subscribe", "foo "+"foochan", m.Elem.String()+m.Channel, nil)
        }
    }()

    time.Sleep(1e8)
    if res, err := c.Publish("barchan", "bar"); err != nil || res != 1 {
        error(t, "publish", 1, res, err)
    }

    time.Sleep(1e8)

    if m := <-sub.Messages; m.Elem.String() != "bar" || m.Channel != "barchan" {
        error(t, "subscribe barchan", "bar", m, nil)
    }

    if err := sub.Unsubscribe("foochan"); err != nil {
        error(t, "unsubscribe", nil, nil, err)
    }

    if res, err := c.Publish("foochan", "foo"); err != nil || res != 0 {
        error(t, "publish", 0, res, err)
    }

    if res, err := c.Publish("barchan", "bar"); err != nil || res != 1 {
        error(t, "publish", 1, res, err)
    }

    if m := <-sub.Messages; m.Elem.String() != "bar" || m.Channel != "barchan" {
        error(t, "subscribe barchan", "bar", m.Elem.String(), nil)
    }

    if err := sub.Psubscribe("*chan"); err != nil {
        error(t, "psubscribe", nil, nil, err)
    }

    if res, err := c.Publish("bazchan", "baz"); err != nil || res != 1 {
        error(t, "publish", 1, res, err)
    }

    if m := <-sub.Messages; m.Elem.String() != "baz" || m.Channel != "bazchan" {
        error(t, "psubscribe bazchan", "baz", m.Elem.String(), nil)
    }

    sub.Close()
    time.Sleep(1e8)

    if _, ok := <-sub.Messages; ok != false {
        error(t, "closed chan", false, ok, nil)
    }

}

func BenchmarkRpush(b *testing.B) {
    c := New("", 0, "")
    start := time.Nanoseconds()
    for i := 0; i < b.N; i++ {
        if _, err := c.Rpush("qux", "qux"); err != nil {
            log.Println("RPUSH", err)
            return
        }
    }
    c.Del("qux")
    stop := time.Nanoseconds() - start
    log.Printf("time: %.2f\n", float32(stop/1.0e+6)/1000.0)
}
