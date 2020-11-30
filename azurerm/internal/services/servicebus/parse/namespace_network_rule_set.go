package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NamespaceNetworkRuleSetId struct {
	ResourceGroup      string
	NamespaceName      string
	NetworkrulesetName string
}

func NamespaceNetworkRuleSetID(input string) (*NamespaceNetworkRuleSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace Network Rule Set ID %q: %+v", input, err)
	}

	rule := NamespaceNetworkRuleSetId{
		ResourceGroup: id.ResourceGroup,
	}

	if rule.NetworkrulesetName, err = id.PopSegment("networkrulesets"); err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace Network Rule Set ID %q: %+v", input, err)
	}

	if rule.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace Network Rule Set ID %q: %+v", input, err)
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace Network Rule Set ID %q: %+v", input, err)
	}

	return &rule, nil
}
