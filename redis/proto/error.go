package proto

var (
	ErrIsNilBulk            *RedisProtoError = &RedisProtoError{"Bulk object is Nil"}
	ErrIsNilMultiBulk       *RedisProtoError = &RedisProtoError{"Multi-Bulk object is Nil"}
	ErrNotIntegerableObject *RedisProtoError = &RedisProtoError{"Not an Integer or String ObjectKind"}
	ErrNotStringableObject  *RedisProtoError = &RedisProtoError{"Not an Integer, String, Status or Error ObjectKind"}
	ErrNotListObject        *RedisProtoError = &RedisProtoError{"Not a List ObjectKind"}
	ErrInt64ReadSize        *RedisProtoError = &RedisProtoError{"protocol error: integer can't be this big"}
	ErrProtocolError        *RedisProtoError = &RedisProtoError{"protocol error"}
)

// Errors introduced by this package.
type RedisProtoError struct {
	ErrorString string
}

func (err *RedisProtoError) String() string { return "redis/proto: " + err.ErrorString }
