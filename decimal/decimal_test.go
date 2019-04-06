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
	{FromI64(2500000), NewRate(25), FromI64(63)},
	{FromI64(250000), NewRate(25), FromI64(6)},
	{Decimal{2, 0}, NewRate(333333), Decimal{0, 0xaaaa9f7b5aea3162}},
	{Decimal{2, 1}, NewRate(333333), Decimal{0, 0xaaaa9f7b5aea3162}},
	{Decimal{2, 2}, NewRate(333333), Decimal{0, 0xaaaa9f7b5aea3162}},
	{Decimal{5, 0}, NewRate(333333), Decimal{1, 0xaaaa8eb463497b74}},
	{FromI64(25000000), NewRate(-25), FromI64(-625)},
	{FromI64(-25000000), NewRate(25), FromI64(-625)},
	{FromI64(-25000000), NewRate(-25), FromI64(625)},
	{FromI64(-2500000), NewRate(25), FromI64(-63)},
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
}

func TestDiv(t *testing.T) {
	for i, v := range divTests {
		if o := v.a.Div(v.r); o != v.c {
			t.Errorf("[%d] %#v != %#v", i, o, v.c)
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
