package calculator_service

type CalculatorService struct {
	calculatorRepository CalculatorRepository
}

type CalculatorRepository interface{}

func NewCalculatorService(repo CalculatorRepository) CalculatorService {
	return CalculatorService{
		calculatorRepository: repo,
	}
}
