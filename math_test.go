package money

import (
	"testing"

	"github.com/zenazn/money/currency"
	"github.com/zenazn/money/decimal"
)

func mxn(i int64) Money {
	return FromMinorUnits(i, currency.MXN)
}

var addSubTests = []struct {
	m1, m2, add, sub Money
}{
	{mxn(10), mxn(35), mxn(45), mxn(-25)},
	{mxn(35), mxn(10), mxn(45), mxn(25)},
	{mxn(-10), mxn(35), mxn(25), mxn(-45)},
	{mxn(10), Money{}, mxn(10), mxn(10)},
	{Money{}, mxn(10), mxn(10), mxn(-10)},
}

func TestAddSub(t *testing.T) {
	for i, test := range addSubTests {
		if r := test.m1.Add(test.m2); !r.Eq(test.add) {
			t.Errorf("[%d] add expected %s got %s", i, test.add, r)
		}
		if r, _ := test.m1.AddErr(test.m2); !r.Eq(test.add) {
			t.Errorf("[%d] add err expected %s got %s", i, test.add, r)
		}
		if r := test.m1.Sub(test.m2); !r.Eq(test.sub) {
			t.Errorf("[%d] sub expected %s got %s", i, test.sub, r)
		}
		if r, _ := test.m1.SubErr(test.m2); !r.Eq(test.sub) {
			t.Errorf("[%d] sub err expected %s got %s", i, test.sub, r)
		}
	}
}

func TestMathErrors(t *testing.T) {
	m1 := FromMinorUnits(1234, currency.CAD)
	m2 := FromMinorUnits(1234, currency.JPY)
	if _, err := m1.AddErr(m2); err == nil {
		t.Error("AddE expected error")
	}
	if _, err := m1.SubErr(m2); err == nil {
		t.Error("SubE expected error")
	}
}

var mulDivTests = []struct {
	m        Money
	r        decimal.Rate
	mul, div Money
}{
	{mxn(1230), decimal.NewRate(2000000), mxn(2460), mxn(615)},
	{mxn(-1230), decimal.NewRate(2000000), mxn(-2460), mxn(-615)},
	{mxn(1230), decimal.NewRate(-2000000), mxn(-2460), mxn(-615)},
}

func TestMulDiv(t *testing.T) {
	for i, test := range mulDivTests {
		if r := test.m.Mul(test.r); r != test.mul {
			t.Errorf("[%d] mul expected %s got %s", i, test.mul, r)
		}
		if r := test.m.Div(test.r); r != test.div {
			t.Errorf("[%d] div expected %s got %s", i, test.div, r)
		}
	}
}

func TestNeg(t *testing.T) {
	if !mxn(123).Neg().Eq(mxn(-123)) {
		t.Error("what's going on")
	}
}
