package azure

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/validation"
)

type KeyVaultChildID struct {
	KeyVaultBaseUrl string
	Name            string
	Version         string
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

func ValidateKeyVaultChildName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]+$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}

	return ws, es
}

// Unfortunately this can't (easily) go in the Validate package
// since there's a circular reference on this package
func ValidateKeyVaultChildId(i interface{}, k string) (s []string, es []error) {
	if s, es = validation.NoZeroValues(i, k); len(es) > 0 {
		return s, es
	}

	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("Expected %s to be a string!", k))
		return s, es
	}

	_, err := ParseKeyVaultChildID(v)
	if err != nil {
		es = append(es, fmt.Errorf("Error parsing Key Vault Child ID: %s", err))
		return s, es
	}

	return s, es
}
