# godis

A simple Redis client for Go.

## todo

*. Write documentation and add some examples.

* Add all tests for sorted set and some server stuff.

*. PipeClient logic is not safe at all. If an error occurs and the user
continues to call ReadReply after that, a new connection will be poped and we
will try to continue reading, making the client hang/timeout ... repeat.

## acknowledgment

The work on this client started as I was hacking around on Michael Hoisie's
original redis client for Go. Also the recent work done by Frank Müller on his
client gave me some pointers how to better handle return values. 
