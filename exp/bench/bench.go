package main

import (
    "flag"
    "fmt"
    "github.com/simonz05/godis/exp"
    "net"
    "os"
    "runtime"
    "runtime/pprof"
    "strings"
    "time"
)

var tests = make(map[string]func(*redis.Client, chan bool))
var C *int = flag.Int("c", 50, "concurrent requests")
var R *int = flag.Int("r", 4, "sample size")
var N *int = flag.Int("n", 10000, "number of request")
var P *int = flag.Int("p", 1, "pipeline requests")
var mock *bool = flag.Bool("mock", false, "run mock redis server")
var cpuprof *string = flag.String("cpuprof", "", "filename for cpuprof")

func init() {
    runtime.GOMAXPROCS(8)
}

func prints(t time.Duration) {
    fmt.Fprintf(os.Stdout, "    %.2f op/sec  real %.4fs\n", float64(*N)/t.Seconds(), t.Seconds())
}

func printsA(avg, tot time.Duration) {
    fmt.Fprintf(os.Stdout, "%.2f op/sec  real %.4fs  tot %.4fs\n", float64(*N)/avg.Seconds(), avg.Seconds(), tot.Seconds())
}

func BenchmarkMock(handle func(*redis.Client, chan bool)) time.Duration {
    ln, err := net.Listen("tcp", "127.0.0.1:6381")

    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }

    go MockRedis(ln)

    ch := make(chan bool)
    start := time.Now()

    for i := 0; i < *C; i++ {
        go handle(nil, ch)
    }

    for i := 0; i < *N; i++ {
        ch <- true
    }

    stop := time.Now().Sub(start)
    ln.Close()
    return stop
}

func BenchmarkRedis(handle func(*redis.Client, chan bool)) time.Duration {
    c := redis.NewClient("tcp:127.0.0.1:6379",13, "")

    //if _, err := c.Call("FLUSHDB"); err != nil {
    //    fmt.Fprintln(os.Stderr, err.Error())
    //    os.Exit(1)
    //}

    ch := make(chan bool)
    start := time.Now()

    for i := 0; i < *C; i++ {
        go handle(c, ch)
    }

    for i := 0; i < *N; i++ {
        ch <- true
    }

    return time.Now().Sub(start)
}

func run(name string) {
    var t, total time.Duration
    test, ok := tests[name]

    if !ok {
        fmt.Fprintf(os.Stderr, "test: `%s` does not exists\n", name)
        os.Exit(1)
    }

    fmt.Printf("%s:\n", strings.ToUpper(name))

    for i := 0; i < *R; i++ {
        if *mock {
            t = BenchmarkMock(test)
        } else {
            t = BenchmarkRedis(test)
        }

        total += t
        prints(t)
    }

    avg := time.Duration(total.Nanoseconds() / int64(*R))

    print("AVG ")
    printsA(avg, total)
    println()
}

func main() {
    flag.Parse()
    fmt.Printf("CONCURRENT: %d SAMPLES: %d REQUESTS: %d\n\n", *C, *R, *N)

    if *cpuprof != "" {
        file, err := os.OpenFile(*cpuprof, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            os.Exit(1)
        }

        defer file.Close()

        pprof.StartCPUProfile(file)
        defer pprof.StopCPUProfile()
    }

    redis.MaxConnections = *C

    for _, name := range flag.Args() {
        run(name)
    }

    println("ConnSum:", redis.ConnSum)

    stats := new(runtime.MemStats)
    runtime.ReadMemStats(stats)
}
