package ptr

import "time"

var (
	NilBool    *bool
	NilInt     *int
	NilFloat64 *float64
	NilString  *string
	NilTime    *time.Time
)

func BoolPtr(b bool) *bool {
	return &b
}

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func StrPtr(s string) *string {
	return &s
}
