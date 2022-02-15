package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ScheduledQueryRulesId struct {
	SubscriptionId         string
	ResourceGroup          string
	ScheduledQueryRuleName string
}

func NewScheduledQueryRulesID(subscriptionId, resourceGroup, scheduledQueryRuleName string) ScheduledQueryRulesId {
	return ScheduledQueryRulesId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ScheduledQueryRuleName: scheduledQueryRuleName,
	}
}

func (id ScheduledQueryRulesId) String() string {
	segments := []string{
		fmt.Sprintf("Scheduled Query Rule Name %q", id.ScheduledQueryRuleName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Scheduled Query Rules", segmentsStr)
}

func (id ScheduledQueryRulesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/scheduledQueryRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ScheduledQueryRuleName)
}

// ScheduledQueryRulesID parses a ScheduledQueryRules ID into an ScheduledQueryRulesId struct
func ScheduledQueryRulesID(input string) (*ScheduledQueryRulesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ScheduledQueryRulesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ScheduledQueryRuleName, err = id.PopSegment("scheduledQueryRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ScheduledQueryRulesIDInsensitively parses an ScheduledQueryRules ID into an ScheduledQueryRulesId struct, insensitively
// This should only be used to parse an ID for rewriting, the ScheduledQueryRulesID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ScheduledQueryRulesIDInsensitively(input string) (*ScheduledQueryRulesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ScheduledQueryRulesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'scheduledQueryRules' segment
	scheduledQueryRulesKey := "scheduledQueryRules"
	for key := range id.Path {
		if strings.EqualFold(key, scheduledQueryRulesKey) {
			scheduledQueryRulesKey = key
			break
		}
	}
	if resourceId.ScheduledQueryRuleName, err = id.PopSegment(scheduledQueryRulesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
