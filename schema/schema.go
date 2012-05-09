package schema

import (
    "fmt"
    "reflect"

    "code.google.com/p/tcgl/monitoring"
    "github.com/simonz05/godis/exp"
)

type UserError string 

func (e UserError) Error() string {
    return string(e)
}

type InternalError string 

func (e InternalError) Error() string {
    return string(e)
}

func newUniqueError(field string, value interface{}) UserError {
    return UserError(fmt.Sprintf("Expected unique `%s`:`%v`", field, value))
}

var (
    nilError     = UserError("Key requested returned a nil reply")
    typeError    = InternalError("Invalid type, expected pointer to struct")
    idFieldError = InternalError("Expected Id field of type int64 on struct")
)

func IsUserError(e error) bool {
    _, ok := e.(UserError)
    return ok
}

func IsInternalError(e error) bool {
    _, ok := e.(InternalError)
    return ok
}

/* logic:
   input: key, *struct

   a) Check if key has Id; yes/no
       no) INCR k:count
           set value as Id in struct

   b) Parse struct; returns [field, value, field, value, ...], [hashField, ...].
      the hashField slice contains any optional properties which should be applied
      to before or after setting the struct in Redis (uniquity, index).

   c) Check unique fields and set unique fields

   d) Add the actual data from struct to a redis hash

   e) Add optional indexes

   f) if c, d, e did raise an error we need to cleanup the things 
      we changed in the database.
*/
func Put(db *redis.Client, k *Key, s interface{}) (*Key, error) {
    mon := monitoring.BeginMeasuring("database:put")
    defer mon.EndMeasuring()
    prep := new(prepare)
    prep.key = k

    e := setId(db, s, prep)

    if e != nil {
        return nil, e
    }

    e = parseStruct(db, s, prep)

    if e != nil {
        return nil, e
    }

    for _, o := range prep.unique {
        uk := k.Unique(o.name, fmt.Sprintf("%v", *o.value))
        var reply *redis.Reply
        // TODO: Watch key
        reply, e = db.Call("GET", uk)

        if e != nil {
            goto Error
        }

        if !reply.Nil() && reply.Elem.Int64() != k.Id() {
            e = newUniqueError(o.name, *o.value)
            goto Error
        }

        reply, e = db.Call("SET", uk, k.Id())

        if e != nil {
            goto Error
        }

        prep.dirty = append(prep.dirty, uk)
    }

    _, e = db.Call(append([]interface{}{"HMSET", k.String()}, prep.args...)...)

    for _, o := range prep.index {
        ik := fmt.Sprintf("%v", *o.value)
        _, e = db.Call("SET", k.Index(o.name, ik), k.Id())

        if e != nil {
            goto Error
        }

        prep.dirty = append(prep.dirty, ik)
    }

    return k, e
Error:
    cleanup(db, prep)
    return nil, e
}

func Get(db *redis.Client, k *Key, s interface{}) error {
    mon := monitoring.BeginMeasuring("database:get")
    reply, e := db.Call("HGETALL", k.String())

    if e != nil {
        return e
    }

    if reply.Len() == 0 {
        return nilError
    }

    e = inflate(db, s, reply.Hash())
    mon.EndMeasuring()
    return e
}

func cleanup(db *redis.Client, prep *prepare) error {
    p := db.AsyncClient()

    for _, k := range prep.dirty {
        p.Call("DEL", k)
    }

    if prep.isNew {
        p.Call("DEL", prep.key.String())
    }

    _, e := p.ReadAll()
    return e
}

func setId(db *redis.Client, s interface{}, prep *prepare) error {
    if prep.key.Id() != 0 {
        return nil
    }

    prep.isNew = true
    reply, e := db.Call("INCR", prep.key.Count())

    if e != nil {
        return e
    }

    prep.key.id = reply.Elem.Int64()

    if e != nil {
        return e
    }
    v := reflect.ValueOf(s)

    if v.Kind() != reflect.Ptr {
        return typeError
    }

    v = v.Elem()

    if v.Kind() != reflect.Struct {
        return typeError
    }

    idField := v.FieldByName("Id")

    if !idField.IsValid() || !idField.CanSet() || idField.Kind() != reflect.Int64 {
        return idFieldError
    }

    idField.SetInt(prep.key.Id())
    return nil
}

type hashField struct {
    name  string
    value *interface{}
}

type prepare struct {
    key    *Key
    unique []*hashField
    index  []*hashField
    args   []interface{}
    dirty  []string
    isNew  bool
}

// parseStruct takes a pointer to a struct or a struct.
// We use the struct to fill an array of key:value from the struct.
func parseStruct(db *redis.Client, s interface{}, prep *prepare) error {
    v := reflect.ValueOf(s)

    if v.Kind() != reflect.Ptr {
        return typeError
    }

    v = v.Elem()

    if v.Kind() != reflect.Struct {
        return typeError
    }

    t := v.Type()
    n := v.NumField()
    args := make([]interface{}, 0, n*2)
    argsLen := cap(args)
    prep.unique = make([]*hashField, 0)
    prep.index = make([]*hashField, 0)

    for i := 0; i < n; i++ {
        fieldType := t.Field(i)

        if len(fieldType.PkgPath) > 0 || fieldType.Anonymous {
            fmt.Printf("unexported field `%s`\n", fieldType.Name)
            argsLen -= 2
            continue
        }

        name, opt := parseTag(fieldType.Tag.Get("redis"))

        if name == "" {
            name = fieldType.Name
        }

        fieldValue := v.Field(i).Interface()
        args = append(args, name, fieldValue)

        unique := opt.Contains("unique")
        index := opt.Contains("index")

        if index || unique {
            hf := &hashField{
                name:  name,
                value: &fieldValue,
            }

            if index {
                prep.index = append(prep.index, hf)
            }

            if unique {
                prep.unique = append(prep.unique, hf)
            }
        }
    }

    prep.args = args[:argsLen]
    return nil
}

// inflate takes a pointer to a struct as dst and a map with 
// values as src. The struct is then filled with the values 
// from the map. 
func inflate(db *redis.Client, dst interface{}, src map[string]redis.Elem) error {
    v := reflect.ValueOf(dst)

    if v.Kind() != reflect.Ptr {
        return typeError
    }

    v = v.Elem()

    if v.Kind() != reflect.Struct {
        return typeError
    }

    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        fieldValue := v.Field(i)
        fieldType := t.Field(i)

        if !fieldValue.CanSet() {
            fmt.Println("Field not setable", fieldType.Name)
            continue
        }

        name, _ := parseTag(fieldType.Tag.Get("json"))

        if name == "" {
            name = fieldType.Name
        }

        value, ok := src[name]

        if !ok {
            fmt.Println("Value for field not in map:", name)
            continue
        }

        switch fieldValue.Kind() {
        case reflect.Int:
            fieldValue.SetInt(int64(value.Int()))
        case reflect.Int64:
            fieldValue.SetInt(value.Int64())
        case reflect.Float64:
            fieldValue.SetFloat(value.Float64())
        case reflect.Bool:
            fieldValue.SetBool(value.Bool())
        case reflect.String:
            fieldValue.SetString(value.String())
        default:
            panic("type not supported: " + fieldValue.Type().String())
        }
    }

    return nil
}
