package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

type CurrenciesFile struct {
	Main map[string]struct {
		Identity struct {
			Version struct {
				Number      string `json:"_number"`
				CLDRVersion string `json:"_cldrVersion"`
			} `json:"version"`
		} `json:"identity"`
		Numbers struct {
			Currencies map[string]struct {
				DisplayName           string `json:"displayName"`
				DisplayNameCountOne   string `json:"displayName-count-one"`
				DisplayNameCountOther string `json:"displayName-count-other"`
				Symbol                string `json:"symbol"`
				SymbolAltNarrow       string `json:"symbol-alt-narrow"`
				SymbolAltVariant      string `json:"symbol-alt-variant"`
			} `json:"currencies`
		} `json:"numbers"`
	} `json:"main"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s locale\n", os.Args[0])
		os.Exit(1)
	}
	locale := os.Args[1]

	url := fmt.Sprintf("https://raw.githubusercontent.com/unicode-cldr/cldr-numbers-modern/master/main/%s/currencies.json", locale)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("while downloading locale=%q from github: %v\n", locale, err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	var cf CurrenciesFile
	jd := json.NewDecoder(resp.Body)
	if err := jd.Decode(&cf); err != nil {
		fmt.Printf("while interpreting file: %v\n", err)
		os.Exit(2)
	}

	fname := fmt.Sprintf("%s.go", locale)
	f, err := os.Create(fname)
	if err != nil {
		fmt.Printf("while opening output file=%q for writing: %v\n", fname, err)
		os.Exit(2)
	}

	w(f, "package cldr")
	w(f, "")
	w(f, "// Generated from CLDR version %s", cf.Main[locale].Identity.Version.CLDRVersion)
	w(f, "var Locale%s = map[string]Localization{", strings.ToUpper(locale))

	currencies := make([]string, 0, len(cf.Main[locale].Numbers.Currencies))
	for s := range cf.Main[locale].Numbers.Currencies {
		currencies = append(currencies, s)
	}
	sort.Strings(currencies)

	for _, s := range currencies {
		cdata := cf.Main[locale].Numbers.Currencies[s]
		w(f, "\t%q: {%q, %q, %q, %q, %q, %q},",
			s,
			cdata.DisplayName,
			cdata.DisplayNameCountOne,
			cdata.DisplayNameCountOther,
			cdata.Symbol,
			cdata.SymbolAltNarrow,
			cdata.SymbolAltVariant,
		)
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
