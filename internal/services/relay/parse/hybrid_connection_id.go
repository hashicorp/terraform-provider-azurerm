package parse

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/sdk/2017-04-01/hybridconnections"
)

func HybridConnectionID(input string) (*hybridconnections.HybridConnectionId, error) {
	return hybridconnections.ParseHybridConnectionID(input)
}
