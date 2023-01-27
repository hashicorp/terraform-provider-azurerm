package zones

func Expand(input []string) Schema {
	out := Schema{}

	if input != nil {
		for _, v := range input {
			out = append(out, v)
		}
	}

	return out
}

func ExpandUntyped(input []interface{}) []string {
	out := make([]string, 0)

	if input != nil {
		for _, v := range input {
			out = append(out, v.(string))
		}
	}

	return out
}
