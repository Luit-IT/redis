package proto

import "os"

var (
	ErrIsNilBulk            os.Error = &RedisProtocolError{"Bulk object is Nil"}
	ErrIsNilMultiBulk       os.Error = &RedisProtocolError{"Multi-Bulk object is Nil"}
	ErrNotIntegerableObject os.Error = &RedisProtocolError{"Not an Integer or String ObjectKind"}
	ErrNotStringableObject  os.Error = &RedisProtocolError{"Not an Integer, String, Status or Error ObjectKind"}
	ErrNotListObject        os.Error = &RedisProtocolError{"Not a List ObjectKind"}
	ErrInt64ReadSize         os.Error = &RedisProtocolError{"integer can't be this big"}
	ErrNestedMultiBulk        os.Error = &RedisProtocolError{"nested/unexpected multi-bulk"}
	ErrProtocolError        os.Error = &RedisProtocolError{"protocol error"}
)

// Errors introduced by this package.
type RedisProtocolError struct {
	ErrorString string
}

func (err *RedisProtocolError) String() string { return "redis/proto: " + err.ErrorString }
