# Godis

A simple Redis client for Go.

## todo

*. PipeClient logic is not safe at all. If an error occurs and the user
continues to call ReadReply after that, a new connection will be poped and we
will try to continue reading, making the client hang/timeout ... repeat.
Read/Write logic should follow more of hiredis's way for the CommandAppend().

* Add all tests for sorted set and some server stuff.
