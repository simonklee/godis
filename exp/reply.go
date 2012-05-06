package redis

import (
    "strconv"
    "strings"
)

type Elem []byte

type Reply struct {
    Err   error
    Elem  Elem
    Elems []*Reply
}

type Message struct {
    Channel string
    Elem    Elem
}

func (e Elem) Bytes() []byte {
    return []byte(e)
}

func (e Elem) String() string {
    return string([]byte(e))
}

func (e Elem) Bool() bool {
    v, _ := strconv.ParseBool(e.String())
    return v
}

func (e Elem) Int() int {
    v, _ := strconv.ParseInt(e.String(), 10, 0)
    return int(v)
}

func (e Elem) Int64() int64 {
    v, _ := strconv.ParseInt(e.String(), 10, 64)
    return v
}

func (e Elem) Float64() float64 {
    v, _ := strconv.ParseFloat(e.String(), 64)
    return v
}

func (r *Reply) Nil() bool {
    return r.Elems == nil && r.Elem == nil && r.Err == nil
}

func (r *Reply) Len() int {
    if r.Elems == nil {
        return 0
    }

    return len(r.Elems)
}

func (r *Reply) BytesArray() [][]byte {
    buf := make([][]byte, len(r.Elems))

    for i, v := range r.Elems {
        buf[i] = v.Elem
    }

    return buf
}

func (r *Reply) StringArray() []string {
    buf := make([]string, len(r.Elems))

    for i, v := range r.Elems {
        buf[i] = v.Elem.String()
    }

    return buf
}

func (r *Reply) IntArray() []int64 {
    buf := make([]int64, len(r.Elems))

    for i, v := range r.Elems {
        v, _ := strconv.ParseInt(v.Elem.String(), 10, 64)
        buf[i] = v
    }

    return buf
}

func (r *Reply) StringMap() map[string]string {
    arr := r.StringArray()
    n := len(arr)
    buf := make(map[string]string, n/2)

    if n%2 == 1 {
        return buf
    }

    for i := 0; i < n; i += 2 {
        buf[arr[i]] = arr[i+1]
    }

    return buf
}

func (r *Reply) Hash() map[string]Elem {
    l := r.Len()
    h := make(map[string]Elem, l/2)

    if l%2 == 1 {
        return h
    }

    var key string

    for i := 0; i < l; i += 2 {
        key = r.Elems[i].Elem.String()
        h[key] = r.Elems[i+1].Elem
    }

    return h
}

func (r *Reply) Message() *Message {
    if len(r.Elems) < 3 {
        return nil
    }

    typ := r.Elems[0].Elem.String()

    switch typ {
    case "message":
        return &Message{r.Elems[1].Elem.String(), r.Elems[2].Elem}
    case "pmessage":
        return &Message{r.Elems[2].Elem.String(), r.Elems[3].Elem}
    }

    if strings.HasSuffix(typ, "subscribe") {
        return nil
    }

    return nil
}
