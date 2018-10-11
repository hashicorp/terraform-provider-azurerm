package validate

import (
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func DevTestLabName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z0-9_-]+$"),
		"Lab Name can only include alphanumeric characters, underscores, hyphens.")
}

func DevTestVirtualNetworkUsagePermissionType() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(dtl.Allow),
		string(dtl.Default),
		string(dtl.Deny),
	}, false)
}
