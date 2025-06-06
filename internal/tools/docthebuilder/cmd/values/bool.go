package values

import "strconv"

type BoolValue bool

func NewBoolValue(val bool, ptr *bool) *BoolValue {
	*ptr = val
	return (*BoolValue)(ptr)
}

func (v *BoolValue) Set(val bool) error {
	*v = BoolValue(val)
	return nil
}

func (v *BoolValue) String() string {
	return strconv.FormatBool(bool(*v))
}

func (v *BoolValue) Type() string {
	return "bool"
}
