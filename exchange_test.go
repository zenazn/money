package money

import (
	"fmt"
	"testing"

	"github.com/zenazn/money/currency"
	"github.com/zenazn/money/decimal"
)

func mkex(ppm int64, d currency.Currency) ExchangeRate {
	return ExchangeRate{
		currency.USD,
		d,
		decimal.NewRate(ppm),
	}
}

func usd(cents int64) Money {
	return FromMinorUnits(cents, currency.USD)
}

type bitcoin struct{}

func (b bitcoin) Symbol() string {
	return "XBT"
}
func (b bitcoin) Units() currency.Units {
	return currency.Units{8, 8}
}

type precise uint8

func (p precise) Symbol() string {
	return fmt.Sprintf("XPC[%d]", p)
}
func (p precise) Units() currency.Units {
	return currency.Units{2, uint8(p)}
}

func btc(sat int64) Money {
	return FromMinorUnits(sat, bitcoin{})
}

var exchangeTests = []struct {
	e    ExchangeRate
	s, d Money
}{
	{mkex(190, bitcoin{}), usd(123456), btc(23456600)},
	{mkex(1829181, bitcoin{}), usd(829171310), btc(1516704405997100)},
	{mkex(1000000, precise(24)), usd(100), FromMinorUnits(100, precise(24))},
	{ExchangeRate{precise(24), currency.USD, decimal.NewRate(1000000)}, FromMinorUnits(100, precise(24)), usd(100)},
}

func TestExchange(t *testing.T) {
	for i, test := range exchangeTests {
		if r := test.s.Exchange(test.e); !r.Eq(test.d) {
			t.Errorf("[%d] exchange expected %s, got %s", i, test.d, r)
		}
	}
}

func TestExchangeRateGetters(t *testing.T) {
	ex := mkex(123400, currency.EUR)
	if src := ex.Source(); src != currency.USD {
		t.Errorf("source %s", src.Symbol())
	}
	if dst := ex.Destination(); dst != currency.EUR {
		t.Errorf("destination %s", dst.Symbol())
	}
	if r := ex.Rate(); r != decimal.NewRate(123400) {
		t.Errorf("rate %v", r)
	}
}
