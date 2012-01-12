package godis

import (
    "errors"
    "strconv"
)

// helpers for handling common return values
func (r *Reply) intOrErr() (int64, error) {
    if r.Err != nil {
        return -1, r.Err
    }

    return r.Elem.Int64(), nil
}

// helpers for handling common return values
func (r *Reply) floatOrErr() (float64, error) {
    if r.Err != nil {
        return -1.0, r.Err
    }

    return r.Elem.Float64(), nil
}

func (r *Reply) boolOrErr() (bool, error) {
    if r.Err != nil {
        return false, r.Err
    }

    return r.Elem.Int64() == 1, nil
}

func (r *Reply) nilOrErr() error {
    return r.Err
}

func (r *Reply) stringOrErr() (string, error) {
    if r.Err != nil {
        return "", r.Err
    }

    return r.Elem.String(), nil
}

func (r *Reply) stringArrOrErr() ([]string, error) {
    if r.Err != nil {
        return nil, r.Err
    }

    return r.StringArray(), nil
}

func (r *Reply) elemOrErr() (Elem, error) {
    if r.Err != nil {
        return nil, r.Err
    }

    return r.Elem, nil
}

func (r *Reply) replyOrErr() (*Reply, error) {
    if r.Err != nil {
        return nil, r.Err
    }

    return r, nil
}

func smapToArr(mapping map[string]string) []string {
    buf := make([]string, len(mapping)*2)
    n := 0

    for k, v := range mapping {
        buf[n], buf[n+1] = k, v
        n += 2
    }

    return buf
}

// generic

// Delete a key
func (c *Client) Del(keys ...string) (int64, error) {
    return SendStr(c, "DEL", keys...).intOrErr()
}

// Determine if a key exists
func (c *Client) Exists(key string) (bool, error) {
    return Send(c, []byte("EXISTS"), []byte(key)).boolOrErr()
}

// Set a key's time to live in seconds
func (c *Client) Expire(key string, seconds int64) (bool, error) {
    return SendStr(c, "EXPIRE", key, strconv.FormatInt(seconds, 10)).boolOrErr()
}

// Set the expiration for a key as a UNIX timestamp
func (c *Client) Expireat(key string, timestamp int64) (bool, error) {
    return SendStr(c, "EXPIREAT", key, strconv.FormatInt(timestamp, 10)).boolOrErr()
}

// Find all keys matching the given pattern
func (c *Client) Keys(pattern string) ([]string, error) {
    return SendStr(c, "KEYS", pattern).stringArrOrErr()
}

// Move a key to another database
func (c *Client) Move(key string, db int) (bool, error) {
    return SendStr(c, "MOVE", key, strconv.Itoa(db)).boolOrErr()
}

// Remove the expiration from a key
func (c *Client) Persist(key string) (bool, error) {
    return SendStr(c, "PERSIST", key).boolOrErr()
}

// Return a random key from the keyspace
func (c *Client) Randomkey() (string, error) {
    return SendStr(c, "RANDOMKEY").stringOrErr()
}

// Rename a key
func (c *Client) Rename(key string, newkey string) error {
    return SendStr(c, "RENAME", key, newkey).nilOrErr()
}

// Rename a key, only if the new key does not exist
func (c *Client) Renamenx(key string, newkey string) (bool, error) {
    return SendStr(c, "RENAMENX", key, newkey).boolOrErr()
}

// Sort the elements in a list, set or sorted set
func (c *Client) Sort(key string, args ...string) (*Reply, error) {
    return SendStr(c, "SORT", append([]string{key}, args...)...).replyOrErr()
    ///out := make([]byte, len(v))

    ///for i, k := range v {
    ///    out[i] = string(k)
    ///}

    ///return out, nil
}

// Get the time to live for a key
func (c *Client) Ttl(key string) (int64, error) {
    return SendStr(c, "TTL", key).intOrErr()
}

// Determine the type stored at key
func (c *Client) Type(key string) (string, error) {
    return SendStr(c, "TYPE", key).stringOrErr()
}

// strings

// Append a value to a key
func (c *Client) Append(key string, value interface{}) (int64, error) {
    return SendIface(c, "APPEND", key, value).intOrErr()
}

// Decrement the integer value of a key by one
func (c *Client) Decr(key string) (int64, error) {
    return SendStr(c, "DECR", key).intOrErr()
}

