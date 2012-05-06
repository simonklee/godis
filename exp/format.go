package redis

import (
    "bytes"
    "fmt"
    "strconv"
)

const (
    cr     byte = 13
    lf     byte = 10
    dollar byte = 36
    colon  byte = 58
    minus  byte = 45
    plus   byte = 43
    star   byte = 42
)

var (
    delim = []byte{cr, lf}
)

func intlen(n int) int {
    l := 1

    if n < 0 {
        n = -n
        l++
    }

    n /= 10

    for n > 9 {
        l++
        n /= 10
    }

    return l
}

func arglen(arg []byte) int {
    //     $   datalen   \r\n data \r\n
    return 1 + intlen(len(arg)) + 2 + len(arg) + 2
}

/* Build a new command by concencate an array 
 * of bytes which create a redis command.
 * Returns a byte array */
func formatArgs(args [][]byte) []byte {
    //   *   args count         \r\n
    n := 1 + intlen(len(args)) + 2

    for i := 0; i < len(args); i++ {
        n += arglen(args[i])
    }

    buf := make([]byte, 0, n)
    buf = append(buf, star)
    buf = strconv.AppendUint(buf, uint64(len(args)), 10)
    buf = append(buf, delim...)

    for _, arg := range args {
        buf = append(buf, dollar)
        buf = strconv.AppendUint(buf, uint64(len(arg)), 10)
        buf = append(buf, delim...)
        buf = append(buf, arg...)
        buf = append(buf, delim...)
    }

    return buf
}

/* Build a new command by concencate an array 
 * of strings which create a redis command.
 * Returns a byte array */
func format(args ...interface{}) []byte {
    buf := make([][]byte, len(args))

    for i, arg := range args {
        switch v := arg.(type) {
        case []byte:
            buf[i] = v
        case nil:
            buf[i] = []byte(nil)
        case string:
            buf[i] = []byte(v)
        default:
            var b bytes.Buffer
            fmt.Fprint(&b, v)
            buf[i] = b.Bytes()
        }
    }

    return formatArgs(buf)
}
