# Godis

A simple Redis client for Go.

## todo

*. Refactor Client/PipeClient so that New() returns a SyncClient struct which
has the Methods of Client. PipeClient should then have an a pointer to a Client
struct, but not SyncClient struct. This should make it possible to use regular
commands on either Sync- or PipeClient

*. PipeClient logic is not safe at all. If an error occurs and the user
continues to call ReadReply after that, a new connection will be poped and we
will try to continue reading, making the client hang/timeout ... repeat.
Read/Write logic should follow more of hiredis's way for the CommandAppend().

* Add all tests for Set. Add missing tests for Hash.
