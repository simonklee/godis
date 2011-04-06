package godis

import (
    "os"
    "strconv"
)

// helpers for handling common return values
func (r *Reply) intOrErr() (int64, os.Error) {
    if r.Err != nil {
        return -1, r.Err
    }

    return r.Elem.Int64(), nil
}

// helpers for handling common return values
func (r *Reply) floatOrErr() (float64, os.Error) {
    if r.Err != nil {
        return -1.0, r.Err
    }

    return r.Elem.Float64(), nil
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
func Del(rw ReaderWriter, keys ...string) (int64, os.Error) {
    return SendStr(rw, "DEL", keys...).intOrErr()
}

// Determine if a key exists
func Exists(rw ReaderWriter, key string) (bool, os.Error) {
    return Send(rw, []byte("EXISTS"), []byte(key)).boolOrErr()
}

// Set a key's time to live in seconds
func Expire(rw ReaderWriter, key string, seconds int64) (bool, os.Error) {
    return SendStr(rw, "EXPIRE", key, strconv.Itoa64(seconds)).boolOrErr()
}

// Set the expiration for a key as a UNIX timestamp
func Expireat(rw ReaderWriter, key string, timestamp int64) (bool, os.Error) {
    return SendStr(rw, "EXPIREAT", key, strconv.Itoa64(timestamp)).boolOrErr()
}

// Find all keys matching the given pattern
func Keys(rw ReaderWriter, pattern string) ([]string, os.Error) {
    return SendStr(rw, "KEYS", pattern).stringArrOrErr()
}

// Move a key to another database
func Move(rw ReaderWriter, key string, db int) (bool, os.Error) {
    return SendStr(rw, "MOVE", key, strconv.Itoa(db)).boolOrErr()
}

// Remove the expiration from a key
func Persist(rw ReaderWriter, key string) (bool, os.Error) {
    return SendStr(rw, "PERSIST", key).boolOrErr()
}

// Return a random key from the keyspace
func Randomkey(rw ReaderWriter, ) (string, os.Error) {
    return SendStr(rw, "RANDOMKEY").stringOrErr()
}

// Rename a key
func Rename(rw ReaderWriter, key string, newkey string) os.Error {
    return SendStr(rw, "RENAME", key, newkey).nilOrErr()
}

// Rename a key, only if the new key does not exist
func Renamenx(rw ReaderWriter, key string, newkey string) (bool, os.Error) {
    return SendStr(rw, "RENAMENX", key, newkey).boolOrErr()
}

// Sort the elements in a list, set or sorted set
func Sort(rw ReaderWriter, key string, args ...string) (*Reply, os.Error) {
    return SendStr(rw, "SORT", append([]string{key}, args...)...).replyOrErr()
    ///out := make([]byte, len(v))

    ///for i, k := range v {
    ///    out[i] = string(k)
    ///}

    ///return out, nil
}

// Get the time to live for a key
func Ttl(rw ReaderWriter, key string) (int64, os.Error) {
    return SendStr(rw, "TTL", key).intOrErr()
}

// Determine the type stored at key
func Type(rw ReaderWriter, key string) (string, os.Error) {
    return SendStr(rw, "TYPE", key).stringOrErr()
}

// strings

// Append a value to a key
func Append(rw ReaderWriter, key string, value interface{}) (int64, os.Error) {
    return SendIface(rw, "APPEND", key, value).intOrErr()
}

// Decrement the integer value of a key by one
func Decr(rw ReaderWriter, key string) (int64, os.Error) {
    return SendStr(rw, "DECR", key).intOrErr()
}

// Decrement the integer value of a key by the given number
func Decrby(rw ReaderWriter, key string, decrement int64) (int64, os.Error) {
    return SendStr(rw, "DECRBY", key, strconv.Itoa64(decrement)).intOrErr()
}

// Get the value of a key
func Get(rw ReaderWriter, key string) (string, os.Error) {
    r := SendStr(rw, "GET", key)

    if r.Err == nil && r.Elem == nil {
        r.Err = os.NewError("key `" + key + "`does not exist")
    }

    return r.stringOrErr()
}

// Returns the bit value at offset in the string value stored at key
func Getbit(rw ReaderWriter, key string, offset int) (int64, os.Error) {
    return SendStr(rw, "GETBIT", key, strconv.Itoa(offset)).intOrErr()
}

// Get a substring of the string stored at a key
func Getrange(rw ReaderWriter, key string, start int, end int) (string, os.Error) {
    return SendStr(rw, "GETRANGE", key, strconv.Itoa(start), strconv.Itoa(end)).stringOrErr()
}

// Set the string value of a key and return its old value
func Getset(rw ReaderWriter, key string, value interface{}) (string, os.Error) {
    return SendIface(rw, "GETSET", key, value).stringOrErr()
}

// Increment the integer value of a key by one
func Incr(rw ReaderWriter, key string) (int64, os.Error) {
    return SendStr(rw, "INCR", key).intOrErr()
}

// Increment the integer value of a key by the given number
func Incrby(rw ReaderWriter, key string, increment int64) (int64, os.Error) {
    return SendStr(rw, "INCRBY", key, strconv.Itoa64(increment)).intOrErr()
}

// Get the values of all the given keys
func Mget(rw ReaderWriter, keys ...string) ([]string, os.Error) {
    return SendStr(rw, "MGET", keys...).stringArrOrErr()
}

// Set multiple keys to multiple values
func Mset(rw ReaderWriter, mapping map[string]string) os.Error {
    return SendStr(rw, "MSET", smapToArr(mapping)...).nilOrErr()
}

// Set multiple keys to multiple values, only if none of the keys exist
func Msetnx(rw ReaderWriter, mapping map[string]string) (bool, os.Error) {
    return SendStr(rw, "MSETNX", smapToArr(mapping)...).boolOrErr()
}

// Set the string value of a key
func Set(rw ReaderWriter, key string, value interface{}) os.Error {
    return SendIface(rw, "SET", key, value).nilOrErr()
}

// Sets or clears the bit at offset in the string value stored at key
func Setbit(rw ReaderWriter, key string, offset int, value int) (int64, os.Error) {
    return SendStr(rw, "SETBIT", key, strconv.Itoa(offset), strconv.Itoa(value)).intOrErr()
}

// Set the value and expiration of a key
func Setex(rw ReaderWriter, key string, seconds int64, value interface{}) os.Error {
    return SendIface(rw, "SETEX", key, seconds, value).nilOrErr()
}

// Set the value of a key, only if the key does not exist
func Setnx(rw ReaderWriter, key string, value interface{}) (bool, os.Error) {
    return SendIface(rw, "SETNX", key, value).boolOrErr()
}

// Overwrite part of a string at key starting at the specified offset
func Setrange(rw ReaderWriter, key string, offset int, value interface{}) (int64, os.Error) {
    return SendIface(rw, "SETRANGE", key, offset, value).intOrErr()
}

// Get the length of the value stored in a key
func Strlen(rw ReaderWriter, key string) (int64, os.Error) {
    return SendStr(rw, "STRLEN", key).intOrErr()
}

// list

// Remove and get the first element in a list, or block until one is available
//func Blpop(rw ReaderWriter, key []string, timeout int64) [][]byte {
//
//}

// Remove and get the last element in a list, or block until one is available
//func Brpop(rw ReaderWriter, key []string, timeout int64) [][]byte {
//
//}

// Pop a value from a list, push it to another list and return it; or block until one is available
//func Brpoplpush(rw ReaderWriter, source string, destination string, timeout int64) []byte {
//
//}

// Get an element from a list by its index
func Lindex(rw ReaderWriter, key string, index int) (Elem, os.Error) {
    return SendStr(rw, "LINDEX", key, strconv.Itoa(index)).elemOrErr()
}

// Insert an element before or after another element in a list
func Linsert(rw ReaderWriter, key, where string, pivot, value interface{}) (int64, os.Error) {
    return SendIface(rw, "LINSERT", key, where, pivot, value).intOrErr()
}

// Get the length of a list
func Llen(rw ReaderWriter, key string) (int64, os.Error) {
    return SendStr(rw, "LLEN", key).intOrErr()
}

// Remove and get the first element in a list
func Lpop(rw ReaderWriter, key string) (Elem, os.Error) {
    return SendStr(rw, "LPOP", key).elemOrErr()
}

// Prepend a value to a list
func Lpush(rw ReaderWriter, key string, value interface{}) (int64, os.Error) {
    return SendIface(rw, "LPUSH", key, value).intOrErr()
}

// Prepend a value to a list, only if the list exists
func Lpushx(rw ReaderWriter, key string, value interface{}) (int64, os.Error) {
    return SendIface(rw, "LPUSHX", key, value).intOrErr()
}

// Get a range of elements from a list
func Lrange(rw ReaderWriter, key string, start, stop int) (*Reply, os.Error) {
    return SendStr(rw, "LRANGE", key, strconv.Itoa(start), strconv.Itoa(stop)).replyOrErr()
}

// Remove elements from a list
func Lrem(rw ReaderWriter, key string, count int, value interface{}) (int64, os.Error) {
    return SendIface(rw, "LREM", key, count, value).intOrErr()
}

// Set the value of an element in a list by its index
func Lset(rw ReaderWriter, key string, index int, value interface{}) os.Error {
    return SendIface(rw, "LSET", key, strconv.Itoa(index), value).nilOrErr()
}

// Trim a list to the specified range
func Ltrim(rw ReaderWriter, key string, start int, stop int) os.Error {
    return SendStr(rw, "LTRIM", key, strconv.Itoa(start), strconv.Itoa(stop)).nilOrErr()
}

// Remove and get the last element in a list
func Rpop(rw ReaderWriter, key string) (Elem, os.Error) {
    return SendStr(rw, "RPOP", key).elemOrErr()
}

// Remove the last element in a list, append it to another list and return it
func Rpoplpush(rw ReaderWriter, source string, destination string) (Elem, os.Error) {
    return SendStr(rw, "RPOPLPUSH", source, destination).elemOrErr()
}

// Append a value to a list
func Rpush(rw ReaderWriter, key string, value interface{}) (int64, os.Error) {
    return SendIface(rw, "RPUSH", key, value).intOrErr()
}

// Append a value to a list, only if the list exists
func Rpushx(rw ReaderWriter, key string, value interface{}) (int64, os.Error) {
    return SendIface(rw, "RPUSHX", key, value).intOrErr()
}

// hash

// Delete a hash field
func Hdel(rw ReaderWriter, key string, field string) (bool, os.Error) {
    return SendStr(rw, "HDEL", key, field).boolOrErr()
}

// Determine if a hash field exists
func Hexists(rw ReaderWriter, key string, field string) (bool, os.Error) {
    return SendStr(rw, "HEXISTS", key, field).boolOrErr()
}

// Get the value of a hash field
func Hget(rw ReaderWriter, key string, field string) (Elem, os.Error) {
    return SendStr(rw, "HGET", key, field).elemOrErr()
}

// Get all the fields and values in a hash
func Hgetall(rw ReaderWriter, key string) (*Reply, os.Error) {
    return SendStr(rw, "HGETALL", key).replyOrErr()
}

// Increment the integer value of a hash field by the given number
func Hincrby(rw ReaderWriter, key string, field string, increment int64) (int64, os.Error) {
    return SendStr(rw, "HINCRBY", key, field, strconv.Itoa64(increment)).intOrErr()
}

// Get all the fields in a hash
func Hkeys(rw ReaderWriter, key string) ([]string, os.Error) {
    return SendStr(rw, "HKEYS", key).stringArrOrErr()
}

// Get the number of fields in a hash
func Hlen(rw ReaderWriter, key string) (int64, os.Error) {
    return SendStr(rw, "HLEN", key).intOrErr()
}

// Get the values of all the given hash fields
func Hmget(rw ReaderWriter, key string, fields ...string) (*Reply, os.Error) {
    return SendStr(rw, "HMGET", append([]string{key}, fields...)...).replyOrErr()
}

// Set multiple hash fields to multiple values
func Hmset(rw ReaderWriter, key string, mapping map[string]interface{}) os.Error {
    buf := make([]interface{}, len(mapping)*2+1)
    buf[0] = key
    n := 1

    for k, v := range mapping {
        buf[n], buf[n+1] = k, v
        n += 2
    }

    return SendIface(rw, "HMSET", buf...).nilOrErr()
}

// Set the string value of a hash field
func Hset(rw ReaderWriter, key string, field string, value interface{}) (bool, os.Error) {
    return SendIface(rw, "HSET", key, field, value).boolOrErr()
}

// Set the value of a hash field, only if the field does not exist
func Hsetnx(rw ReaderWriter, key string, field string, value interface{}) (bool, os.Error) {
    return SendIface(rw, "HSETNX", key, field, value).boolOrErr()
}

// Get all the values in a hash
func Hvals(rw ReaderWriter, key string) (*Reply, os.Error) {
    return SendStr(rw, "HVALS", key).replyOrErr()
}

// set

// Add a member to a set
func Sadd(rw ReaderWriter, key string, member interface{}) (bool, os.Error) {
    return SendIface(rw, "SADD", key, member).boolOrErr()
}

// Get the number of members in a set
func Scard(rw ReaderWriter, key string) (int64, os.Error) {
    return SendStr(rw, "SCARD", key).intOrErr()
}

// Subtract multiple sets
func Sdiff(rw ReaderWriter, keys ...string) (*Reply, os.Error) {
    return SendStr(rw, "SDIFF", keys...).replyOrErr()
}

// Subtract multiple sets and store the resulting set in a key
func Sdiffstore(rw ReaderWriter, destination string, keys ...string) (int64, os.Error) {
    return SendStr(rw, "SDIFFSTORE", append([]string{destination}, keys...)...).intOrErr()
}

// Intersect multiple sets
func Sinter(rw ReaderWriter, keys ...string) (*Reply, os.Error) {
    return SendStr(rw, "SINTER", keys...).replyOrErr()
}

// Intersect multiple sets and store the resulting set in a key
func Sinterstore(rw ReaderWriter, destination string, keys ...string) (int64, os.Error) {
    return SendStr(rw, "SINTERSTORE", append([]string{destination}, keys...)...).intOrErr()
}

// Determine if a given value is a member of a set
func Sismember(rw ReaderWriter, key string, member interface{}) (bool, os.Error) {
    return SendIface(rw, "SISMEMBER", key, member).boolOrErr()
}

// Get all the members in a set
func Smembers(rw ReaderWriter, key string) (*Reply, os.Error) {
    return SendStr(rw, "SMEMBERS", key).replyOrErr()
}

// Move a member from one set to another
func Smove(rw ReaderWriter, source string, destination string, member interface{}) (bool, os.Error) {
    return SendIface(rw, "SMOVE", source, destination, member).boolOrErr()
}

// Remove and return a random member from a set
func Spop(rw ReaderWriter, key string) (Elem, os.Error) {
    return SendStr(rw, "SPOP", key).elemOrErr()
}

// Get a random member from a set
func Srandmember(rw ReaderWriter, key string) (Elem, os.Error) {
    return SendStr(rw, "SRANDMEMBER", key).elemOrErr()
}

// Remove a member from a set
func Srem(rw ReaderWriter, key string, member interface{}) (bool, os.Error) {
    return SendIface(rw, "SREM", key, member).boolOrErr()
}

// Add multiple sets
func Sunion(rw ReaderWriter, keys ...string) (*Reply, os.Error) {
    return SendStr(rw, "SUNION", keys...).replyOrErr()
}

// Add multiple sets and store the resulting set in a key
func Sunionstore(rw ReaderWriter, destination string, keys ...string) (int64, os.Error) {
    return SendStr(rw, "SUNIONSTORE", append([]string{destination}, keys...)...).intOrErr()
}

// sorted_set

// Add a member to a sorted set, or update its score if it already exists
func Zadd(rw ReaderWriter, key string, score float64, member interface{}) (bool, os.Error) {
    return SendIface(rw, "ZADD", key, score, member).boolOrErr()
}

// Get the number of members in a sorted set
func Zcard(rw ReaderWriter, key string) (int64, os.Error) {
    return SendStr(rw, "ZCARD", key).intOrErr()
}

// Count the members in a sorted set with scores within the given values
func Zcount(rw ReaderWriter, key string, min float64, max float64) (int64, os.Error) {
    return SendStr(rw, "ZCOUNT", key, strconv.Ftoa64(min, 'f', -1), strconv.Ftoa64(max, 'f', -1)).intOrErr()
}

// Increment the score of a member in a sorted set
func Zincrby(rw ReaderWriter, key string, increment float64, member interface{}) (float64, os.Error) {
    return SendIface(rw, "ZINCRBY", key, strconv.Ftoa64(increment, 'f', -1), member).floatOrErr()
}

// Intersect multiple sorted sets and store the resulting sorted set in a new key
// `numkeys` is determined by the len of `keys` param
func Zinterstore(rw ReaderWriter, destination string, keys []string, args ...string) (int64, os.Error) {
    a := append([]string{destination, strconv.Itoa(len(keys))}, args...)
    return SendStr(rw, "ZINTERSTORE", a...).intOrErr()
}

// Return a range of members in a sorted set, by index
func Zrange(rw ReaderWriter, key string, start int, stop int) (*Reply, os.Error) {
    return SendStr(rw, "ZRANGE", key, strconv.Itoa(start), strconv.Itoa(stop)).replyOrErr()
}

// Return a range of members in a sorted set, by score
func Zrangebyscore(rw ReaderWriter, key string, min float64, max float64, args ...string) (*Reply, os.Error) {
    a := append([]string{key, strconv.Ftoa64(min, 'f', -1), strconv.Ftoa64(max, 'f', -1)}, args...)
    return SendStr(rw, "ZRANGEBYSCORE", a...).replyOrErr()
}

// Determine the index of a member in a sorted set
func Zrank(rw ReaderWriter, key string, member interface{}) (int64, os.Error) {
    return SendIface(rw, "ZRANK", key, member).intOrErr()
}

// Remove a member from a sorted set
func Zrem(rw ReaderWriter, key string, member interface{}) (int64, os.Error) {
    return SendIface(rw, "ZREM", key, member).intOrErr()
}

// Remove all members in a sorted set within the given indexes
func Zremrangebyrank(rw ReaderWriter, key string, start int, stop int) (int64, os.Error) {
    return SendStr(rw, "ZREMRANGEBYRANK", key, strconv.Itoa(start), strconv.Itoa(stop)).intOrErr()
}

// Remove all members in a sorted set within the given scores
func Zremrangebyscore(rw ReaderWriter, key string, min float64, max float64) (int64, os.Error) {
    return SendStr(rw, "ZREMRANGEBYSCORE", key, strconv.Ftoa64(min, 'f', -1), strconv.Ftoa64(max, 'f', -1)).intOrErr()
}

// Return a range of members in a sorted set, by index, with scores ordered from high to low
func Zrevrange(rw ReaderWriter, key string, start int, stop int, args ...string) (*Reply, os.Error) {
    a := append([]string{key, strconv.Itoa(start), strconv.Itoa(stop)}, args...)
    return SendStr(rw, "ZREVRANGE", a...).replyOrErr()
}

// Return a range of members in a sorted set, by score, with scores ordered from high to low
func Zrevrangebyscore(rw ReaderWriter, key string, max float64, min float64, args ...string) (*Reply, os.Error) {
    a := append([]string{key, strconv.Ftoa64(max, 'f', -1), strconv.Ftoa64(min, 'f', -1)}, args...)
    return SendStr(rw, "ZREVRANGEBYSCORE", a...).replyOrErr()
}

// Determine the index of a member in a sorted set, with scores ordered from high to low
func Zrevrank(rw ReaderWriter, key string, member interface{}) (int64, os.Error) {
    return SendIface(rw, "ZREVRANK", key, member).intOrErr()
}

// Get the score associated with the given member in a sorted set
func Zscore(rw ReaderWriter, key string, member interface{}) (float64, os.Error) {
    return SendIface(rw, "ZSCORE", key, member).floatOrErr()
}

// Add multiple sorted sets and store the resulting sorted set in a new key
// `numkeys` is determined by the len of `keys` param
func Zunionstore(rw ReaderWriter, destination string, keys []string, args ...string) (int64, os.Error) {
    a := append([]string{destination, strconv.Itoa(len(keys))}, args...)
    return SendStr(rw, "ZUNIONSTORE", a...).intOrErr()
}

//// pubsub
//
//// Listen for messages published to channels matching the given patterns
//func Psubscribe(rw ReaderWriter, mapping map[string][]byte) os.Error {
//    return Send(rw, "PSUBSCRIBE").
//}
//
//// Post a message to a channel
//func Publish(rw ReaderWriter, channel string, message string) (int64, os.Error) {
//    return Send(rw, "PUBLISH").intOrErr()
//}
//
//// Stop listening for messages posted to channels matching the given patterns
//func Punsubscribe(rw ReaderWriter, ) os.Error {
//    return Send(rw, "PUNSUBSCRIBE").
//}
//
//// Listen for messages published to the given channels
//func Subscribe(rw ReaderWriter, mapping map[string][]byte) os.Error {
//    return Send(rw, "SUBSCRIBE").
//}
//
//// Stop listening for messages posted to the given channels
//func Unsubscribe(rw ReaderWriter, ) os.Error {
//    return Send(rw, "UNSUBSCRIBE").
//}
//
//// transactions
//
//// Discard all commands issued after MULTI
//func Discard(rw ReaderWriter, ) (bool, os.Error) {
//    return Send(rw, "DISCARD").boolOrErr()
//}
//
//// Execute all commands issued after MULTI
//func Exec(rw ReaderWriter, ) (*Reply, os.Error) {
//    return Send(rw, "EXEC").replyOrErr()
//}
//
//// Mark the start of a transaction block
//func Multi(rw ReaderWriter, ) (bool, os.Error) {
//    return Send(rw, "MULTI").boolOrErr()
//}
//
//// Forget about all watched keys
//func Unwatch(rw ReaderWriter, ) (bool, os.Error) {
//    return Send(rw, "UNWATCH").boolOrErr()
//}
//
// server

// Asynchronously rewrite the append-only file
func Bgrewriteaof(rw ReaderWriter, ) os.Error {
    return Send(rw, []byte("BGREWRITEAOF")).nilOrErr()
}

// Asynchronously save the dataset to disk
func Bgsave(rw ReaderWriter, ) os.Error {
    return Send(rw, []byte("BGSAVE")).nilOrErr()
}

// Get the value of a configuration parameter
func ConfigGet(rw ReaderWriter, parameter string) (Elem, os.Error) {
    return SendStr(rw, "CONFIG GET", parameter).elemOrErr()
}

// Reset the stats returned by INFO
func ConfigResetstat(rw ReaderWriter, ) os.Error {
    return Send(rw, []byte("CONFIG RESETSTAT")).nilOrErr()
}

// Set a configuration parameter to the given value
func ConfigSet(rw ReaderWriter, parameter string, value string) os.Error {
    return SendStr(rw, "CONFIG SET", parameter, value).nilOrErr()
}

// Return the number of keys in the selected database
func Dbsize(rw ReaderWriter, ) (int64, os.Error) {
    return Send(rw, []byte("DBSIZE")).intOrErr()
}

// Get debugging information about a key
func DebugObject(rw ReaderWriter, key string) (Elem, os.Error) {
    return SendStr(rw, "DEBUG OBJECT", key).elemOrErr()
}

// Make the server crash
func DebugSegfault(rw ReaderWriter, ) os.Error {
    return Send(rw, []byte("DEBUG SEGFAULT")).nilOrErr()
}

// Remove all keys from all databases
func Flushall(rw ReaderWriter, ) os.Error {
    return Send(rw, []byte("FLUSHALL")).nilOrErr()
}

// Remove all keys from the current database
func Flushdb(rw ReaderWriter, ) os.Error {
    return Send(rw, []byte("FLUSHDB")).nilOrErr()
}

// Get information and statistics about the server
func Info(rw ReaderWriter, ) (Elem, os.Error) {
    return Send(rw, []byte("INFO")).elemOrErr()
}

// Get the UNIX time stamp of the last successful save to disk
func Lastsave(rw ReaderWriter, ) (int64, os.Error) {
    return Send(rw, []byte("LASTSAVE")).intOrErr()
}

// Listen for all requests received by the server in real time
func Monitor(rw ReaderWriter, ) (*Reply, os.Error) {
    return Send(rw, []byte("MONITOR")).replyOrErr()
}

// Synchronously save the dataset to disk
func Save(rw ReaderWriter, ) os.Error {
    return Send(rw, []byte("SAVE")).nilOrErr()
}

// Synchronously save the dataset to disk and then shut down the server
func Shutdown(rw ReaderWriter, ) os.Error {
    return Send(rw, []byte("SHUTDOWN")).nilOrErr()
}

// connection

//// Authenticate to the server
//func Auth(rw ReaderWriter, password string) os.Error {
//    return Send(rw, "AUTH").
//}

// Echo the given string
func Echo(rw ReaderWriter, message interface{}) (Elem, os.Error) {
    return SendIface(rw, "ECHO", message).elemOrErr()
}

// Ping the server
func Ping(rw ReaderWriter, ) (Elem, os.Error) {
    return Send(rw, []byte("PING")).elemOrErr()
}

// Close the connection
func Quit(rw ReaderWriter, ) os.Error {
    return Send(rw, []byte("QUIT")).nilOrErr()
}

// Change the selected database for the current connection
func Select(rw ReaderWriter, index int) os.Error {
    return SendStr(rw, "SELECT", strconv.Itoa(index)).nilOrErr()
}
