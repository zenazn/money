package money

import (
	"testing"

	"github.com/zenazn/money/currency"
)

var onecad = FromMinorUnits(1, currency.CAD)
var tencad = FromMinorUnits(10, currency.CAD)

var comparisonTests = []struct {
	a, b                 Money
	lt, lte, eq, gte, gt bool
}{
	{onecad, tencad, true, true, false, false, false},
	{tencad, onecad, false, false, false, true, true},
	{onecad, onecad, false, true, true, true, false},
}

func TestComparisons(t *testing.T) {
	for i, test := range comparisonTests {
		if r := test.a.Lt(test.b); r != test.lt {
			t.Errorf("[%d] lt expected %v got %v", i, test.lt, r)
		}
		if r, _ := test.a.LtErr(test.b); r != test.lt {
			t.Errorf("[%d] lt err expected %v got %v", i, test.lt, r)
		}
		if r := test.a.Lte(test.b); r != test.lte {
			t.Errorf("[%d] lte expected %v got %v", i, test.lte, r)
		}
		if r, _ := test.a.LteErr(test.b); r != test.lte {
			t.Errorf("[%d] lte err expected %v got %v", i, test.lte, r)
		}

		if r := test.a.Eq(test.b); r != test.eq {
			t.Errorf("[%d] eq expected %v got %v", i, test.eq, r)
		}
		if r, _ := test.a.EqErr(test.b); r != test.eq {
			t.Errorf("[%d] eq err expected %v got %v", i, test.eq, r)
		}

		if r := test.a.Gte(test.b); r != test.gte {
			t.Errorf("[%d] gte expected %v got %v", i, test.gte, r)
		}
		if r, _ := test.a.GteErr(test.b); r != test.gte {
			t.Errorf("[%d] gte err expected %v got %v", i, test.gte, r)
		}
		if r := test.a.Gt(test.b); r != test.gt {
			t.Errorf("[%d] gt expected %v got %v", i, test.gt, r)
		}
		if r, _ := test.a.GtErr(test.b); r != test.gt {
			t.Errorf("[%d] gt err expected %v got %v", i, test.gt, r)
		}
	}
}

func TestComparisonErrors(t *testing.T) {
	m1 := FromMinorUnits(1234, currency.CAD)
	m2 := FromMinorUnits(1234, currency.JPY)
	if _, err := m1.LtErr(m2); err == nil {
		t.Error("LtErr expected error")
	}
	if _, err := m1.LteErr(m2); err == nil {
		t.Error("LteErr expected error")
	}
	if _, err := m1.EqErr(m2); err == nil {
		t.Error("EqErr expected error")
	}
	if _, err := m1.GteErr(m2); err == nil {
		t.Error("GteErr expected error")
	}
	if _, err := m1.GtErr(m2); err == nil {
		t.Error("GtErr expected error")
	}
}
