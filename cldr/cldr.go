// Package cldr contains currency data from the Unicode Common Locale Data
// Repository (CLDR).
package cldr

import "time"

type Localization struct {
	DisplayName           string
	DisplayNameCountOne   string
	DisplayNameCountOther string
	Symbol                string
	SymbolAltNarrow       string
	SymbolAltVariant      string
}

type Fractions struct {
	Rounding     int
	Digits       int
	CashRounding int
	CashDigits   int
}

type RegionalUsage struct {
	Symbol string
	From   *time.Time
	To     *time.Time
	Tender bool
}