// Decrement the integer value of a key by the given number
func (c *Client) Decrby(key string, decrement int64) (int64, error) {
    return SendStr(c, "DECRBY", key, strconv.FormatInt(decrement, 10)).intOrErr()
}

// Get the value of a key
func (c *Client) Get(key string) (Elem, error) {
    r := SendStr(c, "GET", key)

    if r.Err == nil && r.Elem == nil {
        r.Err = errors.New("key `" + key + "`does not exist")
    }

    return r.elemOrErr()
}

// Returns the bit value at offset in the string value stored at key
func (c *Client) Getbit(key string, offset int) (int64, error) {
    return SendStr(c, "GETBIT", key, strconv.Itoa(offset)).intOrErr()
}

// Get a substring of the string stored at a key
func (c *Client) Getrange(key string, start int, end int) (Elem, error) {
    return SendStr(c, "GETRANGE", key, strconv.Itoa(start), strconv.Itoa(end)).elemOrErr()
}

// Set the string value of a key and return its old value
func (c *Client) Getset(key string, value interface{}) (Elem, error) {
    return SendIface(c, "GETSET", key, value).elemOrErr()
}

// Increment the integer value of a key by one
func (c *Client) Incr(key string) (int64, error) {
    return SendStr(c, "INCR", key).intOrErr()
}

// Increment the integer value of a key by the given number
func (c *Client) Incrby(key string, increment int64) (int64, error) {
    return SendStr(c, "INCRBY", key, strconv.FormatInt(increment, 10)).intOrErr()
}

// Get the values of all the given keys
func (c *Client) Mget(keys ...string) (*Reply, error) {
    return SendStr(c, "MGET", keys...).replyOrErr()
}

// Set multiple keys to multiple values
func (c *Client) Mset(mapping map[string]string) error {
    return SendStr(c, "MSET", smapToArr(mapping)...).nilOrErr()
}

// Set multiple keys to multiple values, only if none of the keys exist
func (c *Client) Msetnx(mapping map[string]string) (bool, error) {
    return SendStr(c, "MSETNX", smapToArr(mapping)...).boolOrErr()
}

// Set the string value of a key
func (c *Client) Set(key string, value interface{}) error {
    return SendIface(c, "SET", key, value).nilOrErr()
}

// Sets or clears the bit at offset in the string value stored at key
func (c *Client) Setbit(key string, offset int, value int) (int64, error) {
    return SendStr(c, "SETBIT", key, strconv.Itoa(offset), strconv.Itoa(value)).intOrErr()
}

// Set the value and expiration of a key
func (c *Client) Setex(key string, seconds int64, value interface{}) error {
    return SendIface(c, "SETEX", key, seconds, value).nilOrErr()
}

// Set the value of a key, only if the key does not exist
func (c *Client) Setnx(key string, value interface{}) (bool, error) {
    return SendIface(c, "SETNX", key, value).boolOrErr()
}

// Overwrite part of a string at key starting at the specified offset
func (c *Client) Setrange(key string, offset int, value interface{}) (int64, error) {
    return SendIface(c, "SETRANGE", key, offset, value).intOrErr()
}

