package customerrors

type CustomError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *CustomError) Error() string {
	return e.Message
}
