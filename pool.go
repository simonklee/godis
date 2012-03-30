package godis

var MaxConnections = 50

type connPool struct {
    free chan *Conn
}

func newConnPool() *connPool {
    p := connPool{make(chan *Conn, MaxConnections)}

    for i := 0; i < MaxConnections; i++ {
        p.free <- nil
    }

    return &p
}

func (p *connPool) pop() *Conn {
    return <-p.free
}

func (p *connPool) push(c *Conn) {
    p.free <- c
}