// Get the length of the value stored in a key
func (c *Client) Strlen(key string) (int64, error) {
    return SendStr(c, "STRLEN", key).intOrErr()
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
func (c *Client) Lindex(key string, index int) (Elem, error) {
    return SendStr(c, "LINDEX", key, strconv.Itoa(index)).elemOrErr()
}

// Insert an element before or after another element in a list
func (c *Client) Linsert(key, where string, pivot, value interface{}) (int64, error) {
    return SendIface(c, "LINSERT", key, where, pivot, value).intOrErr()
}

// Get the length of a list
func (c *Client) Llen(key string) (int64, error) {
    return SendStr(c, "LLEN", key).intOrErr()
}

// Remove and get the first element in a list
func (c *Client) Lpop(key string) (Elem, error) {
    return SendStr(c, "LPOP", key).elemOrErr()
}

// Prepend a value to a list
func (c *Client) Lpush(key string, value interface{}) (int64, error) {
    return SendIface(c, "LPUSH", key, value).intOrErr()
}

// Prepend a value to a list, only if the list exists
func (c *Client) Lpushx(key string, value interface{}) (int64, error) {
    return SendIface(c, "LPUSHX", key, value).intOrErr()
}

// Get a range of elements from a list
func (c *Client) Lrange(key string, start, stop int) (*Reply, error) {
    return SendStr(c, "LRANGE", key, strconv.Itoa(start), strconv.Itoa(stop)).replyOrErr()
}

// Remove elements from a list
func (c *Client) Lrem(key string, count int, value interface{}) (int64, error) {
    return SendIface(c, "LREM", key, count, value).intOrErr()
}

// Set the value of an element in a list by its index
func (c *Client) Lset(key string, index int, value interface{}) error {
    return SendIface(c, "LSET", key, strconv.Itoa(index), value).nilOrErr()
}

// Trim a list to the specified range
func (c *Client) Ltrim(key string, start int, stop int) error {
    return SendStr(c, "LTRIM", key, strconv.Itoa(start), strconv.Itoa(stop)).nilOrErr()
}

// Remove and get the last element in a list
func (c *Client) Rpop(key string) (Elem, error) {
    return SendStr(c, "RPOP", key).elemOrErr()
}

// Remove the last element in a list, append it to another list and return it
func (c *Client) Rpoplpush(source string, destination string) (Elem, error) {
    return SendStr(c, "RPOPLPUSH", source, destination).elemOrErr()
}

// Append a value to a list
func (c *Client) Rpush(key string, value interface{}) (int64, error) {
    return SendIface(c, "RPUSH", key, value).intOrErr()
}

// Append a value to a list, only if the list exists
func (c *Client) Rpushx(key string, value interface{}) (int64, error) {
    return SendIface(c, "RPUSHX", key, value).intOrErr()
}

// hash

// Delete a hash field
func (c *Client) Hdel(key string, field string) (bool, error) {
    return SendStr(c, "HDEL", key, field).boolOrErr()
}

// Determine if a hash field exists
func (c *Client) Hexists(key string, field string) (bool, error) {
    return SendStr(c, "HEXISTS", key, field).boolOrErr()
}

// Get the value of a hash field
func (c *Client) Hget(key string, field string) (Elem, error) {
    return SendStr(c, "HGET", key, field).elemOrErr()
}

// Get all the fields and values in a hash
func (c *Client) Hgetall(key string) (*Reply, error) {
    return SendStr(c, "HGETALL", key).replyOrErr()
}

// Increment the integer value of a hash field by the given number
func (c *Client) Hincrby(key string, field string, increment int64) (int64, error) {
    return SendStr(c, "HINCRBY", key, field, strconv.FormatInt(increment, 10)).intOrErr()
}

// Get all the fields in a hash
func (c *Client) Hkeys(key string) ([]string, error) {
    return SendStr(c, "HKEYS", key).stringArrOrErr()
}

// Get the number of fields in a hash
func (c *Client) Hlen(key string) (int64, error) {
    return SendStr(c, "HLEN", key).intOrErr()
}

// Get the values of all the given hash fields
func (c *Client) Hmget(key string, fields ...string) (*Reply, error) {
    return SendStr(c, "HMGET", append([]string{key}, fields...)...).replyOrErr()
}

// Set multiple hash fields to multiple values
func (c *Client) Hmset(key string, mapping map[string]interface{}) error {
    buf := make([]interface{}, len(mapping)*2+1)
    buf[0] = key
    n := 1

    for k, v := range mapping {
        buf[n], buf[n+1] = k, v
        n += 2
    }

    return SendIface(c, "HMSET", buf...).nilOrErr()
}

// Set the string value of a hash field
func (c *Client) Hset(key string, field string, value interface{}) (bool, error) {
    return SendIface(c, "HSET", key, field, value).boolOrErr()
}

// Set the value of a hash field, only if the field does not exist
func (c *Client) Hsetnx(key string, field string, value interface{}) (bool, error) {
    return SendIface(c, "HSETNX", key, field, value).boolOrErr()
}

// Get all the values in a hash
func (c *Client) Hvals(key string) (*Reply, error) {
    return SendStr(c, "HVALS", key).replyOrErr()
}

// set

// Add a member to a set
func (c *Client) Sadd(key string, member interface{}) (bool, error) {
    return SendIface(c, "SADD", key, member).boolOrErr()
}

// Get the number of members in a set
func (c *Client) Scard(key string) (int64, error) {
    return SendStr(c, "SCARD", key).intOrErr()
}

// Subtract multiple sets
func (c *Client) Sdiff(keys ...string) (*Reply, error) {
    return SendStr(c, "SDIFF", keys...).replyOrErr()
}

// Subtract multiple sets and store the resulting set in a key
func (c *Client) Sdiffstore(destination string, keys ...string) (int64, error) {
    return SendStr(c, "SDIFFSTORE", append([]string{destination}, keys...)...).intOrErr()
}

// Intersect multiple sets
func (c *Client) Sinter(keys ...string) (*Reply, error) {
    return SendStr(c, "SINTER", keys...).replyOrErr()
}

// Intersect multiple sets and store the resulting set in a key
func (c *Client) Sinterstore(destination string, keys ...string) (int64, error) {
    return SendStr(c, "SINTERSTORE", append([]string{destination}, keys...)...).intOrErr()
}

// Determine if a given value is a member of a set
func (c *Client) Sismember(key string, member interface{}) (bool, error) {
    return SendIface(c, "SISMEMBER", key, member).boolOrErr()
}

// Get all the members in a set
func (c *Client) Smembers(key string) (*Reply, error) {
    return SendStr(c, "SMEMBERS", key).replyOrErr()
}

// Move a member from one set to another
func (c *Client) Smove(source string, destination string, member interface{}) (bool, error) {
    return SendIface(c, "SMOVE", source, destination, member).boolOrErr()
}

// Remove and return a random member from a set
func (c *Client) Spop(key string) (Elem, error) {
    return SendStr(c, "SPOP", key).elemOrErr()
}

// Get a random member from a set
func (c *Client) Srandmember(key string) (Elem, error) {
    return SendStr(c, "SRANDMEMBER", key).elemOrErr()
}

// Remove a member from a set
func (c *Client) Srem(key string, member interface{}) (bool, error) {
    return SendIface(c, "SREM", key, member).boolOrErr()
}

// Add multiple sets
func (c *Client) Sunion(keys ...string) (*Reply, error) {
    return SendStr(c, "SUNION", keys...).replyOrErr()
}

// Add multiple sets and store the resulting set in a key
func (c *Client) Sunionstore(destination string, keys ...string) (int64, error) {
    return SendStr(c, "SUNIONSTORE", append([]string{destination}, keys...)...).intOrErr()
}

// sorted_set

// Add a member to a sorted set, or update its score if it already exists
func (c *Client) Zadd(key string, score float64, member interface{}) (bool, error) {
    return SendIface(c, "ZADD", key, score, member).boolOrErr()
}

// Get the number of members in a sorted set
func (c *Client) Zcard(key string) (int64, error) {
    return SendStr(c, "ZCARD", key).intOrErr()
}

// Count the members in a sorted set with scores within the given values
func (c *Client) Zcount(key string, min float64, max float64) (int64, error) {
    return SendStr(c, "ZCOUNT", key, strconv.FormatFloat(min, 'f', -1, 64), strconv.FormatFloat(max, 'f', -1, 64)).intOrErr()
}

// Increment the score of a member in a sorted set
func (c *Client) Zincrby(key string, increment float64, member interface{}) (float64, error) {
    return SendIface(c, "ZINCRBY", key, strconv.FormatFloat(increment, 'f', -1, 64), member).floatOrErr()
}

// Intersect multiple sorted sets and store the resulting sorted set in a new key
// `numkeys` is determined by the len of `keys` param
func (c *Client) Zinterstore(destination string, keys []string, args ...string) (int64, error) {
    a := append([]string{destination, strconv.Itoa(len(keys))}, keys...)
    a = append(a, args...)
    return SendStr(c, "ZINTERSTORE", a...).intOrErr()
}

// Return a range of members in a sorted set, by index
// TODO: add WITHSCORES keyword
func (c *Client) Zrange(key string, start int, stop int) (*Reply, error) {
    return SendStr(c, "ZRANGE", key, strconv.Itoa(start), strconv.Itoa(stop)).replyOrErr()
}

// Return a range of members in a sorted set, by score
func (c *Client) Zrangebyscore(key string, min string, max string, args ...string) (*Reply, error) {
    a := append([]string{key, min, max}, args...)
    return SendStr(c, "ZRANGEBYSCORE", a...).replyOrErr()
}

// Determine the index of a member in a sorted set
// TODO: should cast an error when member does not exist
func (c *Client) Zrank(key string, member interface{}) (int64, error) {
    return SendIface(c, "ZRANK", key, member).intOrErr()
}

// Remove a member from a sorted set
func (c *Client) Zrem(key string, member interface{}) (bool, error) {
    return SendIface(c, "ZREM", key, member).boolOrErr()
}

// Remove all members in a sorted set within the given indexes
func (c *Client) Zremrangebyrank(key string, start int, stop int) (int64, error) {
    return SendStr(c, "ZREMRANGEBYRANK", key, strconv.Itoa(start), strconv.Itoa(stop)).intOrErr()
}

// Remove all members in a sorted set within the given scores
func (c *Client) Zremrangebyscore(key string, min float64, max float64) (int64, error) {
    return SendStr(c, "ZREMRANGEBYSCORE", key, strconv.FormatFloat(min, 'f', -1, 64), strconv.FormatFloat(max, 'f', -1, 64)).intOrErr()
}

// Return a range of members in a sorted set, by index, with scores ordered from high to low
func (c *Client) Zrevrange(key string, start int, stop int, args ...string) (*Reply, error) {
    a := append([]string{key, strconv.Itoa(start), strconv.Itoa(stop)}, args...)
    return SendStr(c, "ZREVRANGE", a...).replyOrErr()
}

// Return a range of members in a sorted set, by score, with scores ordered from high to low
func (c *Client) Zrevrangebyscore(key string, max float64, min float64, args ...string) (*Reply, error) {
    a := append([]string{key, strconv.FormatFloat(max, 'f', -1, 64), strconv.FormatFloat(min, 'f', -1, 64)}, args...)
    return SendStr(c, "ZREVRANGEBYSCORE", a...).replyOrErr()
}

// Determine the index of a member in a sorted set, with scores ordered from high to low
func (c *Client) Zrevrank(key string, member interface{}) (int64, error) {
    return SendIface(c, "ZREVRANK", key, member).intOrErr()
}

// Get the score associated with the given member in a sorted set
func (c *Client) Zscore(key string, member interface{}) (float64, error) {
    return SendIface(c, "ZSCORE", key, member).floatOrErr()
}

// Add multiple sorted sets and store the resulting sorted set in a new key
// `numkeys` is determined by the len of `keys` param
func (c *Client) Zunionstore(destination string, keys []string, args ...string) (int64, error) {
    a := append([]string{destination, strconv.Itoa(len(keys))}, args...)
    return SendStr(c, "ZUNIONSTORE", a...).intOrErr()
}

//// transactions
//
//// Discard all commands issued after MULTI
//func (c *Client) Discard() (bool, os.Error) {
//    return Send(c, "DISCARD").boolOrErr()
//}
//
//// Execute all commands issued after MULTI
//func (c *Client) Exec() (*Reply, os.Error) {
//    return Send(c, "EXEC").replyOrErr()
//}
//
//// Mark the start of a transaction block
//func (c *Client) Multi() (bool, os.Error) {
//    return Send(c, "MULTI").boolOrErr()
//}
//
//// Forget about all watched keys
//func (c *Client) Unwatch() (bool, os.Error) {
//    return Send(c, "UNWATCH").boolOrErr()
//}
//
// server

// Asynchronously rewrite the append-only file
func (c *Client) Bgrewriteaof() error {
    return Send(c, []byte("BGREWRITEAOF")).nilOrErr()
}

// Asynchronously save the dataset to disk
func (c *Client) Bgsave() error {
    return Send(c, []byte("BGSAVE")).nilOrErr()
}

// Get the value of a configuration parameter
func (c *Client) ConfigGet(parameter string) (Elem, error) {
    return SendStr(c, "CONFIG GET", parameter).elemOrErr()
}

// Reset the stats returned by INFO
func (c *Client) ConfigResetstat() error {
    return Send(c, []byte("CONFIG RESETSTAT")).nilOrErr()
}

// Set a configuration parameter to the given value
func (c *Client) ConfigSet(parameter string, value string) error {
    return SendStr(c, "CONFIG SET", parameter, value).nilOrErr()
}

// Return the number of keys in the selected database
func (c *Client) Dbsize() (int64, error) {
    return Send(c, []byte("DBSIZE")).intOrErr()
}

// Get debugging information about a key
func (c *Client) DebugObject(key string) (Elem, error) {
    return SendStr(c, "DEBUG OBJECT", key).elemOrErr()
}

// Make the server crash
func (c *Client) DebugSegfault() error {
    return Send(c, []byte("DEBUG SEGFAULT")).nilOrErr()
}

// Remove all keys from all databases
func (c *Client) Flushall() error {
    return Send(c, []byte("FLUSHALL")).nilOrErr()
}

// Remove all keys from the current database
func (c *Client) Flushdb() error {
    return Send(c, []byte("FLUSHDB")).nilOrErr()
}

// Get information and statistics about the server
func (c *Client) Info() (Elem, error) {
    return Send(c, []byte("INFO")).elemOrErr()
}

// Get the UNIX time stamp of the last successful save to disk
func (c *Client) Lastsave() (int64, error) {
    return Send(c, []byte("LASTSAVE")).intOrErr()
}

// Listen for all requests received by the server in real time
func (c *Client) Monitor() (*Reply, error) {
    return Send(c, []byte("MONITOR")).replyOrErr()
}

// Synchronously save the dataset to disk
func (c *Client) Save() error {
    return Send(c, []byte("SAVE")).nilOrErr()
}

// Synchronously save the dataset to disk and then shut down the server
func (c *Client) Shutdown() error {
    return Send(c, []byte("SHUTDOWN")).nilOrErr()
}

// connection

//// Authenticate to the server
//func (c *Client) Auth(password string) os.Error {
//    return Send(c, "AUTH").
//}

// Echo the given string
func (c *Client) Echo(message interface{}) (Elem, error) {
    return SendIface(c, "ECHO", message).elemOrErr()
}

// Ping the server
func (c *Client) Ping() (Elem, error) {
    return Send(c, []byte("PING")).elemOrErr()
}

// Close the connection
func (c *Client) Quit() error {
    return Send(c, []byte("QUIT")).nilOrErr()
}

// Change the selected database for the current connection
func (c *Client) Select(index int) error {
    for i := 0; i < MaxClientConn; i++ {
        c.pool.pop()
    }

    c.Db = index
    c.pool = newPool()
    err := SendStr(c, "SELECT", strconv.Itoa(index)).nilOrErr()
    return err
}

// pubsub

// Post a message to a channel
func (c *Client) Publish(channel string, message interface{}) (int64, error) {
    return SendIface(c, "PUBLISH", channel, message).intOrErr()
}

// Listen for messages published to the given channels
func (c *Client) Subscribe(channels ...string) (*Sub, error) {
    s := &Sub{c: c}
    err := s.Subscribe(channels...)
    return s, err
}

// Listen for messages published to channels matching the given patterns
func (c *Client) Psubscribe(patterns ...string) (*Sub, error) {
    s := &Sub{c: c}
    err := s.Psubscribe(patterns...)
    return s, err
}

// Stop listening for messages posted to channels matching the given patterns
func (s *Sub) Punsubscribe(patterns ...string) error {
    if !s.subscribed {
        return errors.New("Cannot PUNSUBSCRIBE before subscribing")
    }

    return appendSendStr(s, "PUNSUBSCRIBE", patterns...).Err
}

// Listen for messages published to channels matching the given patterns
func (s *Sub) Psubscribe(patterns ...string) error {
    if err := appendSendStr(s, "PSUBSCRIBE", patterns...).Err; err != nil {
        return err
    }

    if !s.subscribed {
        s.subscribe()
    }

    return nil
}

// Stop listening for messages posted to the given channels
func (s *Sub) Unsubscribe(channels ...string) error {
    if !s.subscribed {
        return errors.New("Cannot UNSUBSCRIBE before subscribing")
    }

    return appendSendStr(s, "UNSUBSCRIBE", channels...).Err
}

// Listen for messages published to the given channels
func (s *Sub) Subscribe(channels ...string) error {
    if err := appendSendStr(s, "SUBSCRIBE", channels...).Err; err != nil {
        return err
    }

    if !s.subscribed {
        s.subscribe()
    }

    return nil
}
