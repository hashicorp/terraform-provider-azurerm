package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TopicAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroup         string
	NamespaceName         string
	TopicName             string
	AuthorizationRuleName string
}

func NewTopicAuthorizationRuleID(subscriptionId, resourceGroup, namespaceName, topicName, authorizationRuleName string) TopicAuthorizationRuleId {
	return TopicAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		NamespaceName:         namespaceName,
		TopicName:             topicName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

func (id TopicAuthorizationRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Authorization Rule Name %q", id.AuthorizationRuleName),
		fmt.Sprintf("Topic Name %q", id.TopicName),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Topic Authorization Rule", segmentsStr)
}

func (id TopicAuthorizationRuleId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/topics/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.TopicName, id.AuthorizationRuleName)
}

// TopicAuthorizationRuleID parses a TopicAuthorizationRule ID into an TopicAuthorizationRuleId struct
func TopicAuthorizationRuleID(input string) (*TopicAuthorizationRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TopicAuthorizationRuleId{
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
	if resourceId.AuthorizationRuleName, err = id.PopSegment("authorizationRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
