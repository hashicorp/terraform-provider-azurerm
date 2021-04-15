package azure

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
)

func SchemaResourceGroupName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validateResourceGroupName,
	}
}

func SchemaResourceGroupNameDeprecated() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Optional:     true,
		ValidateFunc: validateResourceGroupName,
		Deprecated:   "This field is no longer used and will be removed in the next major version of the Azure Provider",
	}
}

func SchemaResourceGroupNameDeprecatedComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validateResourceGroupName,
		Deprecated:   "This field is no longer used and will be removed in the next major version of the Azure Provider",
	}
}

func SchemaResourceGroupNameDiffSuppress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppress.CaseDifference,
		ValidateFunc:     validateResourceGroupName,
	}
}

func SchemaResourceGroupNameForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ValidateFunc: validateResourceGroupName,
	}
}

func SchemaResourceGroupNameOptionalComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		ForceNew:     true,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validateResourceGroupName,
	}
}

func SchemaResourceGroupNameOptional() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Optional:     true,
		ValidateFunc: validateResourceGroupName,
	}
}

func SchemaResourceGroupNameSetOptional() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validateResourceGroupName,
		},
	}
}

func validateResourceGroupName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) > 90 {
		errors = append(errors, fmt.Errorf("%q may not exceed 90 characters in length", k))
	}

	if strings.HasSuffix(value, ".") {
		errors = append(errors, fmt.Errorf("%q may not end with a period", k))
	}

	// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
	if matched := regexp.MustCompile(`^[-\w._()]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, dash, underscores, parentheses and periods", k))
	}

	return warnings, errors
}
