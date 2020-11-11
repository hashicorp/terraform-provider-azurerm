package storage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
)

func AccountIDSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: ValidateAccountID,
	}
}

func ValidateAccountID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.ParseAccountID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a Resource ID: %v", v, err))
	}

	return warnings, errors
}
