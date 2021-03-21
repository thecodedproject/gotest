package gotest

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

func D(v float64) decimal.Decimal {
	return decimal.NewFromFloat(v)
}

func On(funcName string, args ...interface{}) *mock.Call {
	return new(mock.Mock).On(funcName, args)
}
