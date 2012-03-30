package godis

import (
    "bytes"
    "strconv"
    "testing"
)

func error_(t *testing.T, name string, expected, got interface{}, err error) {
    if err != nil {
        t.Errorf("`%s` expected `%v` got `%v`, err(%v)", name, expected, got, err.Error())
    } else {
        t.Errorf("`%s` expected `%v` got `%v`, err(%v)", name, expected, got, err)
    }
}

func formatTest(t *testing.T, exp string, a ...string) {
    got := format(a...)

    if exp != string(got) {
        t.Errorf("format: expected %s got %s", exp, string(got))
    }
}

func TestFormat(t *testing.T) {
    formatTest(t, "*2\r\n$4\r\nPING\r\n$4\r\nPONG\r\n", "PING", "PONG")
    formatTest(t, "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", "SET", "foo", "bar")
    formatTest(t, "*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n", "GET", "foo")
}

func TestClient(t *testing.T) {
    c := NewClient("")

    if _, err := c.Call("SET", "foo", "foo"); err != nil {
        println("errror call")
        t.Fatal(err.Error())
    }

    p := c.Pipeline()
    p.Call("MULTI")
    p.Call("GET", "foo")
    p.Call("EXEC")

    res, err := p.Read()

    if err != nil || string(res.Elem) != "OK" {
        t.Fatal(err.Error())
    }

    res, err = p.Read()

    if err != nil || string(res.Elem) != "QUEUED" {
        error_(t, "pipe", "foo", string(res.Elem), err)
    }

    res, err = p.Read()

    if err != nil || len(res.Elems) != 1 {
        error_(t, "exec", 1, len(res.Elems), err)
    } else {
        println(string(res.Elems[0].Elem))
    }
}

func BenchmarkItoa(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strconv.Itoa(i)
    }
}

func BenchmarkSet(b *testing.B) {
    c := NewClient("")

    for i := 0; i < b.N; i++ {
        c.Call("SET", "foo", "foo")
    }
}

func BenchmarkAppendUint(b *testing.B) {
    var buf []byte
    buf = make([]byte, 0, 1024*16)

    for i := 0; i < b.N; i++ {
        strconv.AppendUint(buf, uint64(i), 10)
    }
}

func BenchmarkAppendBytes(b *testing.B) {
    var buf []byte
    buf = make([]byte, 0, 1024*16)

    for i := 0; i < b.N; i++ {
        buf = append(buf, '\r')
    }
}

func BenchmarkAppendBuffer(b *testing.B) {
    buf := bytes.NewBuffer(make([]byte, 1024*16))

    for i := 0; i < b.N; i++ {
        buf.WriteByte('\r')
    }
}

//package godis
//import (
//    "bufio"
//    "bytes"
//    "errors"
//    "log"
//    "reflect" //"strconv"
//
//    "testing"
//    "time"
//)

