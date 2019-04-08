package money

// LtErr returns true if the argument is less than the receiver, or an error if
// the two are not comparable (i.e., they have different currencies).
func (m Money) LtErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Lt(o.amt), nil
}

// Lt returns true if the argument is less than the receiver, or panics if the
// two are not comparable (i.e., they have different currencies).
func (m Money) Lt(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Lt(o.amt)
}

// LteErr returns true if the argument is less than or equal to the receiver, or
// an error if the two are not comparable (i.e., they have different
// currencies).
func (m Money) LteErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Lte(o.amt), nil
}

// Lte returns true if the argument is less than or equal to the receiver, or
// panics if the two are not comparable (i.e., they have different currencies).
func (m Money) Lte(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Lte(o.amt)
}

// EqErr returns true if the argument is equal to the receiver, or an error if
// the two are not comparable (i.e., they have different currencies).
func (m Money) EqErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Eq(o.amt), nil
}

// Eq returns true if the argument is equal to the receiver, or panics if the
// two are not comparable (i.e., they have different currencies).
func (m Money) Eq(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Eq(o.amt)
}

// GteErr returns true if the argument is greater than or equal to the receiver,
// or an error if the two are not comparable (i.e., they have different
// currencies).
func (m Money) GteErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Gte(o.amt), nil
}

// Gte returns true if the argument is greater than or equal to the receiver, or
// panics if the two are not comparable (i.e., they have different currencies).
func (m Money) Gte(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Gte(o.amt)
}

// GtErr returns true if the argument is greater than the receiver, or an error
// if the two are not comparable (i.e., they have different currencies).
func (m Money) GtErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Gt(o.amt), nil
}

// Gt returns true if the argument is greater than the receiver, or panics if
// the two are not comparable (i.e., they have different currencies).
func (m Money) Gt(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Gt(o.amt)
}
