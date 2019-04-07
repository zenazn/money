package currency

import "testing"

func TestSymbol(t *testing.T) {
	if s := USD.Symbol(); s != "USD" {
		t.Errorf("USD is %q?!", s)
	}
	if s := EUR.Symbol(); s != "EUR" {
		t.Errorf("EUR is %q?!", s)
	}
}

func TestFromISO(t *testing.T) {
	if c, err := FromISOSymbol("USD"); err != nil || c != USD {
		t.Errorf("USD is %#v, not %#v?! %v", c, USD, err)
	}
	if c, err := FromISOSymbol("EUR"); err != nil || c != EUR {
		t.Errorf("EUR is %#v, not %#v?! %v", c, EUR, err)
	}
	if c, err := FromISOSymbol("bitcoin"); err == nil || c != nil {
		t.Errorf("bitcoin is not a currency! %v, %#v", err, c)
	}
}
