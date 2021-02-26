package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func QueueName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z0-9][\w-./~]{0,258}([a-zA-Z0-9])?$`),
		"The topic name can contain only letters, numbers, periods, hyphens, tildas, forward slashes and underscores. The namespace must start and end with a letter or number and be less then 260 characters long.",
	)
}
