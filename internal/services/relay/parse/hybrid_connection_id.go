package parse

import "github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/sdk/hybridconnections"

func HybridConnectionID(input string) (*hybridconnections.HybridConnectionId, error) {
	return hybridconnections.ParseHybridConnectionID(input)
}
