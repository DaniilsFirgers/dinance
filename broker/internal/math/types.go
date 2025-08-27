package math

type Numeric interface {
	int | int64 | float32 | float64
}

type NumericPtr interface {
	*int | *int64 | *float32 | *float64
}
