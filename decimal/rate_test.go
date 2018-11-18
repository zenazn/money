package decimal

import "testing"

func TestRateAdd(t *testing.T) {
	r1 := NewRate(900000)
	r2 := NewRate(2000000)
	r3 := r1.Add(r2)
	if r3.r != 2900000 {
		t.Errorf("rate add %d", r3.r)
	}
}

func TestRateSub(t *testing.T) {
	r1 := NewRate(900000)
	r2 := NewRate(2000000)
	r3 := r1.Sub(r2)
	if r3.r != -1100000 {
		t.Errorf("rate sub %d", r3.r)
	}
}

func TestRateMul(t *testing.T) {
	r1 := NewRate(900000)
	r2 := NewRate(2000000)
	r3 := r1.Mul(r2)
	if r3.r != 1800000 {
		t.Errorf("rate mul %d", r3.r)
	}
}

func TestRateDiv(t *testing.T) {
	r1 := NewRate(900000)
	r2 := NewRate(2000000)
	r3 := r1.Div(r2)
	if r3.r != 450000 {
		t.Errorf("rate div %d", r3.r)
	}
}
