package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ISO4217 struct {
	Published string `xml:"Pblshd,attr"`
	Entries   []struct {
		CountryName        string `xml:"CtryNm"`
		CurrencyName       string `xml:"CcyNm"`
		Currency           string `xml:"Ccy"`
		CurrencyNumber     int    `xml:"CcyNbr"`
		CurrencyMinorUnits string `xml:"CcyMnrUnts"`
	} `xml:"CcyTbl>CcyNtry"`
}

type currency struct {
	name       string
	number     int
	minorUnits int
}

func main() {
	resp, err := http.Get("https://www.currency-iso.org/dam/downloads/lists/list_one.xml")
	if err != nil {
		fmt.Printf("while downloading currency data from ISO: %v\n", err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	var isof ISO4217
	xd := xml.NewDecoder(resp.Body)
	if err := xd.Decode(&isof); err != nil {
		fmt.Printf("while interpreting file: %v\n", err)
		os.Exit(2)
	}

	f, err := os.Create("data.go")
	if err != nil {
		fmt.Printf("while opening data.go: %v\n", err)
		os.Exit(2)
	}

	w(f, "package iso4217")
	w(f, "")
	w(f, "// Data from ISO 4217, published on %s", isof.Published)
	w(f, "")
	w(f, "type CurrencyEntry struct {")
	w(f, "\tCountryName        string")
	w(f, "\tCurrencyName       string")
	w(f, "\tCurrency           string")
	w(f, "\tCurrencyNumber     int")
	w(f, "\tCurrencyMinorUnits int")
	w(f, "}")
	w(f, "")

	w(f, "var Data = []CurrencyEntry{")

	for i, e := range isof.Entries {
		if e.Currency == "" {
			// ISO 4217 has an entry for e.g., Antarctica, which
			// doesn't have a currency
			continue
		}

		// ISO decided to include a non-breaking space at the end of the
		// country name for XDR. I don't think it's supposed to be there
		e.CountryName = strings.TrimSpace(e.CountryName)

		// We encode currencies for which minor units don't apply (e.g.,
		// XAU: troy ounces of gold) as having minor units of -1
		munits := -1
		if e.CurrencyMinorUnits != "N.A." {
			n, err := strconv.Atoi(e.CurrencyMinorUnits)
			if err != nil {
				fmt.Printf("currency %s (%d) unexpected minor units: %v\n", e.Currency, i, err)
				os.Exit(2)
			}
			munits = n
		}

		w(f, "\t{%q, %q, %q, %d, %d},",
			e.CountryName, e.CurrencyName, e.Currency, e.CurrencyNumber, munits)
	}

	w(f, "}")
}

func w(f *os.File, s string, args ...interface{}) {
	_, err := fmt.Fprintf(f, s+"\n", args...)
	if err != nil {
		fmt.Printf("while writing %q: %v\n", f.Name(), err)
		os.Exit(2)
	}
}