//func error_(t *testing.T, name string, expected, got interface{}, err error) {
//    t.Errorf("`%s` expected `%v` got `%v`, err(%v)", name, expected, got, err)
//}
//
//func printRes(t *testing.T, r *Reply) {
//    if len(r.Elems) > 0 {
//        t.Logf("str arr: %q", r.StringArray())
//    } else {
//        t.Logf("str: %q", r.Elem.String())
//    }
//    if r.Err != nil {
//        t.Logf("err: %q", r.Err)
//    }
//}
//
//func compareReply(t *testing.T, name string, a, b *Reply) {
//    if a.Err != nil && b.Err == nil {
//        t.Fatalf("'%s': expected error `%v`", name, a.Err)
//    } else if b.Err != nil && b.Err.Error() != b.Err.Error() {
//        t.Fatalf("'%s': expected %s got %v", name, a.Err, b.Err)
//    } else if b.Elem != nil {
//        for i, c := range a.Elem {
//            if c != b.Elem[i] {
//                t.Errorf("'%s': expected %v got %v", name, b, a)
//            }
//        }
//    } else if b.Elems != nil {
//        for i, rep := range a.Elems {
//            for j, e := range rep.Elem {
//                if e != b.Elems[i].Elem[j] {
//                    t.Errorf("expected %v got %v", b, a)
//                    break
//                }
//            }
//        }
//    }
//}
//
//type simpleParserTest struct {
//    in   string
//    out  Reply
//    name string
//}
//
//type redisReadWriter struct {
//    writer *bufio.Writer
//    reader *bufio.Reader
//}
//
//func dummyReadWriter(data string) *conn {
//    br := bufio.NewReader(bytes.NewBufferString(data))
//    bw := bufio.NewWriter(bytes.NewBufferString(data))
//    return &conn{rwc: nil, r: br, w: bw}
//}
//
//var simpleParserTests = []simpleParserTest{
//    {"+OK\r\n", Reply{Elem: []byte("OK")}, "ok"},
//    {"-ERR\r\n", Reply{Err: errors.New("ERR")}, "err"},
//    {":1\r\n", Reply{Elem: []byte("1")}, "num"},
//    {"$3\r\nfoo\r\n", Reply{Elem: []byte("foo")}, "bulk"},
//    {"$-1\r\n", Reply{Err: errors.New("Nonexisting Key")}, "bulk-nil"},
//    {"*-1\r\n", Reply{}, "multi-bulk-nil"},
//}
//
//func TestParser(t *testing.T) {
//    for _, test := range simpleParserTests {
//        rw := dummyReadWriter(test.in)
//        r := rw.readReply()
//        compareReply(t, test.name, r, &test.out)
//        t.Log(test.in, r, test.out)
//    }
//}
//
//func s2MultiReply(ss ...string) []*Reply {
//    var r = make([]*Reply, len(ss))
//    for i := 0; i < len(ss); i++ {
//        r[i] = &Reply{Elem: []byte(ss[i])}
//    }
//    return r
//}
//
//type SimpleSendTest struct {
//    cmd  string
//    args []string
//    out  Reply
//}
//
//var simpleSendTests = []SimpleSendTest{
//    {"FLUSHDB", []string{}, Reply{Elem: []byte("OK")}},
//    {"SET", []string{"key", "foo"}, Reply{Elem: []byte("OK")}},
//    {"EXISTS", []string{"key"}, Reply{Elem: []byte("1")}},
//    {"GET", []string{"key"}, Reply{Elem: []byte("foo")}},
//    {"GET", []string{"/dev/null"}, Reply{}},
//    {"RPUSH", []string{"list", "foo"}, Reply{Elem: []byte("1")}},
//    {"RPUSH", []string{"list", "bar"}, Reply{Elem: []byte("2")}},
//    {"LRANGE", []string{"list", "0", "2"}, Reply{Elems: s2MultiReply("foo", "bar")}},
//    {"KEYS", []string{"list"}, Reply{Elems: s2MultiReply("list")}},
//}
//
//func TestSimpleSend(t *testing.T) {
//    c := New("", 0, "")
//    for _, test := range simpleSendTests {
//        r := SendStr(c.Rw, test.cmd, test.args...)
//        compareReply(t, test.cmd, &test.out, r)
//        t.Log(test.cmd, test.args)
//        t.Logf("%q == %q\n", test.out, r)
//    }
//}
//
//func equals(a, b []byte) bool {
//    for i, c := range a {
//        if c != b[i] {
//            return false
//        }
//    }
//    return true
//}
//
//func TestBinarySafe(t *testing.T) {
//    c := New("", 0, "")
//    want1 := make([]byte, 256)
//    for i := 0; i < 256; i++ {
//        want1[i] = byte(i)
//    }
//
//    Send(c.Rw, []byte("SET"), []byte("foo"), want1)
//
//    if res := Send(c.Rw, []byte("GET"), []byte("foo")); !equals(res.Elem.Bytes(), want1) {
//        error_(t, "ascii-table-Send", want1, res.Elem.Bytes(), res.Err)
//    }
//
//    SendIface(c.Rw, "SET", "bar", string(want1))
//
//    if res := SendIface(c.Rw, "GET", "bar"); !equals(res.Elem.Bytes(), want1) {
//        error_(t, "ascii-table-SendIface", want1, res.Elem.Bytes(), res.Err)
//    }
//
//    want2 := []byte("♥\r\nµs\r\n")
//    Send(c.Rw, []byte("SET"), []byte("foo"), want2)
//
//    if res := Send(c.Rw, []byte("GET"), []byte("foo")); !equals(res.Elem.Bytes(), want2) {
//        error_(t, "unicode-Send", want2, res.Elem.Bytes(), res.Err)
//    }
//
//    SendIface(c.Rw, "SET", "bar", want2)
//
//    if res := SendIface(c.Rw, "GET", "bar"); !equals(res.Elem.Bytes(), want2) {
//        error_(t, "unicode-SendIface", want2, res.Elem.Bytes(), res.Err)
//    }
//
//    for _, b := range want2 {
//        SendIface(c.Rw, "SET", "bar", b)
//        res := SendIface(c.Rw, "GET", "bar")
//        if uint8(res.Elem.Int64()) != b {
//            error_(t, "unicode-SendIface", b, res.Elem, res.Err)
//        }
//    }
//}
//
//func TestSimplePipe(t *testing.T) {
//    c := NewPipeClient("", 0, "")
//
//    for _, test := range simpleSendTests {
//        r := SendStr(c.Rw, test.cmd, test.args...)
//
//        if r.Err != nil {
//            t.Error(test.cmd, r.Err, test.args)
//        }
//    }
//
//    replies := c.Exec()
//
//    if len(replies) != len(simpleSendTests) {
//        error_(t, "pipe replies len", len(simpleSendTests), len(replies), nil)
//    }
//
//    for i, test := range simpleSendTests {
//        compareReply(t, test.cmd, &test.out, replies[i])
//    }
//}
//
//func TestSimpleTransaction(t *testing.T) {
//    c := New("", 0, "")
//
//    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
//        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
//    }
//
//    p := NewPipeClientFromClient(c)
//    p.Multi()
//    p.Set("foo", "bar")
//    p.Set("bar", "bar")
//    p.Lpush("bar", "bar")
//    p.Get("bar")
//    replies := p.Exec()
//
//    t.Log(replies)
//    t.Log(replies[2].Err)
//
//    pc := NewPipeClient("", 0, "")
//    pc.Set("baz", "baz")
//    pc.Exec()
//}
//
//// for this test to work redis.conf has to be set timeout to 1sec
//// the test return a nil pointer if failed
//func TestConnTimeout(t *testing.T) {
//    c := New("", 0, "")
//    Send(c.Rw, []byte("FLUSHDB"))
//
//    defer func() {
//        if x := recover(); x != nil {
//            t.Errorf("`conn timeout` expected got `%v`", x)
//        }
//    }()
//
//    c.Set("foo", 1)
//    c.Set("bar", 2)
//
//    time.Sleep(1e+9 * 8)
//
//    rep, err := c.Mget("foo", "bar")
//    // rep.IntArray will invoke a nil-pointer panic if there was an err
//    rep.IntArray()
//
//    if err != nil {
//        error_(t, "connection timeout", nil, nil, err)
//    }
//}
//
//func TestReadingBulk(t *testing.T) {
//    c := New("", 0, "")
//
//    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
//        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
//    }
//
//    var want3 []int64
//
//    for i := 0; i < 600; i++ {
//        want3 = append(want3, int64(i))
//        c.Rpush("foobaz", i)
//
//        if res, err := c.Lrange("foobaz", 0, i); err != nil || !reflect.DeepEqual(want3, res.IntArray()) {
//            error_(t, "Lranges", nil, nil, err)
//            t.FailNow()
//        }
//    }
//}
