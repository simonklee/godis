package godis

import (
    "os"
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

func byteOrErr(res interface{}, err os.Error) ([]byte, os.Error) {
    v, ok := res.([]byte)

    if err != nil || !ok {
        return nil, err
    }

    return v, nil
}

func byteArrOrErr(res interface{}, err os.Error) ([][]byte, os.Error) {
    v, ok := res.([][]byte)

    if err != nil || !ok {
        return nil, err
    }

    return v, nil
}

func strToFaces(args []string) []interface{} {
	interfaces := make([]interface{}, len(args))

	for i, n := range args {
		interfaces[i] = n
	}

	return interfaces
}

func numsToFaces(args []int) []interface{} {
	interfaces := make([]interface{}, len(args))

	for i, n := range args {
		interfaces[i] = n
	}

	return interfaces
}

// generic

// Delete a key
func (c *Client) Del(keys ...string) (int64, os.Error) {
    return intOrErr(c.Send("DEL", strToFaces(keys)...))
}

// Determine if a key exists
func (c *Client) Exists(key string) (bool, os.Error) {
    return boolOrErr(c.Send("EXISTS", key))
}

// Set a key's time to live in seconds
func (c *Client) Expire(key string, seconds int64) (bool, os.Error) {
    return boolOrErr(c.Send("EXPIRE", key, seconds))
}

// Set the expiration for a key as a UNIX timestamp
func (c *Client) Expireat(key string, timestamp int64) (bool, os.Error) {
    return boolOrErr(c.Send("EXPIREAT", key, timestamp))
}

// Find all keys matching the given pattern
func (c *Client) Keys(pattern string) ([]string, os.Error) {
    return stringArrOrErr(c.Send("KEYS", pattern))
}

// Move a key to another database
func (c *Client) Move(key string, db int) (bool, os.Error) {
    return boolOrErr(c.Send("MOVE", key, db))
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
    res, err := c.Send("SORT", strToFaces(append([]string{key}, args...))...)

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
    return intOrErr(c.Send("DECRBY", key, decrement))
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
    return intOrErr(c.Send("GETBIT", key, offset))
}

// Get a substring of the string stored at a key
func (c *Client) Getrange(key string, start int, end int) (string, os.Error) {
    return stringOrErr(c.Send("GETRANGE", key, start, end))
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
    return intOrErr(c.Send("INCRBY", key, increment))
}

// Get the values of all the given keys
func (c *Client) Mget(keys ...string) ([]string, os.Error) {
    return stringArrOrErr(c.Send("MGET", strToFaces(keys)...))
}

// Set multiple keys to multiple values
func (c *Client) Mset(mapping map[string]string) os.Error {
    buf := make([]interface{}, len(mapping) * 2)
    n := 0

    for k, v := range mapping {
        buf[n], buf[n + 1] = k, v
        n += 2
    }

    _, err := c.Send("MSET", buf...)
    return err
}

// Set multiple keys to multiple values, only if none of the keys exist
func (c *Client) Msetnx(mapping map[string]string) (bool, os.Error) {
    buf := make([]interface{}, len(mapping) * 2)
    n := 0

    for k, v := range mapping {
        buf[n], buf[n + 1] = k, v
        n += 2
    }

    return boolOrErr(c.Send("MSETNX", buf...))
}

// Set the string value of a key
func (c *Client) Set(key string, value string) os.Error {
    return nilOrErr(c.Send("SET", key, value))
}

// Sets or clears the bit at offset in the string value stored at key
func (c *Client) Setbit(key string, offset int, value int) (int64, os.Error) {
    return intOrErr(c.Send("SETBIT", key, offset, value))
}

// Set the value and expiration of a key
func (c *Client) Setex(key string, seconds int64, value string) os.Error {
    return nilOrErr(c.Send("SETEX", key, seconds, value))
}

// Set the value of a key, only if the key does not exist
func (c *Client) Setnx(key string, value string) (bool, os.Error) {
    return boolOrErr(c.Send("SETNX", key, value))
}

// Overwrite part of a string at key starting at the specified offset
func (c *Client) Setrange(key string, offset int, value string) (int64, os.Error) {
    return intOrErr(c.Send("SETRANGE", key, offset, value))
}

// Get the length of the value stored in a key
func (c *Client) Strlen(key string) (int64, os.Error) {
    return intOrErr(c.Send("STRLEN", key))
}

// list

// Remove and get the first element in a list, or block until one is available
//func (c *Client) Blpop(key []string, timeout int64) [][]byte {
//
//}

// Remove and get the last element in a list, or block until one is available
//func (c *Client) Brpop(key []string, timeout int64) [][]byte {
//
//}

// Pop a value from a list, push it to another list and return it; or block until one is available
//func (c *Client) Brpoplpush(source string, destination string, timeout int64) []byte {
//
//}

// Get an element from a list by its index
func (c *Client) Lindex(key string, index int) ([]byte, os.Error) {
    return byteOrErr(c.Send("LINDEX", key, index))
}

// Insert an element before or after another element in a list
func (c *Client) Linsert(key, where string, pivot, value interface{}) (int64, os.Error) {
    return intOrErr(c.Send("LINSERT", key, where, pivot, value))
}

// Get the length of a list
func (c *Client) Llen(key string) (int64, os.Error) {
    return intOrErr(c.Send("LLEN", key))
}

// Remove and get the first element in a list
func (c *Client) Lpop(key string) ([]byte, os.Error) {
    return byteOrErr(c.Send("LPOP", key))
}

// Prepend a value to a list
func (c *Client) Lpush(key string, value interface{}) (int64, os.Error) {
    return intOrErr(c.Send("LPUSH", key, value))
}

// Prepend a value to a list, only if the list exists
func (c *Client) Lpushx(key string, value interface{}) (int64, os.Error) {
    return intOrErr(c.Send("LPUSHX", key, value))
}

// Get a range of elements from a list
func (c *Client) Lrange(key string, start, stop int) ([][]byte, os.Error) {
    return byteArrOrErr(c.Send("LRANGE", key, start, stop))
}

// Remove elements from a list
func (c *Client) Lrem(key string, count int, value string) (int64, os.Error) {
    return intOrErr(c.Send("LREM", key, count, value))
}

// Set the value of an element in a list by its index
func (c *Client) Lset(key string, index int, value string) os.Error {
    return nilOrErr(c.Send("LSET", key, index, value))
}

// Trim a list to the specified range
func (c *Client) Ltrim(key string, start int, stop int) os.Error {
    return nilOrErr(c.Send("LTRIM", key, start, stop))
}

// Remove and get the last element in a list
func (c *Client) Rpop(key string) ([]byte, os.Error) {
    return byteOrErr(c.Send("RPOP", key))
}

// Remove the last element in a list, append it to another list and return it
func (c *Client) Rpoplpush(source string, destination string) ([]byte, os.Error) {
    return byteOrErr(c.Send("RPOPLPUSH", source, destination))
}

// Append a value to a list
func (c *Client) Rpush(key string, value interface{}) (int64, os.Error) {
    return intOrErr(c.Send("RPUSH", key, value))
}

// Append a value to a list, only if the list exists
func (c *Client) Rpushx(key string, value interface{}) (int64, os.Error) {
    return intOrErr(c.Send("RPUSHX", key, value))
}
