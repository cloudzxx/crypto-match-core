package decimal

import (
	"github.com/shopspring/decimal"
)

func NewFromFloat64(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

func Add(x, y decimal.Decimal) decimal.Decimal {
	return x.Add(y)
}

func Sub(x, y decimal.Decimal) decimal.Decimal {
	return x.Sub(y)
}

func Mul(x, y decimal.Decimal) decimal.Decimal {
	return x.Mul(y)
}

func Div(x, y decimal.Decimal, prec uint) decimal.Decimal {
	if prec == 0 {
		prec = 12
	}
	return x.DivRound(y, int32(prec))
}

func Cmp(x, y decimal.Decimal) int {
	return x.Cmp(y)
}

func IsZero(x decimal.Decimal) bool {
	return x.IsZero()
}

func IsPositive(x decimal.Decimal) bool {
	return x.GreaterThan(decimal.Zero)
}

func IsNegative(x decimal.Decimal) bool {
	return x.LessThan(decimal.Zero)
}

func Min(x, y decimal.Decimal) decimal.Decimal {
	if x.Cmp(y) > 0 {
		return y
	}
	return x
}

func Max(x, y decimal.Decimal) decimal.Decimal {
	if x.Cmp(y) < 0 {
		return y
	}
	return x
}

func MustNewFromString(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}
