package storage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AccountID struct {
	Name          string
	ResourceGroup string

	ID azure.ResourceID
}

func AccountIDSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: ValidateAccountID,
	}
}

func ParseAccountID(id string) (*AccountID, error) {
	storageID, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return nil, err
	}

	resourceGroup := storageID.ResourceGroup
	if resourceGroup == "" {
		return nil, fmt.Errorf("%q is missing a Resource Group", id)
	}

	storageAccountName := storageID.Path["storageAccounts"]
	if storageAccountName == "" {
		return nil, fmt.Errorf("%q is missing the `storageAccounts` segment", id)
	}

	accountId := AccountID{
		Name:          storageAccountName,
		ResourceGroup: resourceGroup,
		ID:            *storageID,
	}
	return &accountId, nil
}

func ValidateAccountID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	id, err := azure.ParseAzureResourceID(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a Resource Id: %v", v, err))
	}

	if id != nil {
		if id.Path["storageAccounts"] == "" {
			errors = append(errors, fmt.Errorf("The 'storageAccounts' segment is missing from Resource ID %q", v))
		}
	}

	return warnings, errors
}
