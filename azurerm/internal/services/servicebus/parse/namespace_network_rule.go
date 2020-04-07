package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServiceBusNamespaceNetworkRuleId struct {
	Name          string
	NamespaceName string
	ResourceGroup string
}

func ServiceBusNamespaceNetworkRuleID(input string) (*ServiceBusNamespaceNetworkRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace Network Rule ID %q: %+v", input, err)
	}

	rule := ServiceBusNamespaceNetworkRuleId{
		ResourceGroup: id.ResourceGroup,
	}

	if rule.Name, err = id.PopSegment("networkrulesets"); err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace Network Rule ID %q: %+v", input, err)
	}

	if rule.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace Network Rule ID %q: %+v", input, err)
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace Network Rule ID %q: %+v", input, err)
	}

	return &rule, nil
}
