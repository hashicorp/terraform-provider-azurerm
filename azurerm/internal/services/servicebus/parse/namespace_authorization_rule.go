package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NamespaceAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroup         string
	NamespaceName         string
	AuthorizationRuleName string
}

func NewNamespaceAuthorizationRuleID(subscriptionId, resourceGroup, namespaceName, authorizationRuleName string) NamespaceAuthorizationRuleId {
	return NamespaceAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		NamespaceName:         namespaceName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

func (id NamespaceAuthorizationRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Authorization Rule Name %q", id.AuthorizationRuleName),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Namespace Authorization Rule", segmentsStr)
}

func (id NamespaceAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/AuthorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.AuthorizationRuleName)
}

// NamespaceAuthorizationRuleID parses a NamespaceAuthorizationRule ID into an NamespaceAuthorizationRuleId struct
func NamespaceAuthorizationRuleID(input string) (*NamespaceAuthorizationRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NamespaceAuthorizationRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}
	if resourceId.AuthorizationRuleName, err = id.PopSegment("AuthorizationRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
