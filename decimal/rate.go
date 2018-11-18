package decimal

const rateBase = 1000000

type Rate struct {
	// TODO: base?
	r int64
}

func ParseRate(s string) Rate {
	panic("not implemented")
	return Rate{}
}

func NewRate(ppm int) Rate {
	return Rate{int64(ppm)}
}

func (r Rate) Add(o Rate) Rate {
	return Rate{r.r + o.r}
}

func (r Rate) Sub(o Rate) Rate {
	return Rate{r.r - o.r}
}

func (r Rate) Neg() Rate {
	return Rate{-r.r}
}

func (r Rate) signAbs() (Rate, bool) {
	if r.r < 0 {
		return r.Neg(), true
	}
	return r, false
}

func (r Rate) Mul(o Rate) Rate {
	return Rate{(r.r * o.r) / rateBase}
}

func (r Rate) Div(o Rate) Rate {
	return Rate{(r.r * rateBase) / o.r}
}
