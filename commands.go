package godis

import (
    "os"
    "strconv"
)

// helpers for handling common return values
func intOrErr(res interface{}, err os.Error) (int64, os.Error) {
    if err != nil {
        return -1, err
    }

    return res.(int64), nil
}

func boolOrErr(res interface{}, err os.Error) (bool, os.Error) {
    if err != nil {
        return false, err
    }

    return res.(int64) == 1, nil
}

func nilOrErr(res interface{}, err os.Error) os.Error {
    return err
}

func stringOrErr(res interface{}, err os.Error) (string, os.Error) {
    if err != nil {
        return "", err
    }

    switch v := res.(type) {
    case string:
        return v, nil
    case []byte:
        return string(v), nil
    }

    return "", nil
}

func stringArrOrErr(res interface{}, err os.Error) ([]string, os.Error) {
    v, ok := res.([][]byte)

    if err != nil || !ok {
        return nil, err
    }

    out := make([]string, len(v))

    for i, k := range v {
        out[i] = string(k)
    }

    return out, nil
}

// generic

// Delete a key
func (c *Client) Del(keys ...string) (int64, os.Error) {
    return intOrErr(c.Send("DEL", keys...))
}

// Determine if a key exists
func (c *Client) Exists(key string) (bool, os.Error) {
    return boolOrErr(c.Send("EXISTS", key))
}

// Set a key's time to live in seconds
func (c *Client) Expire(key string, seconds int64) (bool, os.Error) {
    return boolOrErr(c.Send("EXPIRE", key, strconv.Itoa64(seconds)))
}

// Set the expiration for a key as a UNIX timestamp
func (c *Client) Expireat(key string, timestamp int64) (bool, os.Error) {
    return boolOrErr(c.Send("EXPIREAT", key, strconv.Itoa64(timestamp)))
}

// Find all keys matching the given pattern
func (c *Client) Keys(pattern string) ([]string, os.Error) {
    return stringArrOrErr(c.Send("KEYS", pattern))
}

// Move a key to another database
func (c *Client) Move(key string, db int) (bool, os.Error) {
    return boolOrErr(c.Send("MOVE", key, strconv.Itoa(db)))
}

// Remove the expiration from a key
func (c *Client) Persist(key string) (bool, os.Error) {
    return boolOrErr(c.Send("PERSIST", key))
}

// Return a random key from the keyspace
func (c *Client) Randomkey() (string, os.Error) {
    return stringOrErr(c.Send("RANDOMKEY")) 
}

// Rename a key
func (c *Client) Rename(key string, newkey string) os.Error {
    return nilOrErr(c.Send("RENAME", key, newkey))
}

// Rename a key, only if the new key does not exist
func (c *Client) Renamenx(key string, newkey string) (bool, os.Error) {
    return boolOrErr(c.Send("RENAMENX", key, newkey))
}

// Sort the elements in a list, set or sorted set
func (c *Client) Sort(key string, args...string) ([][]byte, os.Error) {
    res, err := c.Send("SORT", append([]string{key}, args...)...)

    if err != nil {
        return nil, err
    }

    v, ok := res.([][]byte)

    if !ok {
        return nil, err
    }
    return v, nil
    ///out := make([]byte, len(v))

    ///for i, k := range v {
    ///    out[i] = string(k)
    ///}

    ///return out, nil
}

// Get the time to live for a key
func (c *Client) Ttl(key string) (int64, os.Error) {
    return intOrErr(c.Send("TTL", key))
}

// Determine the type stored at key
func (c *Client) Type(key string) (string, os.Error) {
    return stringOrErr(c.Send("TYPE", key))
}

// strings

// Append a value to a key
func (c *Client) Append(key string, value string) (int64, os.Error) {
    return intOrErr(c.Send("APPEND", key, value))
}

// Decrement the integer value of a key by one
func (c *Client) Decr(key string) (int64, os.Error) {
    return intOrErr(c.Send("DECR", key))
}

// Decrement the integer value of a key by the given number
func (c *Client) Decrby(key string, decrement int64) (int64, os.Error) {
    return intOrErr(c.Send("DECRBY", key, strconv.Itoa64(decrement)))
}

// Get the value of a key
func (c *Client) Get(key string) (string, os.Error) {
    res, err := c.Send("GET", key)

    if err == nil && res == nil {
        err = newError("key `%s` does not exist", key)
    }

    return stringOrErr(res, err)
}

// Returns the bit value at offset in the string value stored at key
func (c *Client) Getbit(key string, offset int) (int64, os.Error) {
    return intOrErr(c.Send("GETBIT", key, strconv.Itoa(offset)))
}

// Get a substring of the string stored at a key
func (c *Client) Getrange(key string, start int, end int) (string, os.Error) {
    return stringOrErr(c.Send("GETBIT", key, strconv.Itoa(start), strconv.Itoa(end)))
}

// Set the string value of a key and return its old value
func (c *Client) Getset(key string, value string) (string, os.Error) {
    return stringOrErr(c.Send("GETSET", key, value))
}

// Increment the integer value of a key by one
func (c *Client) Incr(key string) (int64, os.Error) {
    return intOrErr(c.Send("INCR", key))
}

// Increment the integer value of a key by the given number
func (c *Client) Incrby(key string, increment int64) (int64, os.Error) {
    return intOrErr(c.Send("INCRBY", key, strconv.Itoa64(increment)))
}

// Get the values of all the given keys
func (c *Client) Mget(keys ...string) ([]string, os.Error) {
    return stringArrOrErr(c.Send("MGET", keys...))
}

// Set multiple keys to multiple values
func (c *Client) Mset(mapping map[string][]byte) os.Error {
    buf := make([]string, len(mapping) * 2)
    n := 0

    for k, v := range mapping {
        buf[n] = k
        buf[n + 1] = string(v)
        n += 2
    }

    _, err := c.Send("MSET", buf...)
    return err
}

// Set multiple keys to multiple values, only if none of the keys exist
func (c *Client) Msetnx(mapping map[string][]byte) (bool, os.Error) {
    buf := make([]string, len(mapping) * 2)
    n := 0

    for k, v := range mapping {
        buf[n] = k
        buf[n + 1] = string(v)
        n += 2
    }

    return boolOrErr(c.Send("MSETNX", buf...))
}

// Set the string value of a key
func (c *Client) Set(key string, value string) os.Error {
    _, err := c.Send("SET", key, value)
    return err
}

// Sets or clears the bit at offset in the string value stored at key
func (c *Client) Setbit(key string, offset int, value string) (int64, os.Error) {
    return intOrErr(c.Send("SETBIT", key, strconv.Itoa(offset), value))
}

// Set the value and expiration of a key
func (c *Client) Setex(key string, seconds int64, value string) os.Error {
    return nilOrErr(c.Send("SET", key, strconv.Itoa64(seconds), value))
}

// Set the value of a key, only if the key does not exist
func (c *Client) Setnx(key string, value string) (bool, os.Error) {
    return boolOrErr(c.Send("SETNX", key, value))
}

// Overwrite part of a string at key starting at the specified offset
func (c *Client) Setrange(key string, offset int, value string) (int64, os.Error) {
    return intOrErr(c.Send("SETRANGE", key, strconv.Itoa(offset), value))
}

// Get the length of the value stored in a key
func (c *Client) Strlen(key string) (int64, os.Error) {
    return intOrErr(c.Send("STRLEN", key))
}
 
