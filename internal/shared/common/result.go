package common_result

type Result[T any] struct {
	Success  bool     `json:"success"`
	Data     *T       `json:"data"`
	Messages []string `json:"messages"`
	Code     int      `json:"code"`
}

func Ok[T any](data *T, code int, message string) Result[T] {
    return Result[T]{
        Success: true,
        Data:    data,
		Messages:  []string{message},
        Code:    code,
    }
}

func Fail[T any](message string, code int) Result[T] {
    return Result[T]{
        Success: false,
        Messages:  []string{message},
        Code:    code,
    }
}

func Fails[T any](messages []string, code int) Result[T] {
    return Result[T]{
        Success: false,
        Messages:  messages,
        Code:    code,
    }
}