package ptr

import (
	"testing"
	"time"
)

func TestPtr(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		s := "foo"
		ptr := Ptr(s)
		if *ptr != s {
			t.Errorf("got %v, want %v", *ptr, s)
		}
	})

	t.Run("bool", func(t *testing.T) {
		b := true
		ptr := Ptr(b)
		if *ptr != b {
			t.Errorf("got %v, want %v", *ptr, b)
		}
	})

	t.Run("int", func(t *testing.T) {
		i := 1
		ptr := Ptr(i)
		if *ptr != i {
			t.Errorf("got %v, want %v", *ptr, i)
		}
	})
	t.Run("int8", func(t *testing.T) {
		i := int8(1)
		ptr := Ptr(i)
		if *ptr != i {
			t.Errorf("got %v, want %v", *ptr, i)
		}
	})
	t.Run("int16", func(t *testing.T) {
		i := int16(1)
		ptr := Ptr(i)
		if *ptr != i {
			t.Errorf("got %v, want %v", *ptr, i)
		}
	})
	t.Run("int32", func(t *testing.T) {
		i := int32(1)
		ptr := Ptr(i)
		if *ptr != i {
			t.Errorf("got %v, want %v", *ptr, i)
		}
	})
	t.Run("int64", func(t *testing.T) {
		i := int64(1)
		ptr := Ptr(i)
		if *ptr != i {
			t.Errorf("got %v, want %v", *ptr, i)
		}
	})

	t.Run("uint", func(t *testing.T) {
		u := uint(1)
		ptr := Ptr(u)
		if *ptr != u {
			t.Errorf("got %v, want %v", *ptr, u)
		}
	})
	t.Run("uint8", func(t *testing.T) {
		u := uint8(1)
		ptr := Ptr(u)
		if *ptr != u {
			t.Errorf("got %v, want %v", *ptr, u)
		}
	})
	t.Run("uint16", func(t *testing.T) {
		u := uint16(1)
		ptr := Ptr(u)
		if *ptr != u {
			t.Errorf("got %v, want %v", *ptr, u)
		}
	})
	t.Run("uint32", func(t *testing.T) {
		u := uint32(1)
		ptr := Ptr(u)
		if *ptr != u {
			t.Errorf("got %v, want %v", *ptr, u)
		}
	})
	t.Run("uint64", func(t *testing.T) {
		u := uint64(1)
		ptr := Ptr(u)
		if *ptr != u {
			t.Errorf("got %v, want %v", *ptr, u)
		}
	})

	t.Run("float64", func(t *testing.T) {
		f := float64(1.23456789)
		ptr := Ptr(f)
		if *ptr != f {
			t.Errorf("got %v, want %v", *ptr, f)
		}
	})

	t.Run("float32", func(t *testing.T) {
		f := float32(1.23456789)
		ptr := Ptr(f)
		if *ptr != f {
			t.Errorf("got %v, want %v", *ptr, f)
		}
	})

	t.Run("time.Time", func(t *testing.T) {
		time := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		ptr := Ptr(time)
		if *ptr != time {
			t.Errorf("got %v, want %v", *ptr, time)
		}
	})
}
