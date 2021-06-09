package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SubscriptionRuleId struct {
	SubscriptionId   string
	ResourceGroup    string
	NamespaceName    string
	TopicName        string
	SubscriptionName string
	RuleName         string
}

func NewSubscriptionRuleID(subscriptionId, resourceGroup, namespaceName, topicName, subscriptionName, ruleName string) SubscriptionRuleId {
	return SubscriptionRuleId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		NamespaceName:    namespaceName,
		TopicName:        topicName,
		SubscriptionName: subscriptionName,
		RuleName:         ruleName,
	}
}

func (id SubscriptionRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Rule Name %q", id.RuleName),
		fmt.Sprintf("Subscription Name %q", id.SubscriptionName),
		fmt.Sprintf("Topic Name %q", id.TopicName),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subscription Rule", segmentsStr)
}

func (id SubscriptionRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/topics/%s/subscriptions/%s/rules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.TopicName, id.SubscriptionName, id.RuleName)
}

// SubscriptionRuleID parses a SubscriptionRule ID into an SubscriptionRuleId struct
func SubscriptionRuleID(input string) (*SubscriptionRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SubscriptionRuleId{
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
	if resourceId.TopicName, err = id.PopSegment("topics"); err != nil {
		return nil, err
	}
	if resourceId.SubscriptionName, err = id.PopSegment("subscriptions"); err != nil {
		return nil, err
	}
	if resourceId.RuleName, err = id.PopSegment("rules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
