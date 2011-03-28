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
func (c *Client) Expire(key string, seconds int) (bool, os.Error) {
    return boolOrErr(c.Send("EXPIRE", key, strconv.Itoa(seconds)))
}

// Set the expiration for a key as a UNIX timestamp
func (c *Client) Expireat(key string, timestamp int64) (bool, os.Error) {
    return boolOrErr(c.Send("EXPIREAT", key, strconv.Itoa64(timestamp)))
}

// Find all keys matching the given pattern
func (c *Client) Keys(pattern string) ([]string, os.Error) {
    res, err := c.Send("KEYS", pattern)
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
    _, err := c.Send("RENAME", key, newkey)
    return err
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

// Get the value of a key
func (c *Client) Get(key string) (string, os.Error) {
    res, err := c.Send("GET", key)

    if err == nil && res == nil {
        err = newError("key `%s` does not exist", key)
    }

    return stringOrErr(res, err)
}

// Set the string value of a key
func (c *Client) Set(key string, value string) os.Error {
    _, err := c.Send("SET", key, value)
    return err
}
