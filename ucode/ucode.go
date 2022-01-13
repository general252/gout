package ucode

type Code uint32

const (
	OK                Code = 0
	Unknown           Code = 1
	Fail              Code = 2
	InvalidArgument   Code = 3
	IO                Code = 4
	DB                Code = 5
	DataFormat        Code = 6 // 文件格式错误
	NotFound          Code = 7
	AlreadyExists     Code = 8  // 已经存在
	Canceled          Code = 9  // 取消
	PermissionDenied  Code = 10 // 没有权限
	ResourceExhausted Code = 11 // 资源耗尽
	Aborted           Code = 12 // 终止
	OutOfRange        Code = 13 // 越界
	Unimplemented     Code = 14 // 未实现
	Internal          Code = 15 // 内部错误
	Unavailable       Code = 16 // 无效
	Unauthenticated   Code = 17 // 未验证
	AuthenticatedFail Code = 18 // 验证失败
)

var _s = map[Code]string{
	OK:                "OK",
	Unknown:           "Unknown",
	Fail:              "Fail",
	InvalidArgument:   "InvalidArgument",
	IO:                "IO",
	DB:                "DB",
	DataFormat:        "DataFormat",
	NotFound:          "NotFound",
	AlreadyExists:     "AlreadyExists",
	Canceled:          "Canceled",
	PermissionDenied:  "PermissionDenied",
	ResourceExhausted: "ResourceExhausted",
	Aborted:           "Aborted",
	OutOfRange:        "OutOfRange",
	Unimplemented:     "Unimplemented",
	Internal:          "Internal",
	Unavailable:       "Unavailable",
	Unauthenticated:   "Unauthenticated",
	AuthenticatedFail: "AuthenticatedFail",
}

func (c Code) String() string {
	if str, ok := _s[c]; ok {
		return str
	}

	return _s[Unknown]
}
