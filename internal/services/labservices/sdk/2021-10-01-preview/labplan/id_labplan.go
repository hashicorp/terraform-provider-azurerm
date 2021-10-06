package labplan

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type LabPlanId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewLabPlanID(subscriptionId, resourceGroup, name string) LabPlanId {
	return LabPlanId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id LabPlanId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Lab Plan", segmentsStr)
}

func (id LabPlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.LabServices/labPlans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ParseLabPlanID parses a LabPlan ID into an LabPlanId struct
func ParseLabPlanID(input string) (*LabPlanId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LabPlanId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("labPlans"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseLabPlanIDInsensitively parses an LabPlan ID into an LabPlanId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseLabPlanID method should be used instead for validation etc.
func ParseLabPlanIDInsensitively(input string) (*LabPlanId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LabPlanId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'labPlans' segment
	labPlansKey := "labPlans"
	for key := range id.Path {
		if strings.EqualFold(key, labPlansKey) {
			labPlansKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(labPlansKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
