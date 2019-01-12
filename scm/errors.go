package scm

// UnsupportedType is a string.
type UnsupportedType string

// Error returns the UnSupportedType.
func (a UnsupportedType) Error() string {
	return string(a)
}
