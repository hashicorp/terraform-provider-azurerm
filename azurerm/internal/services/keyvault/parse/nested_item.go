package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
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
	return fmt.Sprintf("%s/%s/%s/%s", n.KeyVaultBaseUrl, n.NestedItemType, n.Name, n.Version)
}

func ParseNestedItemID(input string) (*NestedItemId, error) {
	return nil, nil
}
