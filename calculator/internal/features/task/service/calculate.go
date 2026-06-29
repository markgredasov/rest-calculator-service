package tasks_service

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_errors "github.com/markgredasov/rest-calculator-service/internal/core/errors"
	core_logger "github.com/markgredasov/rest-calculator-service/internal/core/logger"
	"go.uber.org/zap"
)

func (s *TasksService) Calculate(ctx context.Context, task domain.Task) (domain.Task, error) {
	log := core_logger.FromContext(ctx)

	if task.Expression == "" {
		return domain.Task{}, fmt.Errorf("empty expression: %w", core_errors.ErrInvalidArgument)
	}

	result, err := evaluateWithGoParser(task.Expression)
	if err != nil {
		log.Error("failed to evaluate expression", zap.Error(err))
		return domain.Task{}, fmt.Errorf("failed to evaluate expression: %w", err)
	}

	resultTask := domain.Task{
		Expression: task.Expression,
		Result:     result,
	}

	return resultTask, nil
}

func evaluateWithGoParser(expr string) (float64, error) {
	expr = strings.ReplaceAll(expr, " ", "")

	for _, ch := range expr {
		if !isValidMathChar(ch) {
			return 0, fmt.Errorf("invalid character '%c' in expression", ch)
		}
	}

	node, err := parser.ParseExpr(expr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse expression: %w", err)
	}

	return evalNode(node)
}

func isValidMathChar(ch rune) bool {
	return (ch >= '0' && ch <= '9') ||
		ch == '+' || ch == '-' || ch == '*' || ch == '/' ||
		ch == '(' || ch == ')' || ch == '.'
}

func evalNode(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		left, err := evalNode(n.X)
		if err != nil {
			return 0, err
		}
		right, err := evalNode(n.Y)
		if err != nil {
			return 0, err
		}

		switch n.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("unsupported operator: %v", n.Op)
		}

	case *ast.ParenExpr:
		return evalNode(n.X)

	case *ast.BasicLit:
		if n.Kind == token.INT || n.Kind == token.FLOAT {
			val, err := strconv.ParseFloat(n.Value, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse number: %s", n.Value)
			}
			return val, nil
		}
		return 0, fmt.Errorf("unsupported literal: %v", n.Kind)

	case *ast.UnaryExpr:
		if n.Op == token.SUB {
			val, err := evalNode(n.X)
			if err != nil {
				return 0, err
			}
			return -val, nil
		}
		return 0, fmt.Errorf("unsupported unary operator: %v", n.Op)

	default:
		return 0, fmt.Errorf("unsupported node type: %T", n)
	}
}
