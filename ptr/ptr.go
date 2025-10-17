package ptr

import "time"

var (
	NilBool *bool

	NilInt   *int
	NilInt8  *int8
	NilInt16 *int16
	NilInt32 *int32
	NilInt64 *int64

	NilUint   *uint
	NilUint8  *uint8
	NilUint16 *uint16
	NilUint32 *uint32
	NilUint64 *uint64

	NilFloat32 *float32
	NilFloat64 *float64
	NilString  *string
	NilTime    *time.Time
)

func Ptr[T any](v T) *T {
	return &v
}
