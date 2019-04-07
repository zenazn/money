package money

func (m Money) LtErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Lt(o.amt), nil
}

func (m Money) Lt(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Lt(o.amt)
}

func (m Money) LteErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Lte(o.amt), nil
}

func (m Money) Lte(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Lte(o.amt)
}

func (m Money) EqErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Eq(o.amt), nil
}

func (m Money) Eq(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Eq(o.amt)
}

func (m Money) GteErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Gte(o.amt), nil
}

func (m Money) Gte(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Gte(o.amt)
}

func (m Money) GtErr(o Money) (bool, error) {
	if err := compat(m.ccy, o.ccy); err != nil {
		return false, err
	}
	return m.amt.Gt(o.amt), nil
}

func (m Money) Gt(o Money) bool {
	if err := compat(m.ccy, o.ccy); err != nil {
		panic(err.Error())
	}
	return m.amt.Gt(o.amt)
}
