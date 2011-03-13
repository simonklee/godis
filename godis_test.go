package godis

import (
    "testing"
    "bufio"
    "bytes"
    "log"
    "os"
)

const (
    LOG = true
)

var logger = log.New(os.Stderr, "", 0x00)

func l(args ...interface{}) {
    if LOG {
        logger.Println(args...)
    }
}

func s2bytes(s string) []byte {
    return bytes.NewBufferString(s).Bytes()
}

type CmdGoodTest struct {
    cmd  string
    args []string
    out  interface{}
}

var cmdGoodTests = []CmdGoodTest{
    {"FLUSHDB", []string{}, "OK"},
    {"SET", []string{"key", "foo"}, "OK"},
    {"EXISTS", []string{"key"}, int64(1)},
    {"GET", []string{"key"}, s2bytes("foo")},
    {"RPUSH", []string{"list", "foo"}, int64(1)},
}

func TestGoodSend(t *testing.T) {
    var c Client
    for _, test := range cmdGoodTests {
        res, err := c.Send(test.cmd, test.args...)

        if err != nil {
            t.Errorf("unexpeced error %q", err)
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
                t.Errorf("'%s': expected %v got %v", test.cmd, test.out, res)
            }
        }
        l(test.cmd, test.args, test.out)
    }
}

type simpleParserTest struct {
    in   string
    out  interface{}
    name string
    err  os.Error
}

var simpleParserTests = []simpleParserTest{
    {"+OK\r\n", "OK", "ok", nil},
    {"-ERR\r\n", nil, "err", os.NewError("ERR")},
    {":1\r\n", int64(1), "num", nil},
    {"$3\r\nfoo\r\n", s2bytes("foo"), "bulk", nil},
    {"$-1\r\n", nil, "bulk-nil", nil},
    {"*-1\r\n", nil, "multi-bulk-nil", nil},
}

func reader(data string) *bufio.Reader {
    b := bufio.NewReader(bytes.NewBufferString(data))
    return b
}

func TestParser(t *testing.T) {
    for _, test := range simpleParserTests {
        res, err := Read(reader(test.in))

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
        //l(test.in, res, test.out)
    }
}

//type CmdBadTest struct {
//    cmd string
//    args []string
//    out interface{}
//    err os.Error
//}
//
//var cmdBadTests = []CmdBadTest{
//    {"EXISTS", []string{"bar"}, 0, nil},
//}
//
//func TestBadSend(t *testing.T) {
//    var c Client
//    for _, test := range cmdBadTests {
//        // TODO: implement test
//        res, err := c.Send(test.cmd, test.args...)
//
//        if err != test.err {
//            t.Errorf("expected error %v got %v", test.err, res)
//        }
//
//        if res != test.out {
//        }
//    }
//}
//    // client.write(bytesCommand("GET", "keylist"))
//    // client.write(bytesCommand("GET", "nonexistant"))
//    //client.send("LRANGE", "keylist", "0", "4")
//    //client.send("KEYS", "*")
