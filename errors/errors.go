package errors

const (
	// 定义可预见的异常
	UserNotFound = 10001
	PasswrodErr  = 10002
)

var ErrorCodeMapping = map[int]string{
	UserNotFound: "用户不存在",
}
