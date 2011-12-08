# Redis client for Go

This Redis binding for the Go programming language is still very incomplete.

It will be built bottom up; first a correct and complete implementation of the
Redis protocol and a solid connection context, then pipelining using or
imitating [net/textproto][], conversion of requests and replies to and from
the protocol data types, a method for sending any command and getting
sanitized replies, specific command methods with fixed return types, really
bottom up. All while trying to fully keep up with writing tests, and any
possible changes in Redis or Go itself.

Maybe, one day, it'll make it into the standard package set. Maybe not, that's
okay too.

[net/textproto]: http://golang.org/pkg/net/textproto/
