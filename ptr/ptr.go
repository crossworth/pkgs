package ptr

// Of returns the value that v points to.
func Of[T any](v T) *T {
	return &v
}

// ValueOrDefault returns the value if the v is not nil
// or the default value provided on d.
func ValueOrDefault[T any](v *T, d T) T {
	if v == nil {
		return d
	}
	return *v
}

// ValuePtrOrNil returns the value pointer or nil if returnNil is true.
func ValuePtrOrNil[T any](v T, returnNil bool) *T {
	if returnNil {
		return nil
	}

	return &v
}

// ValueOf returns the value of the pointer.
func ValueOf[T any](v *T) T {
	return *v
}

// Float32To64 converts a *float32 to a *float64.
func Float32To64(f *float32) *float64 {
	if f == nil {
		return nil
	}
	return Of(float64(*f))
}
