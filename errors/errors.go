package errors

const (
	// 定义可预见的异常
	UserNotFound = 10001
)

var ErrorCodeMapping = map[int]string{
	UserNotFound: "用户不存在",
}
