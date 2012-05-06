package redis

var MaxConnections = 50

type connPool struct {
    free chan Connection
}

func newConnPool() *connPool {
    p := connPool{make(chan Connection, MaxConnections)}

    for i := 0; i < MaxConnections; i++ {
        p.free <- nil
    }

    return &p
}

func (p *connPool) pop() Connection {
    return <-p.free
}

func (p *connPool) push(c Connection) {
    p.free <- c
}
