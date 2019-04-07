package decimal

import "testing"

var addTests = []struct {
	a, b, c Decimal
}{
	{FromI64(2500000), FromI64(5000000), FromI64(7500000)},
	{FromI64(-2500000), FromI64(5000000), FromI64(2500000)},
	{FromI64(-2500000), FromI64(1700000), FromI64(-800000)},
	{FromI64(1700000), FromI64(-1700000), FromI64(0)},
	{Decimal{1, 0xfffffffffffffff0}, FromI64(17), Decimal{2, 1}},
}

func TestAdd(t *testing.T) {
	for i, v := range addTests {
		if o := v.a.Add(v.b); o != v.c {
			t.Errorf("[%d] %#v != %#v", i, o, v.c)
		}
	}
}

var subTests = []struct {
	a, b, c Decimal
}{
	{FromI64(2500000), FromI64(5000000), FromI64(-2500000)},
	{FromI64(-2500000), FromI64(5000000), FromI64(-7500000)},
	{FromI64(1700000), FromI64(2500000), FromI64(-800000)},
	{FromI64(1700000), FromI64(1700000), FromI64(0)},
	{Decimal{2, 0}, FromI64(17), Decimal{1, 0xffffffffffffffef}},
}

func TestSub(t *testing.T) {
	for i, v := range subTests {
		if o := v.a.Sub(v.b); o != v.c {
			t.Errorf("[%d] %#v != %#v", i, o, v.c)
		}
	}
}

var mulTests = []struct {
	a Decimal
	r Rate
	c Decimal
}{
	{FromI64(25000000), NewRate(25), FromI64(625)},
	{FromI64(2500000), NewRate(25), FromI64(62)},
	{FromI64(250000), NewRate(25), FromI64(6)},
	{FromI64(25000000), NewRate(27), FromI64(675)},
	{FromI64(2500000), NewRate(27), FromI64(68)},
	{FromI64(250000), NewRate(27), FromI64(7)},
	{Decimal{2, 0}, NewRate(333333), Decimal{0, 0xaaaa9f7b5aea3162}},
	{Decimal{2, 1}, NewRate(333333), Decimal{0, 0xaaaa9f7b5aea3162}},
	{Decimal{2, 2}, NewRate(333333), Decimal{0, 0xaaaa9f7b5aea3162}},
	{Decimal{5, 0}, NewRate(333333), Decimal{1, 0xaaaa8eb463497b74}},
	{FromI64(25000000), NewRate(-25), FromI64(-625)},
	{FromI64(-25000000), NewRate(25), FromI64(-625)},
	{FromI64(-25000000), NewRate(-25), FromI64(625)},
	{FromI64(-2500000), NewRate(25), FromI64(-62)},
	{FromI64(-2500000), NewRate(27), FromI64(-68)},
}

func TestMul(t *testing.T) {
	for i, v := range mulTests {
		if o := v.a.Mul(v.r); o != v.c {
			t.Errorf("[%d] %#v != %#v", i, o, v.c)
		}
	}
}

var divTests = []struct {
	a Decimal
	r Rate
	c Decimal
}{
	{FromI64(25000000), NewRate(1000000), FromI64(25000000)},
	{FromI64(25000000), NewRate(100000), FromI64(250000000)},
	{FromI64(25000000), NewRate(10000000), FromI64(2500000)},
	{FromI64(83000000), NewRate(333333), FromI64(249000249)},
	{FromI64(83000000), NewRate(3333333), FromI64(24900002)},
	{FromI64(83000000), NewRate(3000000), FromI64(27666667)},
	{FromI64(83000000), NewRate(-3000000), FromI64(-27666667)},
	{FromI64(-83000000), NewRate(3000000), FromI64(-27666667)},
	{FromI64(-83000000), NewRate(-3000000), FromI64(27666667)},
	{Decimal{0x1234, 0x5}, NewRate(3832922), Decimal{0x4bf, 0xc85a9d723ac83f4c}},
	{Decimal{0x1234, 0x6}, NewRate(3832922), Decimal{0x4bf, 0xc85a9d723ac83f4c}},
	{Decimal{0x1234, 0x7}, NewRate(3832922), Decimal{0x4bf, 0xc85a9d723ac83f4d}},
	{FromI64(25), NewRate(10000000), FromI64(2)},
	{FromI64(75), NewRate(10000000), FromI64(8)},
	{FromI64(75), NewRate(2000001), FromI64(37)},
	{FromI64(75), NewRate(1999999), FromI64(38)},
	{FromI64(3), NewRate(2000001), FromI64(1)},
	{FromI64(3), NewRate(1999999), FromI64(2)},
}

