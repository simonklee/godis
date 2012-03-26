package godis

import (
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

/* Build a new command by concencate an array 
 * of bytes which create a redis command.
 * Returns a byte array */
func formatArgs(args [][]byte) []byte {
    buf := make([]byte, 0, 16)
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
func format(args ...string) []byte {
    buf := make([][]byte, len(args))

    for i, arg := range args {
        buf[i] = []byte(arg)
    }

    return formatArgs(buf)
}
