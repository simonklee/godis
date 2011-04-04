package godis

import (
    "os"
)

// helpers for handling common return values
func (r *Reply) intOrErr() (int64, os.Error) {
    if r.Err != nil {
        return -1, r.Err
    }

    return r.Elem.Int64(), nil
}

func (r *Reply) boolOrErr() (bool, os.Error) {
    if r.Err != nil {
        return false, r.Err
    }

    return r.Elem.Int64() == 1, nil
}

func (r *Reply) nilOrErr() os.Error {
    return r.Err
}

func (r *Reply) stringOrErr() (string, os.Error) {
    if r.Err != nil {
        return "", r.Err
    }

    return r.Elem.String(), nil
}

func (r *Reply) stringArrOrErr() ([]string, os.Error) {
    if r.Err != nil {
        return nil, r.Err
    }

    return r.StringArray(), nil
}

func (r *Reply) elemOrErr() (Elem, os.Error) {
    if r.Err != nil {
        return nil, r.Err
    }

    return r.Elem, nil
}

func (r *Reply) replyOrErr() (*Reply, os.Error) {
    if r.Err != nil {
        return nil, r.Err
    }

    return r, nil
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

func smapToArr(mapping map[string] string) []interface{} {
    buf := make([]interface{}, len(mapping)*2)
    n := 0

    for k, v := range mapping {
        buf[n], buf[n+1] = k, v
        n += 2
    }

    return buf
}

// generic

// Delete a key
func (c *Client) Del(keys ...string) (int64, os.Error) {
    return Send(c, "DEL", strToFaces(keys)...).intOrErr()
}

// Determine if a key exists
func (c *Client) Exists(key string) (bool, os.Error) {
    return Send(c, "EXISTS", key).boolOrErr()
}

// Set a key's time to live in seconds
func (c *Client) Expire(key string, seconds int64) (bool, os.Error) {
    return Send(c, "EXPIRE", key, seconds).boolOrErr()
}

// Set the expiration for a key as a UNIX timestamp
func (c *Client) Expireat(key string, timestamp int64) (bool, os.Error) {
    return Send(c, "EXPIREAT", key, timestamp).boolOrErr()
}

// Find all keys matching the given pattern
func (c *Client) Keys(pattern string) ([]string, os.Error) {
    return Send(c, "KEYS", pattern).stringArrOrErr()
}

// Move a key to another database
func (c *Client) Move(key string, db int) (bool, os.Error) {
    return Send(c, "MOVE", key, db).boolOrErr()
}

// Remove the expiration from a key
func (c *Client) Persist(key string) (bool, os.Error) {
    return Send(c, "PERSIST", key).boolOrErr()
}

// Return a random key from the keyspace
func (c *Client) Randomkey() (string, os.Error) {
    return Send(c, "RANDOMKEY").stringOrErr()
}

// Rename a key
func (c *Client) Rename(key string, newkey string) os.Error {
    return Send(c, "RENAME", key, newkey).nilOrErr()
}

// Rename a key, only if the new key does not exist
func (c *Client) Renamenx(key string, newkey string) (bool, os.Error) {
    return Send(c, "RENAMENX", key, newkey).boolOrErr()
}

// Sort the elements in a list, set or sorted set
func (c *Client) Sort(key string, args ...string) (*Reply, os.Error) {
    a := strToFaces(append([]string{key}, args...))
    return Send(c, "SORT", a...).replyOrErr()
    ///out := make([]byte, len(v))

    ///for i, k := range v {
    ///    out[i] = string(k)
    ///}

    ///return out, nil
}

// Get the time to live for a key
func (c *Client) Ttl(key string) (int64, os.Error) {
    return Send(c, "TTL", key).intOrErr()
}

// Determine the type stored at key
func (c *Client) Type(key string) (string, os.Error) {
    return Send(c, "TYPE", key).stringOrErr()
}

// strings

// Append a value to a key
func (c *Client) Append(key string, value interface{}) (int64, os.Error) {
    return Send(c, "APPEND", key, value).intOrErr()
}

// Decrement the integer value of a key by one
func (c *Client) Decr(key string) (int64, os.Error) {
    return Send(c, "DECR", key).intOrErr()
}

// Decrement the integer value of a key by the given number
func (c *Client) Decrby(key string, decrement int64) (int64, os.Error) {
    return Send(c, "DECRBY", key, decrement).intOrErr()
}

// Get the value of a key
func (c *Client) Get(key string) (string, os.Error) {
    r := Send(c, "GET", key)

    if r.Err == nil && r.Elem == nil {
        r.Err = os.NewError("key `" + key + "`does not exist")
    }

    return r.stringOrErr()
}

// Returns the bit value at offset in the string value stored at key
func (c *Client) Getbit(key string, offset int) (int64, os.Error) {
    return Send(c, "GETBIT", key, offset).intOrErr()
}

// Get a substring of the string stored at a key
func (c *Client) Getrange(key string, start int, end int) (string, os.Error) {
    return Send(c, "GETRANGE", key, start, end).stringOrErr()
}

// Set the string value of a key and return its old value
func (c *Client) Getset(key string, value interface{}) (string, os.Error) {
    return Send(c, "GETSET", key, value).stringOrErr()
}

// Increment the integer value of a key by one
func (c *Client) Incr(key string) (int64, os.Error) {
    return Send(c, "INCR", key).intOrErr()
}

// Increment the integer value of a key by the given number
func (c *Client) Incrby(key string, increment int64) (int64, os.Error) {
    return Send(c, "INCRBY", key, increment).intOrErr()
}

// Get the values of all the given keys
func (c *Client) Mget(keys ...string) ([]string, os.Error) {
    return Send(c, "MGET", strToFaces(keys)...).stringArrOrErr()
}

// Set multiple keys to multiple values
func (c *Client) Mset(mapping map[string]string) os.Error {
    return Send(c, "MSET", smapToArr(mapping)...).nilOrErr()
}

// Set multiple keys to multiple values, only if none of the keys exist
func (c *Client) Msetnx(mapping map[string]string) (bool, os.Error) {
    return Send(c, "MSETNX", smapToArr(mapping)...).boolOrErr()
}

// Set the string value of a key
func (c *Client) Set(key string, value interface{}) os.Error {
    return Send(c, "SET", key, value).nilOrErr()
}

// Sets or clears the bit at offset in the string value stored at key
func (c *Client) Setbit(key string, offset int, value int) (int64, os.Error) {
    return Send(c, "SETBIT", key, offset, value).intOrErr()
}

// Set the value and expiration of a key
func (c *Client) Setex(key string, seconds int64, value interface{}) os.Error {
    return Send(c, "SETEX", key, seconds, value).nilOrErr()
}

// Set the value of a key, only if the key does not exist
func (c *Client) Setnx(key string, value interface{}) (bool, os.Error) {
    return Send(c, "SETNX", key, value).boolOrErr()
}

// Overwrite part of a string at key starting at the specified offset
func (c *Client) Setrange(key string, offset int, value interface{}) (int64, os.Error) {
    return Send(c, "SETRANGE", key, offset, value).intOrErr()
}

// Get the length of the value stored in a key
func (c *Client) Strlen(key string) (int64, os.Error) {
    return Send(c, "STRLEN", key).intOrErr()
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
func (c *Client) Lindex(key string, index int) (Elem, os.Error) {
    return Send(c, "LINDEX", key, index).elemOrErr()
}

// Insert an element before or after another element in a list
func (c *Client) Linsert(key, where string, pivot, value interface{}) (int64, os.Error) {
    return Send(c, "LINSERT", key, where, pivot, value).intOrErr()
}

// Get the length of a list
func (c *Client) Llen(key string) (int64, os.Error) {
    return Send(c, "LLEN", key).intOrErr()
}

// Remove and get the first element in a list
func (c *Client) Lpop(key string) (Elem, os.Error) {
    return Send(c, "LPOP", key).elemOrErr()
}

// Prepend a value to a list
func (c *Client) Lpush(key string, value interface{}) (int64, os.Error) {
    return Send(c, "LPUSH", key, value).intOrErr()
}

// Prepend a value to a list, only if the list exists
func (c *Client) Lpushx(key string, value interface{}) (int64, os.Error) {
    return Send(c, "LPUSHX", key, value).intOrErr()
}

// Get a range of elements from a list
func (c *Client) Lrange(key string, start, stop int) (*Reply, os.Error) {
    return Send(c, "LRANGE", key, start, stop).replyOrErr()
}

// Remove elements from a list
func (c *Client) Lrem(key string, count int, value string) (int64, os.Error) {
    return Send(c, "LREM", key, count, value).intOrErr()
}

// Set the value of an element in a list by its index
func (c *Client) Lset(key string, index int, value string) os.Error {
    return Send(c, "LSET", key, index, value).nilOrErr()
}

// Trim a list to the specified range
func (c *Client) Ltrim(key string, start int, stop int) os.Error {
    return Send(c, "LTRIM", key, start, stop).nilOrErr()
}

// Remove and get the last element in a list
func (c *Client) Rpop(key string) (Elem, os.Error) {
    return Send(c, "RPOP", key).elemOrErr()
}

// Remove the last element in a list, append it to another list and return it
func (c *Client) Rpoplpush(source string, destination string) (Elem, os.Error) {
    return Send(c, "RPOPLPUSH", source, destination).elemOrErr()
}

// Append a value to a list
func (c *Client) Rpush(key string, value interface{}) (int64, os.Error) {
    return Send(c, "RPUSH", key, value).intOrErr()
}

// Append a value to a list, only if the list exists
func (c *Client) Rpushx(key string, value interface{}) (int64, os.Error) {
    return Send(c, "RPUSHX", key, value).intOrErr()
}

// hash

// Delete a hash field
func (c *Client) Hdel(key string, field string) (bool, os.Error) {
    return Send(c, "HDEL", key, field).boolOrErr()
}

// Determine if a hash field exists
func (c *Client) Hexists(key string, field string) (bool, os.Error) {
    return Send(c, "HEXISTS", key, field).boolOrErr()
}

// Get the value of a hash field
func (c *Client) Hget(key string, field string) (Elem, os.Error) {
    return Send(c, "HGET", key, field).elemOrErr()
}

// Get all the fields and values in a hash
func (c *Client) Hgetall(key string) (*Reply, os.Error) {
    return Send(c, "HGETALL", key).replyOrErr()
}

// Increment the integer value of a hash field by the given number
func (c *Client) Hincrby(key string, field string, increment int64) (int64, os.Error) {
    return Send(c, "HINCRBY", key, field, increment).intOrErr()
}

// Get all the fields in a hash
func (c *Client) Hkeys(key string) ([]string, os.Error)  {
    return Send(c, "HKEYS", key).stringArrOrErr()
}

// Get the number of fields in a hash
func (c *Client) Hlen(key string) (int64, os.Error) {
    return Send(c, "HLEN", key).intOrErr()
}

// Get the values of all the given hash fields
func (c *Client) Hmget(key string, fields []string) (*Reply, os.Error) {
    a := strToFaces(append([]string{key}, fields...))
    return Send(c, "HMGET", a...).replyOrErr()
}

// Set multiple hash fields to multiple values
func (c *Client) Hmset(key string, mapping map[string] interface{}) os.Error {
    buf := make([]interface{}, len(mapping)*2 + 1)
    buf[0] = key
    n := 1
    
    for k, v := range mapping {
        buf[n], buf[n+1] = k, v
        n += 2
    }

    return Send(c, "HMSET",  buf...).nilOrErr()
}

// Set the string value of a hash field
func (c *Client) Hset(key string, field string, value interface{}) (bool, os.Error) {
    return Send(c, "HSET", key, field, value).boolOrErr()
}

// Set the value of a hash field, only if the field does not exist
func (c *Client) Hsetnx(key string, field string, value interface{}) (int64, os.Error) {
    return Send(c, "HSETNX", key, field, value).intOrErr()
}

// Get all the values in a hash
func (c *Client) Hvals(key string) (*Reply, os.Error) {
    return Send(c, "HVALS", key).replyOrErr()
}

// set

// Add a member to a set
func (c *Client) Sadd(key string, member interface{}) (bool, os.Error) {
    return Send(c, "SADD", key, member).boolOrErr()
}

// Get the number of members in a set
func (c *Client) Scard(key string) (int64, os.Error) {
    return Send(c, "SCARD", key).intOrErr()
}

// Subtract multiple sets
func (c *Client) Sdiff(keys ...string) (*Reply, os.Error) {
    return Send(c, "SDIFF", strToFaces(keys)...).replyOrErr()
}

// Subtract multiple sets and store the resulting set in a key
func (c *Client) Sdiffstore(destination string, keys ...string) (int64, os.Error) {
    b := append([]string{destination}, keys...)
    return Send(c, "SDIFFSTORE", strToFaces(b)...).intOrErr()
}

// Intersect multiple sets
func (c *Client) Sinter(keys ...string) (*Reply, os.Error) {
    return Send(c, "SINTER", strToFaces(keys)...).replyOrErr()
}

// Intersect multiple sets and store the resulting set in a key
func (c *Client) Sinterstore(destination string, keys ...string) (int64, os.Error) {
    b := append([]string{destination}, keys...)
    return Send(c, "SINTERSTORE", strToFaces(b)...).intOrErr()
}

// Determine if a given value is a member of a set
func (c *Client) Sismember(key string, member interface{}) (bool, os.Error) {
    return Send(c, "SISMEMBER", key, member).boolOrErr()
}

// Get all the members in a set
func (c *Client) Smembers(key string) (*Reply, os.Error) {
    return Send(c, "SMEMBERS", key).replyOrErr()
}

// Move a member from one set to another
func (c *Client) Smove(source string, destination string, member interface{}) (bool, os.Error) {
    return Send(c, "SMOVE", source, destination, member).boolOrErr()
}

// Remove and return a random member from a set
func (c *Client) Spop(key string) (Elem, os.Error) {
    return Send(c, "SPOP", key).elemOrErr()
}

// Get a random member from a set
func (c *Client) Srandmember(key string) (Elem, os.Error) {
    return Send(c, "SRANDMEMBER", key).elemOrErr()
}

// Remove a member from a set
func (c *Client) Srem(key string, member interface{}) (bool, os.Error) {
    return Send(c, "SREM", key, member).boolOrErr()
}

// Add multiple sets
func (c *Client) Sunion(key ...string) (*Reply, os.Error) {
    return Send(c, "SUNION", strToFaces(key)...).replyOrErr()
}

// Add multiple sets and store the resulting set in a key
func (c *Client) Sunionstore(destination string, keys ...string) (int64, os.Error) {
    b := append([]string{destination}, keys...)
    return Send(c, "SUNIONSTORE", strToFaces(b)...).intOrErr()
}

// sorted_set

// Add a member to a sorted set, or update its score if it already exists
func (c *Client) Zadd(key string, score float64, member interface{}) (int64, os.Error) {
    return Send(c, "ZADD").intOrErr()
}

// Get the number of members in a sorted set
func (c *Client) Zcard(key string) (int64, os.Error) {
    return Send(c, "ZCARD").intOrErr()
}

// Count the members in a sorted set with scores within the given values
func (c *Client) Zcount(key string, min float64, max float64) (int64, os.Error) {
    return Send(c, "ZCOUNT").intOrErr()
}

// Increment the score of a member in a sorted set
func (c *Client) Zincrby(key string, increment int, member string) ([]byte, os.Error) {
    return Send(c, "ZINCRBY").elemOrErr()
}

// Intersect multiple sorted sets and store the resulting sorted set in a new key
func (c *Client) Zinterstore(destination string, numkeys int, key []string) (int64, os.Error) {
    return Send(c, "ZINTERSTORE").intOrErr()
}

// Return a range of members in a sorted set, by index
func (c *Client) Zrange(key string, start int, stop int) (*Reply, os.Error) {
    return Send(c, "ZRANGE").replyOrErr()
}

// Return a range of members in a sorted set, by score
func (c *Client) Zrangebyscore(key string, min float64, max float64) (*Reply, os.Error) {
    return Send(c, "ZRANGEBYSCORE").replyOrErr()
}

// Determine the index of a member in a sorted set
func (c *Client) Zrank(key string, member string) (int64, os.Error) {
    return Send(c, "ZRANK").intOrErr()
}

// Remove a member from a sorted set
func (c *Client) Zrem(key string, member string) (int64, os.Error) {
    return Send(c, "ZREM").intOrErr()
}

// Remove all members in a sorted set within the given indexes
func (c *Client) Zremrangebyrank(key string, start int, stop int) (int64, os.Error) {
    return Send(c, "ZREMRANGEBYRANK").intOrErr()
}

// Remove all members in a sorted set within the given scores
func (c *Client) Zremrangebyscore(key string, min float64, max float64) (int64, os.Error) {
    return Send(c, "ZREMRANGEBYSCORE").intOrErr()
}

// Return a range of members in a sorted set, by index, with scores ordered from high to low
func (c *Client) Zrevrange(key string, start int, stop int) (*Reply, os.Error) {
    return Send(c, "ZREVRANGE").replyOrErr()
}

// Return a range of members in a sorted set, by score, with scores ordered from high to low
func (c *Client) Zrevrangebyscore(key string, max float64, min float64) (*Reply, os.Error) {
    return Send(c, "ZREVRANGEBYSCORE").replyOrErr()
}

// Determine the index of a member in a sorted set, with scores ordered from high to low
func (c *Client) Zrevrank(key string, member string) (int64, os.Error) {
    return Send(c, "ZREVRANK").intOrErr()
}

// Get the score associated with the given member in a sorted set
func (c *Client) Zscore(key string, member string) ([]byte, os.Error) {
    return Send(c, "ZSCORE").elemOrErr()
}

// Add multiple sorted sets and store the resulting sorted set in a new key
func (c *Client) Zunionstore(destination string, numkeys int, key []string) (int64, os.Error) {
    return Send(c, "ZUNIONSTORE").intOrErr()
}
