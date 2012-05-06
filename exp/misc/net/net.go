package main

import (
    "flag"
    "fmt"
    "log"
    "math/rand"
    "net"
    "os"
    "runtime"
    "runtime/pprof"
    "strconv"
    "time"
)

var C *int = flag.Int("c", 1, "concurrent connections")
var R *int = flag.Int("r", 4, "sample size")
var N *int = flag.Int("n", 10000, "number of request")
var D *int = flag.Int("d", 512, "max data size")
var scale *bool = flag.Bool("scale", false, "scale read byffer dynamically")
var redis *bool = flag.Bool("redis", false, "use redis as server backend")
var cpuprof *string = flag.String("cpuprof", "", "filename for cpuprof")

var (
    maxIOBuf = uint16(1024)
    minIOBuf = uint16(8)
    data     [][]byte
)

func init() {
    runtime.GOMAXPROCS(8)
}

func prints(t time.Duration) {
    fmt.Fprintf(os.Stdout, "    %.2f op/sec  real %.4fs\n", float64(*N)/t.Seconds(), t.Seconds())
}

func printsA(avg, tot time.Duration) {
    fmt.Fprintf(os.Stdout, "%.2f op/sec  real %.4fs  tot %.4fs\n", float64(*N)/avg.Seconds(), avg.Seconds(), tot.Seconds())
}

func createDataTable() {
    data = make([][]byte, *D)
    for i := 0; i < *D; i++ {
        s := make([]byte, i)

        for j := range s {
            if j != i-1 {
                s[j] = 'a'
            } else {
                s[j] = '\n'
            }
        }

        data[i] = s
    }
}

func serve(ln net.Listener, open chan net.Conn) {
    cnt := 0
    defer ln.Close()

    for {
        conn, err := ln.Accept()
        cnt++

        if err != nil {
            //log.Println(err.Error())
            break
        }

        go handle(conn, cnt)
        open <- conn
    }
}

func handle(c net.Conn, nr int) {
    buf := make([]byte, 16)

    for {
        _, err := c.Read(buf)

        if err != nil {
            break
        }

        _, err = c.Write(data[rand.Int31n(int32(*D))])

        if err != nil {
            break
        }

        //s := string(buf)
        //log.Printf("#%d, nwrite: %d\n", nr, n)
    }
}

func round(v uint16) uint16 {
    v--
    v |= v >> 1
    v |= v >> 2
    v |= v >> 4
    v |= v >> 8
    v++
    return v
}

func max(v uint16) uint16 {
    if v > maxIOBuf {
        return maxIOBuf
    }

    return v
}

func min(v uint16) uint16 {
    if v < minIOBuf {
        return minIOBuf
    }

    return v
}

func client(done chan bool, netaddr string) {
    defer func() {
        done <- true
    }()

    conn, err := net.Dial("tcp", netaddr)

    if err != nil {
        log.Fatalln("dial error", err.Error())
    }

    defer conn.Close()

    //cmd := []byte("*2\r\n$3\r\nGET\r\n$3\r\n255\r\n")

    avg := uint16(maxIOBuf / 2)
    lastavg := avg
    buf := make([]byte, maxIOBuf)
    l := *N / *C

    for i := 0; i < l; i++ {
        n := strconv.Itoa(int(rand.Int31n(int32(*D))))
        s := fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%s\r\n", len(n), n)
        //println(s)

        if _, err := conn.Write([]byte(s)); err != nil {
            log.Println(err.Error())
            return
        }

        ok := false

        for !ok {
            nread, err := conn.Read(buf)

            if err != nil {
                log.Println(err.Error())
                return
            }

            avg += uint16(nread)

            for j := 0; j < nread; j++ {
                if buf[j] == '\n' {
                    //println(string(buf[:j]))
                    ok = true
                    break
                }
            }
        }

        if i%100 == 0 {
            avg /= 100
            avg = round(avg)
            avg = max(avg)

            if avg != lastavg {
                println(avg, avg, lastavg)

                if *scale {
                    buf = make([]byte, avg)
                }

                lastavg = avg
            }
        }
    }
}

func run() time.Duration {
    done := make(chan bool)
    open := make(chan net.Conn, 128)
    ln, err := net.Listen("tcp", "127.0.0.1:6381")

    if err != nil {
        log.Fatalln(err.Error())
    }

    start := time.Now()

    go serve(ln, open)
    go client(done, "127.0.0.1:6381")

    <-done
    stop := time.Now().Sub(start)

    c := <-open
    c.Close()
    ln.Close()
    return stop
}

func runRedis() time.Duration {
    done := make(chan bool)

    start := time.Now()
    for i := 0; i < *C; i++ {
        go client(done, "127.0.0.1:6379")
    }

    for i := 0; i < *C; i++ {
        <-done
    }

    return time.Now().Sub(start)
}

func main() {
    flag.Parse()
    log.Printf("REQUESTS: %d\n\n", *N)

    createDataTable()

    if *cpuprof != "" {
        file, err := os.OpenFile(*cpuprof, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

        if err != nil {
            log.Fatalln(err.Error())
        }

        defer file.Close()

        pprof.StartCPUProfile(file)
        defer pprof.StopCPUProfile()
    }

    var t, total time.Duration

    for i := 0; i < *R; i++ {
        if *redis {
            t = runRedis()
        } else {
            t = run()
        }

        total += t
        prints(t)
    }

    avg := time.Duration(total.Nanoseconds() / int64(*R))

    print("AVG ")
    printsA(avg, total)
    println()
}
