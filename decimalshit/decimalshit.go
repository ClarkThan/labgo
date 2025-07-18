package decimalshit

import (
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	Pivot = decimal.NewFromFloat(0.015) // 误差比例
)

func demo1() {
	amount1, _ := decimal.NewFromString("2547.8681")
	amount2, _ := decimal.NewFromString("2515.1463")

	diffRate := amount1.Sub(amount2).DivRound(amount1, 6)
	ok := diffRate.LessThanOrEqual(Pivot)
	fmt.Println(ok)
}

func Main() {
	demo1()
}
