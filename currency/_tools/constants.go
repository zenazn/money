package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/zenazn/money/iso4217"
)

func main() {
	// currency codes are decently dense in [0, 1000) and currency symbols
	// are constant-width, so store the ISO number-to-symbol mapping in a
	// big string
	var lut [3000]byte
	var currencies []string
	numberMap := make(map[string]int)
	var minor [1000]uint8

	for i := 0; i < len(minor); i++ {
		minor[i] = 0xff
	}

	for _, ce := range iso4217.Data {
		n := ce.CurrencyNumber
		if lut[n*3] != 0 {
			// Duplicate: we've already included this currency
			continue
		}
		lut[n*3+0] = ce.Currency[0]
		lut[n*3+1] = ce.Currency[1]
		lut[n*3+2] = ce.Currency[2]

		currencies = append(currencies, ce.Currency)
		numberMap[ce.Currency] = n
		minor[n] = uint8(ce.CurrencyMinorUnits)
		// The distinction between -1 and 0 isn't super important here
		if ce.CurrencyMinorUnits == -1 {
			minor[n] = 0
		}
	}
	sort.Strings(currencies)

	for i := 0; i < len(lut); i++ {
		if lut[i] == 0 {
			lut[i] = '.'
		}
	}

	f, err := os.Create("constants.go")
	if err != nil {
		fmt.Printf("while creating constants.go: %v", err)
		os.Exit(2)
	}

	w(f, "package currency")
	w(f, "")
	w(f, "const (")
	for _, currency := range currencies {
		w(f, "\t%s iso = %d", currency, numberMap[currency])
	}
	w(f, ")")
	w(f, "")

	w(f, "const numberToSymbol = %q", lut)
	w(f, "")
	w(f, "var symbolToCurrency = map[string]iso{")
	for _, currency := range currencies {
		w(f, "\t%q: %d,", currency, numberMap[currency])
	}
	w(f, "}")
	w(f, "")

	w(f, "var numberToMinor = %#v", minor)
}

func w(f *os.File, s string, args ...interface{}) {
	_, err := fmt.Fprintf(f, s+"\n", args...)
	if err != nil {
		fmt.Printf("while writing %q: %v\n", f.Name(), err)
		os.Exit(2)
	}
}
