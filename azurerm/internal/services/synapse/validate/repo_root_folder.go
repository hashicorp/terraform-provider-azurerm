package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func RepoRootFolder() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^\/(.*\/?)*$`),
		"Root folder must start with '/' and needs to be a valid git path")
}
