package godis

import (
    "testing"
    "reflect"
    "os"
)

type CmdGoodTest struct {
    cmd string
    args []string
    out interface{}
}

var cmdGoodTests = []CmdGoodTest{
    {"EXISTS", []string{"key"}, 1},
}

func TestGoodSend(t *testing.T) {
    var c Client
    for _, test := range cmdGoodTests {
        res, err := c.Send(test.cmd, test.args...)

        if err != nil {
            t.Errorf("unexpeced error %q", err)
            t.FailNow()
        }

        r_typ := reflect.Typeof(res)
        t_typ := reflect.Typeof(test.out)

        if  r_typ != t_typ {
            t.Errorf("expected typeof %v got %v", t_typ, r_typ)
        }

        if res != test.out {
            t.Errorf("expected %v got %v", test.out, res)
        }
    }
}

type CmdBadTest struct {
    cmd string
    args []string
    out interface{}
    err os.Error
}

var cmdBadTests = []CmdBadTest{
    {"EXISTS", []string{"bar"}, 0, nil},
}

func TestBadSend(t *testing.T) {
    var c Client
    for _, test := range cmdBadTests {
        // TODO: implement test
        res, err := c.Send(test.cmd, test.args...)

        if err != test.err {
            t.Errorf("expected error %v got %v", test.err, res)
        }

        if res != test.out {
        }
    }
}
//func main() {
//    //var client Client = Client{"127.0.0.1", 6379, 0, nil} 
//    var client Client
//    log(client.Host)
//    //client := new(Client)
//
//    // var enc_set []byte = bytesCommand("SET", "key", "hello")
//    // fmt.Printf("%q\n", enc_set)
//
//    // var enc_get []byte = bytesCommand("GET", "key")
//    // fmt.Printf("%q\n", enc_get)
//
//    // client.write(enc_set)
//    // client.write(enc_get)
//    //client.send("RPUSH", "keylist", "two")
//    // client.write(bytesCommand("GET", "keylist"))
//    // client.write(bytesCommand("GET", "nonexistant"))
//    //client.send("GET", "key")
//    // client.send("SET", "key", "Hello")
//    //client.send("LRANGE", "keylist", "0", "4")
//    //client.send("KEYS", "*")
//    //client.send("EXISTS", "key")
//}
