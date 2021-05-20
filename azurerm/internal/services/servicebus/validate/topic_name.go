package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func TopicName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9]([-._~a-zA-Z0-9]{0,258}[a-zA-Z0-9])?$"),
		"The topic name can contain only letters, numbers, periods, hyphens, tildas and underscores. The namespace must start with a letter or number, and it must end with a letter or number and be less then 260 characters long.",
	)
}
