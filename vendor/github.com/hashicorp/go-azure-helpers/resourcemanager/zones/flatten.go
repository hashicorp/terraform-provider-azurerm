package zones

func Flatten(input *[]string) []interface{} {
	out := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			out = append(out, v)
		}
	}

	return out
}
