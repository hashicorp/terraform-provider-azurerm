package azure

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
)

type KeyVaultChildID = keyVaultParse.NestedItemId

func ParseKeyVaultChildID(id string) (*KeyVaultChildID, error) {
	// example: https://tharvey-keyvault.vault.azure.net/type/bird/fdf067c93bbb4b22bff4d8b7a9a56217
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure KeyVault Child Id: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 3 {
		return nil, fmt.Errorf("Azure KeyVault Child Id should have 3 segments, got %d: '%s'", len(components), path)
	}

	childId := KeyVaultChildID{
		KeyVaultBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
		Name:            components[1],
		Version:         components[2],
	}

	return &childId, nil
}

func ParseKeyVaultChildIDVersionOptional(id string) (*KeyVaultChildID, error) {
	// example: https://tharvey-keyvault.vault.azure.net/type/bird/fdf067c93bbb4b22bff4d8b7a9a56217
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure KeyVault Child Id: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 2 && len(components) != 3 {
		return nil, fmt.Errorf("Azure KeyVault Child Id should have 2 or 3 segments, got %d: '%s'", len(components), path)
	}

	version := ""
	if len(components) == 3 {
		version = components[2]
	}

	childId := KeyVaultChildID{
		KeyVaultBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
		Name:            components[1],
		Version:         version,
	}

	return &childId, nil
}

func ValidateKeyVaultChildName(v interface{}, k string) (warnings []string, errors []error) {
	return keyVaultValidate.NestedItemName(v, k)
}

// Unfortunately this can't (easily) go in the Validate package
// since there's a circular reference on this package
func ValidateKeyVaultChildId(i interface{}, k string) (warnings []string, errors []error) {
	if warnings, errors = validation.StringIsNotEmpty(i, k); len(errors) > 0 {
		return warnings, errors
	}

	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("Expected %s to be a string!", k))
		return warnings, errors
	}

	if _, err := ParseKeyVaultChildID(v); err != nil {
		errors = append(errors, fmt.Errorf("Error parsing Key Vault Child ID: %s", err))
		return warnings, errors
	}

	return warnings, errors
}

// Unfortunately this can't (easily) go in the Validate package
// since there's a circular reference on this package
func ValidateKeyVaultChildIdVersionOptional(i interface{}, k string) (warnings []string, errors []error) {
	if warnings, errors = validation.StringIsNotEmpty(i, k); len(errors) > 0 {
		return warnings, errors
	}

	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("Expected %s to be a string!", k))
		return warnings, errors
	}

	if _, err := ParseKeyVaultChildIDVersionOptional(v); err != nil {
		errors = append(errors, fmt.Errorf("Error parsing Key Vault Child ID: %s", err))
		return warnings, errors
	}

	return warnings, errors
}
