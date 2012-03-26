package godis

import (
    "errors"
    "net"
)

// New connection
func NewConn(addr, proto string) (net.Conn, error) {
    conn, err := net.Dial(proto, addr)

    if err != nil {
        return nil, errors.New("ERR " + proto + ":" + addr)
    }

    return conn, nil
}
