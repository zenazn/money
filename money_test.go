package money

import (
	"testing"

	"github.com/zenazn/money/currency"
	"github.com/zenazn/money/decimal"
)

func mustparse(amt, ccy string) Money {
	m, err := Parse(amt, ccy)
	if err != nil {
		panic(err)
	}
	return m
}

var constructorTests = []struct {
	v   Money
	amt string
	ccy currency.Currency
}{
	{Money{}, "0", nil},
	{Zero(currency.USD), "0", currency.USD},
	{New(decimal.FromI64(42), currency.EUR), "42", currency.EUR},
	{FromMinorUnits(123, currency.USD), "1230000", currency.USD},
	{FromMinorUnits(123, currency.JPY), "123000000", currency.JPY},
	{mustparse("1.23", "USD"), "1230000", currency.USD},
	{mustparse("0.000023", "USD"), "23", currency.USD},
	{mustparse("1234", "USD"), "1234000000", currency.USD},
}

func TestConstructors(t *testing.T) {
	for i, test := range constructorTests {
		if s := test.v.amt.String(); s != test.amt {
			t.Errorf("[%d] amount expected %q got %q", i, test.amt, s)
		}
		if c := test.v.ccy; c != test.ccy {
			t.Errorf("[%d] amount expected %v got %v", i, test.ccy, c)
		}
	}
}

var parseFailureTests = []string{
	".123",
	"123.",
	"123.4.2",
	"123,45",
	"orange",
	"12.3456789",
}

func TestParseFailures(t *testing.T) {
	for i, s := range parseFailureTests {
		if v, err := Parse(s, "USD"); err == nil {
			t.Errorf("[%d] unexpectedly passed: %v", i, v)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Parse("1234.56", "USD")
	}
}

func BenchmarkParseLarge(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Parse("26499352303014264292828635027055.993969", "USD")
	}
}

var basicTests = []struct {
	m    Money
	ccy  currency.Currency
	zero bool
}{
	{FromMinorUnits(123, currency.CAD), currency.CAD, false},
	{Money{}, nil, true},
	{Zero(currency.CAD), currency.CAD, true},
}

func TestBasic(t *testing.T) {
	for i, test := range basicTests {
		if r := (test.m.Amount() == decimal.Decimal{}); r != test.zero {
			t.Errorf("[%d] amount", i)
		}
		if r := test.m.Zero(); r != test.zero {
			t.Errorf("[%d] zero", i)
		}
		if r := test.m.Currency(); r != test.ccy {
			t.Errorf("[%d] ccy", i)
		}
	}
}

type fakeUSD struct{}

func (f fakeUSD) Symbol() string        { return currency.USD.Symbol() }
func (f fakeUSD) Units() currency.Units { return currency.USD.Units() }

func TestComparableTo(t *testing.T) {
	m1 := FromMinorUnits(123, currency.USD)
	m2 := FromMinorUnits(567, currency.USD)
	m3 := FromMinorUnits(890, currency.EUR)
	m4 := Money{}
	m5 := Zero(fakeUSD{})

	if !m1.ComparableTo(m2) {
		t.Error("m1 / m2")
	}
	if m1.ComparableTo(m3) {
		t.Error("m1 / m3")
	}
	if !m1.ComparableTo(m4) {
		t.Error("m1 / m4")
	}
	if !m1.ComparableTo(m5) {
		t.Error("m1 / m5")
	}
}

var stringTests = []struct {
	v Money
	s string
}{
	{Money{}, "0"},
	{Zero(currency.USD), "USD 0.00"},
	{FromMinorUnits(12300, currency.EUR), "EUR 123.00"},
	{FromMinorUnits(1230, currency.EUR), "EUR 12.30"},
	{FromMinorUnits(123, currency.EUR), "EUR 1.23"},
	{FromMinorUnits(23, currency.EUR), "EUR 0.23"},
	{FromMinorUnits(3, currency.EUR), "EUR 0.03"},
	{FromMinorUnits(-1234, currency.MXN), "MXN -12.34"},
	{New(decimal.FromI64(42), currency.CAD), "CAD 0.000042"},
	{New(decimal.FromI64(12345678), currency.CAD), "CAD 12.345678"},
	{New(decimal.FromI64(12345600), currency.CAD), "CAD 12.3456"},
	{FromMinorUnits(123, currency.JPY), "JPY 123"},
	{New(decimal.FromI64(12300000), currency.JPY), "JPY 12.3"},
	{New(decimal.FromI64(12000000), currency.JPY), "JPY 12"},
}

func TestString(t *testing.T) {
	for i, test := range stringTests {
		if s := test.v.String(); s != test.s {
			t.Errorf("[%d] expected %q got %q", i, test.s, s)
		}
	}
}
