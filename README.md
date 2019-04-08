Money
=====

Money is an opinionated library for representing and manipulating monetary
values.

## Example

``` go
m1 := money.FromMinorUnits(123, currency.USD) // $1.23
m2, _ := money.Parse("4.567", "USD")
twoish := decimal.NewRate(2200000) // parts per million, so 2.2
fmt.Println(m1.Add(m2).Mul(twoish))
// Prints "USD 12.7534"
```
