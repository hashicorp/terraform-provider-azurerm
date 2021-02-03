package validate

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ServiceBusMaxSizeInMegabytes() schema.SchemaValidateFunc {
	return validation.IntInSlice([]int{
		1024,
		2048,
		3072,
		4096,
		5120,
		10240,
		20480,
		40960,
		81920,
	})
}
