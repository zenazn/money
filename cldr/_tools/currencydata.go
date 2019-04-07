package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"
)

type stringDate struct {
	t time.Time
}

func (sd *stringDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if t, err := time.Parse("2006-01-02", s); err != nil {
		return err
	} else {
		sd.t = t
		return nil
	}
}

type stringBool struct {
	// Inverted so the default is true
	nb bool
}

func (sb *stringBool) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	var b bool
	if err := json.Unmarshal([]byte(s), &b); err != nil {
		return err
	}
	sb.nb = !b
	return nil
}

type CurrencyDataFile struct {
	Supplemental struct {
		Version struct {
			Number         string `json:"_number"`
			UnicodeVersion string `json:"_unicodeVersion"`
			CLDRVersion    string `json:"_cldrVersion"`
		} `json:"version"`
		CurrencyData struct {
			Fractions map[string]struct {
				Rounding     *int `json:"_rounding,string"`
				Digits       *int `json:"_digits,string"`
				CashRounding *int `json:"_cashRounding,string"`
				CashDigits   *int `json:"_cashDigits,string"`
			} `json:"fractions"`
			Region map[string][]map[string]struct {
				From   *stringDate `json:"_from"`
				To     *stringDate `json:"_to"`
				Tender stringBool  `json:"_tender"`
			} `json:"region"`
		} `json:"currencyData"`
	} `json:"supplemental"`
}

func main() {
	resp, err := http.Get("https://raw.githubusercontent.com/unicode-cldr/cldr-core/master/supplemental/currencyData.json")
	if err != nil {
		fmt.Printf("while downloading currency data from github: %v", err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	var cf CurrencyDataFile
	jd := json.NewDecoder(resp.Body)
	if err := jd.Decode(&cf); err != nil {
		fmt.Printf("while interpreting file: %v", err)
		os.Exit(2)
	}

	f, err := os.Create("data.go")
	if err != nil {
		fmt.Printf("while opening output file=%q for writing: %v", f.Name(), err)
		os.Exit(2)
	}

	w(f, "package cldr")
	w(f, "")
	w(f, `import "time"`)
	w(f, "")
	w(f, "// Generated from CLDR version %s", cf.Supplemental.Version.CLDRVersion)
	w(f, "var CurrencyFractions = map[string]Fractions{")

	default_ := cf.Supplemental.CurrencyData.Fractions["DEFAULT"]
	delete(cf.Supplemental.CurrencyData.Fractions, "DEFAULT")

	defRounding := *default_.Rounding
	defDigits := *default_.Digits

	currencies := make([]string, 0, len(cf.Supplemental.CurrencyData.Fractions))
	for s := range cf.Supplemental.CurrencyData.Fractions {
		currencies = append(currencies, s)
	}
	sort.Strings(currencies)

	for _, s := range currencies {
		fdata := cf.Supplemental.CurrencyData.Fractions[s]
		rounding := defaults(fdata.Rounding, defRounding)
		digits := defaults(fdata.Digits, defDigits)
		w(f, "\t%q: {%d, %d, %d, %d},",
			s,
			rounding, digits,
			defaults(fdata.CashRounding, rounding),
			defaults(fdata.CashDigits, digits),
		)
	}

	w(f, "}")
	w(f, "")
	w(f, "var DefaultFractions = Fractions{%d, %d, %d, %d}", defRounding, defDigits, defRounding, defDigits)
	w(f, "")

	regions := make([]string, 0, len(cf.Supplemental.CurrencyData.Region))
	for r := range cf.Supplemental.CurrencyData.Region {
		regions = append(regions, r)
	}
	sort.Strings(regions)

	w(f, `func dp(s string) *time.Time {`)
	w(f, "\t"+`if t, err := time.Parse("2016-01-02", s); err != nil {`)
	w(f, "\t\t"+`panic("invalid date")`)
	w(f, "\t"+`} else {`)
	w(f, "\t\t"+`return &t`)
	w(f, "\t"+`}`)
	w(f, `}`)
	w(f, "")

	w(f, "var RegionCurrencies = map[string][]RegionalUsage{")
	for _, r := range regions {
		w(f, "\t%q: {", r)
		for _, rrange := range cf.Supplemental.CurrencyData.Region[r] {
			// Should only be one
			for ccy, cdata := range rrange {
				if cdata.To != nil {
					// CLDR date ranges to until midnight at
					// the end of the day
					cdata.To.t = cdata.To.t.Add(24 * time.Hour)
				}
				w(f, "\t\t{%q, %s, %s, %v},",
					ccy, fmtTime(cdata.From), fmtTime(cdata.To), !cdata.Tender.nb)
			}
		}
		w(f, "\t},")
	}
	w(f, "}")
}

func w(f *os.File, s string, args ...interface{}) {
	_, err := fmt.Fprintf(f, s+"\n", args...)
	if err != nil {
		fmt.Printf("while writing %q: %v", f.Name(), err)
		os.Exit(2)
	}
}

func defaults(a *int, b int) int {
	if a == nil {
		return b
	}
	return *a
}

func fmtTime(sd *stringDate) string {
	if sd == nil {
		return "nil"
	}
	return fmt.Sprintf("dp(%q)", sd.t.Format("2006-01-02"))
}
