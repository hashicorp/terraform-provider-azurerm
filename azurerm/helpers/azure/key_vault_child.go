package azure

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type KeyVaultChildID struct {
	KeyVaultBaseUrl string
	Name            string
	Version         string
}

func NewKeyVaultChildResourceID(keyVaultBaseUrl, childType, name, version string) (string, error) {
	fmtString := "%s/%s/%s/%s"
	keyVaultUrl, err := url.Parse(keyVaultBaseUrl)
	if err != nil || keyVaultBaseUrl == "" {
		return "", fmt.Errorf("failed to parse Key Vault Base URL %q: %+v", keyVaultBaseUrl, err)
	}
	// (@jackofallops) - Log Analytics service adds the port number to the API returns, so we strip it here
	if hostParts := strings.Split(keyVaultUrl.Host, ":"); len(hostParts) > 1 {
		keyVaultUrl.Host = hostParts[0]
	}

	return fmt.Sprintf(fmtString, keyVaultUrl.String(), childType, name, version), nil
}

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
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}

	return warnings, errors
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
