package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/sdk/hybridconnections"

func HybridConnectionID(input string) (*hybridconnections.HybridConnectionId, error) {
	return hybridconnections.ParseHybridConnectionID(input)
}
