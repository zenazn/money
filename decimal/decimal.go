// Package decimal implements a 128-bit, 38-decimal-digit signed integer. With
// an implicit scaling factor, it is appropriate for use as a fixed-point
// decimal, for instance to store monetary values.
package decimal

import (
	"encoding/binary"
	"math/bits"
)

// Decimal is a 128-bit, 38-decimal-digit signed integer.
type Decimal struct {
	hi, lo uint64
}

// FromI64 returns a Decimal representation of the given int64.
func FromI64(i int64) Decimal {
	if i < 0 {
		return Decimal{^uint64(0), uint64(i)}
	}
	return Decimal{0, uint64(i)}
}

// Add returns a Decimal that is the sum of its two operands.
func (d Decimal) Add(o Decimal) Decimal {
	lo, carry := bits.Add64(d.lo, o.lo, 0)
	hi, _ := bits.Add64(d.hi, o.hi, carry)
	return Decimal{hi, lo}
}

// Sub returns a Decimal that is the result of subtracting its second operand
// from its first operand.
func (d Decimal) Sub(o Decimal) Decimal {
	lo, carry := bits.Sub64(d.lo, o.lo, 0)
	hi, _ := bits.Sub64(d.hi, o.hi, carry)
	return Decimal{hi, lo}
}

// Neg returns the negative of the given Decimal.
func (d Decimal) Neg() Decimal {
	lo, carry := bits.Add64(^d.lo, 1, 0)
	hi, _ := bits.Add64(^d.hi, 0, carry)
	return Decimal{hi, lo}
}

func (d Decimal) signAbs() (Decimal, bool) {
	if d.hi>>63 == 1 {
		return d.Neg(), true
	}
	return d, false
}

// Mul returns a Decimal that is the result of multiplying the given Decimal by
// the given Rate. It performs calculations at full precision, and rounds the
// fractional component of the result using Bankers Rounding (round-half-even).
// If the result is larger than can be represented by a Decimal, Mul panics.
func (d Decimal) Mul(r Rate) Decimal {
	// We can multiply two's compliment negative numbers without too much
	// difficulty, but later on we'll need to divide as well, and that's not
	// something I know how to do. So instead, just normalize all the signs
	// up front and deal with things at the very end.
	r, rSign := r.signAbs()
	d, dSign := d.signAbs()

	//       hi  lo
	//           ru
	//     --------
	//      m1h m1l
	//  m2h m2l
	// ------------
	//   o3  o2  o1
	ru := uint64(r.r)
	m1h, m1l := bits.Mul64(d.lo, ru)
	m2h, m2l := bits.Mul64(d.hi, ru)
	o1 := m1l
	o2, c2 := bits.Add64(m1h, m2l, 0)
	o3, _ := bits.Add64(m2h, 0, c2)

	// Normalize by the rate's base
	t3, r3 := bits.Div64(0, o3, rateBase)
	t2, r2 := bits.Div64(r3, o2, rateBase)
	t1, r1 := bits.Div64(r2, o1, rateBase)

	// Round 0.5 towards the nearest even number
	if r1 > rateBase/2 || (r1 == rateBase/2 && t1&1 == 1) {
		t1, r1 = bits.Add64(t1, 1, 0)
		t2, r2 = bits.Add64(t2, 0, r1)
		t3, _ = bits.Add64(t3, 0, r2)
	}

	if t3 != 0 {
		panic("decimal: mul: overflow")
	}

	// Flip signs if necessary
	if rSign != dSign {
		var c1 uint64
		t1, c1 = bits.Add64(^t1, 1, 0)
		t2, _ = bits.Add64(^t2, 0, c1)
	}

	return Decimal{t2, t1}
}

