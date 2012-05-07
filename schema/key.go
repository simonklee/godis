package schema 

import (
    "fmt"
)

type Key struct {
    kind string
    id   int64
}

func (k *Key) String() string {
    return fmt.Sprintf("%s:%d", k.kind, k.id)
}

func (k *Key) Count() string {
    return fmt.Sprintf("%s:count", k.kind)
}

func (k *Key) Unique(field, value string) string {
    return fmt.Sprintf("%s:unique:%s:%s", k.kind, field, value)
}

func (k *Key) Index(field, value string) string {
    return fmt.Sprintf("%s:index:%s:%s", k.kind, field, value)
}

func (k *Key) Id() int64 {
    return k.id
}

func NewKey(kind string, id int64) *Key {
    return &Key{kind, id}
}
