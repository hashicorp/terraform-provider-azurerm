package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServiceBusNamespaceNetworkRuleSetId struct {
	Name          string
	NamespaceName string
	ResourceGroup string
}

func ServiceBusNamespaceNetworkRuleSetID(input string) (*ServiceBusNamespaceNetworkRuleSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace Network Rule Set ID %q: %+v", input, err)
	}

	rule := ServiceBusNamespaceNetworkRuleSetId{
		ResourceGroup: id.ResourceGroup,
	}

	if rule.Name, err = id.PopSegment("networkrulesets"); err != nil {
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
