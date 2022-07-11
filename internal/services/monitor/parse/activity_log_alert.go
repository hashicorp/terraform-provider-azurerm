package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ActivityLogAlertId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewActivityLogAlertID(subscriptionId, resourceGroup, name string) ActivityLogAlertId {
	return ActivityLogAlertId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ActivityLogAlertId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Activity Log Alert", segmentsStr)
}

func (id ActivityLogAlertId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/activityLogAlerts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ActivityLogAlertID parses a ActivityLogAlert ID into an ActivityLogAlertId struct
func ActivityLogAlertID(input string) (*ActivityLogAlertId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ActivityLogAlertId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("activityLogAlerts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ActivityLogAlertIDInsensitively parses an ActivityLogAlert ID into an ActivityLogAlertId struct, insensitively
// This should only be used to parse an ID for rewriting, the ActivityLogAlertID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ActivityLogAlertIDInsensitively(input string) (*ActivityLogAlertId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ActivityLogAlertId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'activityLogAlerts' segment
	activityLogAlertsKey := "activityLogAlerts"
	for key := range id.Path {
		if strings.EqualFold(key, activityLogAlertsKey) {
			activityLogAlertsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(activityLogAlertsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
