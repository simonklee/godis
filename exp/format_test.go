package redis

import (
    "testing"
)

func formatTest(t *testing.T, exp string, a ...interface{}) {
    got := format(a...)

    if exp != string(got) {
        t.Errorf("format: expected %s got %s", exp, string(got))
    }
}

func TestFormat(t *testing.T) {
    formatTest(t, "*2\r\n$4\r\nPING\r\n$4\r\nPONG\r\n", "PING", "PONG")
    formatTest(t, "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", "SET", "foo", "bar")
    formatTest(t, "*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n", "GET", "foo")
}

func BenchmarkFormat(t *testing.B) {
    for i := 0; i < t.N; i++ {
        format("SET", "foo", "bar")
    }
}
