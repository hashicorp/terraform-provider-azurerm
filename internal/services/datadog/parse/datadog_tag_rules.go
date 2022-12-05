package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DatadogTagRulesId struct {
	SubscriptionId string
	ResourceGroup  string
	MonitorName    string
	TagRuleName    string
}

func NewDatadogTagRulesID(subscriptionId, resourceGroup, monitorName, tagRuleName string) DatadogTagRulesId {
	return DatadogTagRulesId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		MonitorName:    monitorName,
		TagRuleName:    tagRuleName,
	}
}

func (id DatadogTagRulesId) String() string {
	segments := []string{
		fmt.Sprintf("Tag Rule Name %q", id.TagRuleName),
		fmt.Sprintf("Monitor Name %q", id.MonitorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Datadog Tag Rules", segmentsStr)
}

func (id DatadogTagRulesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Datadog/monitors/%s/tagRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MonitorName, id.TagRuleName)
}

// DatadogTagRulesID parses a DatadogTagRules ID into an DatadogTagRulesId struct
func DatadogTagRulesID(input string) (*DatadogTagRulesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DatadogTagRulesId{
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
