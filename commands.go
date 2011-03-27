package godis 

// generic

// Delete a key
func (c *Client) Del(key ...string) int64 {
    return 0
}

// Determine if a key exists
func (c *Client) Exists(key string) int64 {
    return 0
}

// Set a key's time to live in seconds
func (c *Client) Expire(key string, seconds int) int64 {
    return 0
}

// Set the expiration for a key as a UNIX timestamp
func (c *Client) Expireat(key string, timestamp int64) int64 {
    return 0
}

// Find all keys matching the given pattern
func (c *Client) Keys(pattern string) [][]byte {
    return [][]byte{}
}

// Move a key to another database
func (c *Client) Move(key string, db int) int64 {
    return 0
}

// Remove the expiration from a key
func (c *Client) Persist(key string) int64 {
    return 0
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
func (c *Client) Ttl(key string) int64 {
    return 0
}

// Determine the type stored at key
func (c *Client) Type(key string) bool {
    return false
}

