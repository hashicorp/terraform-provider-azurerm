package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var _ resourceid.Formatter = NestedItemId{}

type NestedItemId struct {
	KeyVaultBaseUrl string
	NestedItemType  string
	Name            string
	Version         string
}

func NewNestedItemID(keyVaultBaseUrl, nestedItemType, name, version string) (*NestedItemId, error) {
	keyVaultUrl, err := url.Parse(keyVaultBaseUrl)
	if err != nil || keyVaultBaseUrl == "" {
		return nil, fmt.Errorf("parsing %q: %+v", keyVaultBaseUrl, err)
	}
	// (@jackofallops) - Log Analytics service adds the port number to the API returns, so we strip it here
	if hostParts := strings.Split(keyVaultUrl.Host, ":"); len(hostParts) > 1 {
		keyVaultUrl.Host = hostParts[0]
	}

	return &NestedItemId{
		KeyVaultBaseUrl: keyVaultUrl.String(),
		NestedItemType:  nestedItemType,
		Name:            name,
		Version:         version,
	}, nil
}

func (n NestedItemId) ID() string {
	// example: https://tharvey-keyvault.vault.azure.net/type/bird/fdf067c93bbb4b22bff4d8b7a9a56217
	elements := []string{strings.TrimSuffix(n.KeyVaultBaseUrl, "/"), n.NestedItemType, n.Name, n.Version}
	return strings.Join(utils.RemoveFromStringArray(elements, ""), "/")
}

// ParseNestedItemID parses a Key Vault Nested Item ID (such as a Certificate, Key or Secret)
// containing a version into a NestedItemId object
func ParseNestedItemID(input string) (*NestedItemId, error) {
	item, err := parseNestedItemId(input)
	if err != nil {
		return nil, err
	}

	if item.Version == "" {
		return nil, fmt.Errorf("expected a versioned ID but no version in %q", input)
	}

	return item, nil
}

// ParseOptionallyVersionedNestedItemID parses a Key Vault Nested Item ID (such as a Certificate, Key or Secret)
// optionally containing a version into a NestedItemId object
func ParseOptionallyVersionedNestedItemID(input string) (*NestedItemId, error) {
	return parseNestedItemId(input)
}

func parseNestedItemId(id string) (*NestedItemId, error) {
	// versioned example: https://tharvey-keyvault.vault.azure.net/type/bird/fdf067c93bbb4b22bff4d8b7a9a56217
	// versionless example: https://tharvey-keyvault.vault.azure.net/type/bird/
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure KeyVault Child Id: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 2 && len(components) != 3 {
		return nil, fmt.Errorf("KeyVault Nested Item should contain 2 or 3 segments, got %d from %q", len(components), path)
	}

	version := ""
	if len(components) == 3 {
		version = components[2]
	}

	childId := NestedItemId{
		KeyVaultBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
		NestedItemType:  components[0],
		Name:            components[1],
		Version:         version,
	}

	return &childId, nil
}
