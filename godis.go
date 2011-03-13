package godis

import (
    "fmt"
    "os"
    "bufio"
    "bytes"
    "strings"
    "strconv"
    "log"
)

type Client struct {
    Addr     string
    Db       int
    Password string
    pool     *Pool
}

var (
    defaultAddr = "localhost:6379"
)

func newError(format string, args ...interface{}) os.Error {
    return os.NewError(fmt.Sprintf(format, args...))
}

func errorReply(line string) (interface{}, os.Error) {
    // log.Println("GODIS: " + res)
    return nil, newError(line)
}

func singleReply(line string) (string, os.Error) {
    // log.Println("GODIS: " + res)
    return line, nil
}

func integerReply(line string) (int64, os.Error) {
    // log.Println("GODIS: " + res)
    n, err := strconv.Atoi64(line)
    return n, err
}

func bulkReply(line string, head *bufio.Reader) ([]byte, os.Error) {
    l, _ := strconv.Atoi(line)
    if l == -1 {
        return nil, nil
    }

    data := make([]byte, l)

    n, err := head.Read(data)
    if n != l || err != nil {
        if n != l {
            err = newError("expected %d bytes got %d bytes", l, n)
        }
        return nil, err
    }
    // log.Println("GODIS: bulk-len: " + strconv.Itoa(l))
    // log.Println("GODIS: bulk-value: " + string(data))
    // log.Printf("GODIS: %q\n", data)
    return data, nil
}

func multiBulkReply(line string, head *bufio.Reader) ([][]byte, os.Error) {
    l, _ := strconv.Atoi(line)

    if l == -1 {
        return nil, nil
    }

    // log.Println("GODIS: multi-bulk-len: " + strconv.Itoa(l))
    var data = make([][]byte, l)
    for i := 0; i < l; i++ {
        d, err := Read(head)
        if err != nil {
            log.Println("GODIS: returned with error")
            return nil, err
        }
        data[i] = d.([]byte)
    }

    // fmt.Printf("GODIS: %q\n", data)
    return data, nil
}

func Read(head *bufio.Reader) (interface{}, os.Error) {
    res, err := head.ReadString('\n')
    if err != nil {
        return nil, err
    }

    typ := res[0]
    line := strings.TrimSpace(res[1:])

    switch typ {
    case '-':
        return errorReply(line)
    case '+':
        return singleReply(line)
    case ':':
        return integerReply(line)
    case '$':
        return bulkReply(line, head)
    case '*':
        return multiBulkReply(line, head)
    }
    return nil, newError("Unknown response ")
}

func buildCommand(args ...string) []byte {
    cmd := bytes.NewBufferString(fmt.Sprintf("*%d\r\n", len(args)))
    for _, arg := range args {
        cmd.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
    }
    return cmd.Bytes()
}

func (c *Client) Send(cmd string, args ...string) (data interface{}, err os.Error) {
    if c.Addr == "" {
        c.Addr = defaultAddr
    }

    if c.pool == nil {
        c.pool = NewPool(c.Addr)
    }

    conn, err := c.pool.Pop()
    defer c.pool.Push(conn)

    if err != nil {
        return nil, err
    }

    cmds := append([]string{cmd}, args...)
    _, err = conn.Write(buildCommand(cmds...))
    if err != nil {
        return nil, err
    }

    data, err = Read(bufio.NewReader(conn))
    if err != nil {
        return nil, err
    }

    return
}
