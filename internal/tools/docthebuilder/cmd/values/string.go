package values

type StringValue string

func NewStringValue(val string, ptr *string) *StringValue {
	*ptr = val
	return (*StringValue)(ptr)
}

func (v *StringValue) Set(val string) error {
	*v = StringValue(val)
	return nil
}

func (v *StringValue) String() string {
	return string(*v)
}

func (v *StringValue) Type() string {
	return "string"
}
