package app

type Code int

// Don’t change the numbers, as external scripts could rely on them.
const (
	GENERAL_ERROR             Code = 1
	NO_INPUT_ERROR            Code = 2
	NO_TARGET_FILE            Code = 3
	IO_ERROR                  Code = 4
	BOOKMARK_ACCESS_ERROR     Code = 5
	BOOKMARK_NOT_SET          Code = 6
	NO_SUCH_FILE_OR_DIRECTORY Code = 7
)

func (c Code) ToInt() int {
	return int(c)
}

type Error interface {
	Error() string
	Details() string
	Original() error
	Code() Code
}

type AppError struct {
	code     Code
	message  string
	details  string
	original error
}

func NewError(message string, details string, original error) Error {
	return NewErrorWithCode(GENERAL_ERROR, message, details, original)
}

func NewErrorWithCode(code Code, message string, details string, original error) Error {
	return AppError{code, message, details, original}
}

func (e AppError) Error() string {
	return e.message
}

func (e AppError) Details() string {
	return e.details
}

func (e AppError) Original() error {
	return e.original
}

func (e AppError) Code() Code {
	return e.code
}
