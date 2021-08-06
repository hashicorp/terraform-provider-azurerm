package parse

import "github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/sdk/namespaces"

func NamespaceID(input string) (*namespaces.NamespaceId, error) {
	return namespaces.ParseNamespaceID(input)
}