// Div returns a Decimal that is the result of dividing the given Decimal by the
// given Rate. It performs calculations at full precision, and rounds the
// fractional component of the result using Bankers Rounding (round-half-even).
// Div will panic if the given rate is zero, or if the result of the computation
// is larger than can be represented by a Decimal.
func (d Decimal) Div(r Rate) Decimal {
	if r.r == 0 {
		panic("decimal: div: divide by zero")
	}

	// I have no idea how to divide negative two's compliment numbers, so
	// just normalize the signs up front and be done with it.
	r, rSign := r.signAbs()
	d, dSign := d.signAbs()

	// Multiply by the rate base now instead of later, since it'll produce
	// more accurate results
	m1h, m1l := bits.Mul64(d.lo, rateBase)
	m2h, m2l := bits.Mul64(d.hi, rateBase)
	o1 := m1l
	o2, c2 := bits.Add64(m1h, m2l, 0)
	o3, _ := bits.Add64(m2h, 0, c2)

	// Now divide by the rate
	ru := uint64(r.r)
	t3, r3 := bits.Div64(0, o3, ru)
	t2, r2 := bits.Div64(r3, o2, ru)
	t1, r1 := bits.Div64(r2, o1, ru)

	// Round 0.5 towards the nearest even number
	if r1 > ru/2 || (r1 == ru/2 && t1&1 == 1) {
		t1, r1 = bits.Add64(t1, 1, 0)
		t2, r2 = bits.Add64(t2, 0, r1)
		t3, _ = bits.Add64(t3, 0, r2)
	}

	if t3 != 0 {
		panic("decimal: div: overflow")
	}

	// Flip signs if necessary
	if rSign != dSign {
		var c1 uint64
		t1, c1 = bits.Add64(^t1, 1, 0)
		t2, _ = bits.Add64(^t2, 0, c1)
	}

	return Decimal{t2, t1}
}

func (d Decimal) divmod(n uint64) (Decimal, uint64) {
	if d.hi>>63 == 1 {
		panic("decimal: divmod arg is negative")
	}

	q2, r2 := bits.Div64(0, d.hi, n)
	q1, r1 := bits.Div64(r2, d.lo, n)
	return Decimal{q2, q1}, r1
}

// Lt returns true if the Decimal argument is less than the given Decimal.
func (d Decimal) Lt(o Decimal) bool {
	da, dn := d.signAbs()
	oa, on := o.signAbs()
	if dn && !on {
		return true
	} else if !dn && on {
		return false
	} else if dn {
		if da.hi > oa.hi {
			return true
		} else if da.hi < oa.hi {
			return false
		}
		return da.lo > oa.lo
	} else {
		if da.hi < oa.hi {
			return true
		} else if da.hi > oa.hi {
			return false
		}
		return da.lo < oa.lo
	}
}

// Lte returns true if the Decimal argument is less than or equal to the given Decimal.
func (d Decimal) Lte(o Decimal) bool {
	if d == o {
		return true
	}
	return d.Lt(o)
}

// Eq returns true if the two Decimal arguments are equal.
func (d Decimal) Eq(o Decimal) bool {
	return d == o
}

// Gt returns true if the Decimal argument is greater than the given Decimal.
func (d Decimal) Gt(o Decimal) bool {
	return !d.Lte(o)
}

// Gte returns true if the Decimal argument is greater than or equal to the
// given Decimal.
func (d Decimal) Gte(o Decimal) bool {
	return !d.Lt(o)
}

// String returns a decimal string representing the given value.
func (d Decimal) String() string {
	var buf [40]byte
	k := len(buf) - 1
	d, s := d.signAbs()
	var rem uint64

	// I'm sure we can do better than this with a little bit of effort, but
	// this runs in a couple hundred nanoseconds on my computer, which is
	// fast enough for me for now.
	for d.hi != 0 || d.lo != 0 {
		d, rem = d.divmod(10000000000)
		for i := 0; i < 10; i++ {
			buf[k] = '0' + byte(rem%10)
			k = k - 1
			rem = rem / 10
		}
	}

	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 || buf[i] == '0' {
			continue
		}

		if s {
			i = i - 1
			buf[i] = '-'
		}
		return string(buf[i:])
	}

	return "0"
}

// ReadDecimal reads a Decimal, encoded as a 128-bit little-endian value, from
// the given byte slice.
func ReadDecimal(buf []byte) Decimal {
	lo := binary.LittleEndian.Uint64(buf[0:8])
	hi := binary.LittleEndian.Uint64(buf[8:16])
	return Decimal{hi, lo}
}

// Write writes the Decimal as a 128-bit little-endian value to the first 16
// bytes of the given byte slice.
func (d Decimal) Write(buf []byte) {
	binary.LittleEndian.PutUint64(buf[0:8], d.lo)
	binary.LittleEndian.PutUint64(buf[8:16], d.hi)
}
