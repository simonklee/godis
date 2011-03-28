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

    return res.(string), nil
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
func (c *Client) Keys(pattern string) [][]byte {
    return [][]byte{}
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
func (c *Client) Randomkey() []byte {
    return []byte{}
}

// Rename a key, only if the new key does not exist
func (c *Client) Renamenx(key string, newkey string) int64 {
    return 0
}

// Sort the elements in a list, set or sorted set
func (c *Client) Sort(key string) [][]byte {
    return [][]byte{}
}

// Get the time to live for a key
func (c *Client) Ttl(key string) (int64, os.Error) {
    return intOrErr(c.Send("TTL", key))
}

// Determine the type stored at key
func (c *Client) Type(key string) (string, os.Error) {
    return stringOrErr(c.Send("TYPE", key))
}
