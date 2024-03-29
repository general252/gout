package uvptr

import "time"

func Time(v time.Time) *time.Time { return &v }

func Bool(v bool) *bool { return &v }

func Int(v int) *int { return &v }

func Int8(v int8) *int8 { return &v }

func Int16(v int16) *int16 { return &v }

func Int32(v int32) *int32 { return &v }

func Int64(v int64) *int64 { return &v }

func Uint(v uint) *uint { return &v }

func Uint8(v uint8) *uint8 { return &v }

func Uint16(v uint16) *uint16 { return &v }

func Uint32(v uint32) *uint32 { return &v }

func Uint64(v uint64) *uint64 { return &v }

func Float32(v float32) *float32 { return &v }

func Float64(v float64) *float64 { return &v }

func String(v string) *string { return &v }

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TimeV(v *time.Time) time.Time {
	if v == nil {
		return time.Time{}
	}
	return *v
}

func BoolV(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func IntV(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func Int8V(v *int8) int8 {
	if v == nil {
		return 0
	}
	return *v
}

func Int16V(v *int16) int16 {
	if v == nil {
		return 0
	}
	return *v
}

func Int32V(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

func Int64V(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func UintV(v *uint) uint {
	if v == nil {
		return 0
	}
	return *v
}

func Uint8V(v *uint8) uint8 {
	if v == nil {
		return 0
	}
	return *v
}

func Uint16V(v *uint16) uint16 {
	if v == nil {
		return 0
	}
	return *v
}

func Uint32V(v *uint32) uint32 {
	if v == nil {
		return 0
	}
	return *v
}

func Uint64V(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func Float32V(v *float32) float32 {
	if v == nil {
		return 0
	}
	return *v
}

func Float64V(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func StringV(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
