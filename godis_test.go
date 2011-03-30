package godis

import (
    "testing"
    "bytes"
    "bufio"
    "os"
    "strconv"
    "time"
    "fmt"
)

type simpleParserTest struct {
    in   string
    out  interface{}
    name string
    err  os.Error
}

func dummyReadWriter(data string) *redisReadWriter {
    r := bufio.NewReader(bytes.NewBufferString(data))
    w := bufio.NewWriter(bytes.NewBufferString(data))
    return &redisReadWriter{w, r}
}

var simpleParserTests = []simpleParserTest{
    {"+OK\r\n", "OK", "ok", nil},
    {"-ERR\r\n", nil, "err", os.NewError("ERR")},
    {":1\r\n", int64(1), "num", nil},
    {"$3\r\nfoo\r\n", s2Bytes("foo"), "bulk", nil},
    {"$-1\r\n", nil, "bulk-nil", nil},
    {"*-1\r\n", nil, "multi-bulk-nil", nil},
}

func TestParser(t *testing.T) {
    for _, test := range simpleParserTests {
        rw := dummyReadWriter(test.in)
        res, err := rw.read()

        if err != nil && test.err == nil {
            t.Errorf("'%s': unexpected error %v", test.name, err)
            t.FailNow()
        }

        switch v := res.(type) {
            case []byte:
                for i, c := range res.([]byte) {
                    if c != test.out.([]byte)[i] {
                        t.Errorf("expected %v got %v", test.out, res)
                    }
                }
            case [][]byte:
                for _, b := range res.([][]byte) {
                    for j, c := range b {
                        if c != test.out.([]byte)[j] {
                            t.Errorf("expected %v got %v", test.out, res)
                        }
                    }
                }
            default:
                if res != test.out {
                    t.Errorf("'%s': expected %s got %v", test.name, test.out, res)
                }
        }
        t.Log(test.in, res, test.out)
    }
}

func s2Bytes(s string) []byte {
    return bytes.NewBufferString(s).Bytes()
}

func s2MultiBytes(ss ...string) [][]byte {
    var buf = make([][]byte, len(ss))
    for i := 0; i < len(ss); i++ {
        buf[i] = s2Bytes(ss[i])
    }
    return buf
}

type SimpleSendTest struct {
    cmd  string
    args []string
    out  interface{}
}

var simpleSendTests = []SimpleSendTest{
    {"FLUSHDB", []string{}, "OK"},
    {"SET", []string{"key", "foo"}, "OK"},
    {"EXISTS", []string{"key"}, int64(1)},
    {"GET", []string{"key"}, s2Bytes("foo")},
    {"RPUSH", []string{"list", "foo"}, int64(1)},
    {"RPUSH", []string{"list", "bar"}, int64(2)},
    {"LRANGE", []string{"list", "0", "2"}, s2MultiBytes("foo", "bar")},
    {"KEYS", []string{"list"}, s2MultiBytes("list")},
    {"GET", []string{"/dev/null"}, nil},
}

func TestSimpleSend(t *testing.T) {
    c := New("", 0, "")
    for _, test := range simpleSendTests {
        res, err := c.SendStr(test.cmd, test.args...)

        if err != nil {
            t.Errorf("'%s': unexpeced error %q", test.cmd, err)
            t.FailNow()
        }

        switch v := res.(type) {
        case []byte:
            for i, c := range res.([]byte) {
                if c != test.out.([]byte)[i] {
                    t.Errorf("'%s': expected %v got %v", test.cmd, test.out, res)
                }
            }
        case [][]byte:
            res_arr := res.([][]byte)
            out_arr := test.out.([][]byte)

            for i := 0; i < len(res_arr); i++ {
                for j := 0; j < len(res_arr[i]); j++ {
                    if res_arr[i][j] != out_arr[i][j] {
                        t.Errorf("'%s': expected %v got %v", test.cmd, test.out, res)
                    }
                }
            }
        default:
            if res != test.out {
                t.Errorf("'%s': expected %v got %v", test.cmd, test.out, res)
            }
        }
        t.Log(test.cmd, test.args, test.out)
    }
}

func BenchmarkParsing(b *testing.B) {
    c := New("", 0, "")

    for i := 0; i < 1000; i++ {
        c.SendStr("RPUSH", "list", "foo")
    }

    start := time.Nanoseconds()

    for i := 0; i < b.N; i++ {
        c.SendStr("LRANGE", "list", "0", "50")
    }

    stop := time.Nanoseconds() - start

    fmt.Fprintf(os.Stdout, "time: %.2f\n", float32(stop / 1.0e+6) / 1000.0)
    c.SendStr("FLUSHDB")
}

func error(t *testing.T, name string, expected, got interface{}, err os.Error) {
    t.Errorf("`%s` expected `%v` got `%v`, err(%v)", name, expected, got, err)
}

func TestGeneric(t *testing.T) {
    c := New("", 0, "")
    c.Send("FLUSHDB")

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
    if res, err := c.Expireat("foo", time.Seconds() + 10); !res {
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
    c.SendStr("FLUSHDB")
    c.SendStr("MSET", "foo", "one", "bar", "two", "baz", "three")

    res, err := c.Keys("foo"); 

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
    c.SendStr("FLUSHDB")
    c.SendStr("RPUSH", "foo", "2") 
    c.SendStr("RPUSH", "foo", "3") 
    c.SendStr("RPUSH", "foo", "1") 

    res, err := c.Sort("foo")

    if err != nil {
        error(t, "sort", nil, nil, err)
    }

    expected := []int{1, 2, 3}

    if len(res) != len(expected) {
        error(t, "sort", len(res), len(expected), nil)
    }
    
    for i, v := range res {
        r, _ := strconv.Atoi(string(v))
        if r != expected[i] {
            error(t, "sort", expected[i], v, nil)
        }
    }
}

func TestString(t *testing.T) {
    c := New("", 0, "")
    c.SendStr("FLUSHDB")

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
    in := map[string]string{"foo":"foo", "bar":"bar", "qux":"qux"}

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
    c.SendStr("FLUSHDB")

    if res, err := c.Lpush("foobar", []byte(strconv.Itoa(1))); err != nil || res != 1 {
        error(t, "LPUSH", 1, res, err)
    }
}
