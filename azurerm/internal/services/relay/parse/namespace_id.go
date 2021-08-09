package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/sdk/namespaces"

func NamespaceID(input string) (*namespaces.NamespaceId, error) {
	return namespaces.ParseNamespaceID(input)
}
