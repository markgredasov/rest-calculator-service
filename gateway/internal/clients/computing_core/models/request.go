package clients_computing_core_models

type Request struct {
	Numbers   []int  `json:"numbers"`
	Operation string `json:"operation"`
}

func NewRequest(numbers []int, operation string) Request {
	return Request{
		Numbers:   numbers,
		Operation: operation,
	}
}
