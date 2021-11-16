package tags

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func Flatten(input *map[string]string) map[string]*string {
	output := make(map[string]*string)
	if input == nil {
		return output
	}

	for k, v := range *input {
		output[k] = pointer.ToString(v)
	}

	return output
}
