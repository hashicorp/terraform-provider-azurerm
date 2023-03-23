package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type LogzSubAccountTagRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	MonitorName    string
	AccountName    string
	TagRuleName    string
}

func NewLogzSubAccountTagRuleID(subscriptionId, resourceGroup, monitorName, accountName, tagRuleName string) LogzSubAccountTagRuleId {
	return LogzSubAccountTagRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		MonitorName:    monitorName,
		AccountName:    accountName,
		TagRuleName:    tagRuleName,
	}
}

func (id LogzSubAccountTagRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Tag Rule Name %q", id.TagRuleName),
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Monitor Name %q", id.MonitorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Logz Sub Account Tag Rule", segmentsStr)
}

func (id LogzSubAccountTagRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logz/monitors/%s/accounts/%s/tagRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MonitorName, id.AccountName, id.TagRuleName)
}

// LogzSubAccountTagRuleID parses a LogzSubAccountTagRule ID into an LogzSubAccountTagRuleId struct
func LogzSubAccountTagRuleID(input string) (*LogzSubAccountTagRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogzSubAccountTagRuleId{
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
	if resourceId.AccountName, err = id.PopSegment("accounts"); err != nil {
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
