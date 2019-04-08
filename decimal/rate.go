package decimal

const rateBase = 1000000

// Rate represents a scalar rate as a fixed-point decimal.
type Rate struct {
	// TODO: dynamic base?
	r int64
}

// NewRate returns a Rate with the given number of parts-per-million.
func NewRate(ppm int64) Rate {
	return Rate{ppm}
}

// Add returns a new Rate that is the sum of the two arguments.
func (r Rate) Add(o Rate) Rate {
	return Rate{r.r + o.r}
}

// Sub returns a new Rate that is the result of subtracting the second rate from
// the first.
func (r Rate) Sub(o Rate) Rate {
	return Rate{r.r - o.r}
}

// Neg returns a new Rate that negates the given Rate.
func (r Rate) Neg() Rate {
	return Rate{-r.r}
}

func (r Rate) signAbs() (Rate, bool) {
	if r.r < 0 {
		return r.Neg(), true
	}
	return r, false
}

// Mul returns a new Rate that is the result of multiplying the two given Rates.
func (r Rate) Mul(o Rate) Rate {
	// TODO: the intermediate result here can overflow
	return Rate{(r.r * o.r) / rateBase}
}

// Div returns a new Rate that is the result of dividing the first Rate by the
// second Rate.
func (r Rate) Div(o Rate) Rate {
	// TODO: the intermediate result here can overflow
	return Rate{(r.r * rateBase) / o.r}
}
