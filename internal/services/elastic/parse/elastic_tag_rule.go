package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ElasticTagRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	MonitorName    string
	TagRuleName    string
}

func NewElasticTagRuleID(subscriptionId, resourceGroup, monitorName, tagRuleName string) ElasticTagRuleId {
	return ElasticTagRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		MonitorName:    monitorName,
		TagRuleName:    tagRuleName,
	}
}

func (id ElasticTagRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Tag Rule Name %q", id.TagRuleName),
		fmt.Sprintf("Monitor Name %q", id.MonitorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Elastic Tag Rule", segmentsStr)
}

func (id ElasticTagRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Elastic/monitors/%s/tagRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MonitorName, id.TagRuleName)
}

// ElasticTagRuleID parses a ElasticTagRule ID into an ElasticTagRuleId struct
func ElasticTagRuleID(input string) (*ElasticTagRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ElasticTagRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MonitorName, err = id.PopSegment("monitors"); err != nil {
		return nil, err
	}
	if resourceId.TagRuleName, err = id.PopSegment("tagRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