func TestDiv(t *testing.T) {
	for i, v := range divTests {
		if o := v.a.Div(v.r); o != v.c {
			t.Errorf("[%d] %#v != %#v", i, o, v.c)
		}
	}
}

var comparisonTests = []struct {
	a, b                 Decimal
	lt, lte, eq, gte, gt bool
}{
	{FromI64(-10), FromI64(10), true, true, false, false, false},
	{FromI64(10), FromI64(-10), false, false, false, true, true},

	{Decimal{0xf000000000000000, 0}, FromI64(-10), true, true, false, false, false},
	{FromI64(-10), FromI64(-10), false, true, true, true, false},
	{FromI64(-10), Decimal{0xf000000000000000, 0}, false, false, false, true, true},

	{FromI64(10), Decimal{1, 0}, true, true, false, false, false},
	{FromI64(10), FromI64(10), false, true, true, true, false},
	{Decimal{1, 0}, FromI64(10), false, false, false, true, true},
}

func TestComparisons(t *testing.T) {
	for i, test := range comparisonTests {
		if r := test.a.Lt(test.b); r != test.lt {
			t.Errorf("[%d] lt expected %v got %v", i, test.lt, r)
		}
		if r := test.a.Lte(test.b); r != test.lte {
			t.Errorf("[%d] lte expected %v got %v", i, test.lte, r)
		}
		if r := test.a.Eq(test.b); r != test.eq {
			t.Errorf("[%d] eq expected %v got %v", i, test.eq, r)
		}
		if r := test.a.Gte(test.b); r != test.gte {
			t.Errorf("[%d] gte expected %v got %v", i, test.gte, r)
		}
		if r := test.a.Gt(test.b); r != test.gt {
			t.Errorf("[%d] gt expected %v got %v", i, test.gt, r)
		}
	}
}

var stringTests = []struct {
	a Decimal
	s string
}{
	{FromI64(25000000), "25000000"},
	{FromI64(-25000000), "-25000000"},
	{FromI64(8740302187228643401), "8740302187228643401"},
	{FromI64(-8740302187228643401), "-8740302187228643401"},
	{Decimal{0x1234, 0x5}, "85961827383486510530565"},
	{Decimal{0x4b3b4ca85a86c47a, 0x098a223fffffffff}, "99999999999999999999999999999999999999"},
	{Decimal{0x4b3b4ca85a86c47a, 0x098a224000000000}, "100000000000000000000000000000000000000"},
	{Decimal{0xb4c4b357a5793b85, 0xf675ddc000000001}, "-99999999999999999999999999999999999999"},
	{Decimal{0xb4c4b357a5793b85, 0xf675ddc000000000}, "-100000000000000000000000000000000000000"},
}

func TestString(t *testing.T) {
	for i, v := range stringTests {
		if o := v.a.String(); o != v.s {
			t.Errorf("[%d] %q != %q", i, o, v.s)
		}
	}
}

func BenchmarkString(b *testing.B) {
	d := Decimal{0x5897e7bd6715a370, 0x17c4aea0fd62d52b}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.String()
	}
}

func BenchmarkStringSmall(b *testing.B) {
	d := FromI64(302187286)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.String()
	}
}

func TestRoundTrip(t *testing.T) {
	d := Decimal{0x5897e7bd6715a370, 0x17c4aea0fd62d52b}
	buf := make([]byte, 16)
	d.Write(buf)

	d2 := ReadDecimal(buf)
	if d != d2 {
		t.Errorf("%#v != %#v", d, d2)
	}
}
