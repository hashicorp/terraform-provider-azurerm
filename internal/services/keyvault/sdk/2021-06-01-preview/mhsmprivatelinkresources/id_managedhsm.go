package mhsmprivatelinkresources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedHSMId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewManagedHSMID(subscriptionId, resourceGroup, name string) ManagedHSMId {
	return ManagedHSMId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ManagedHSMId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed H S M", segmentsStr)
}

func (id ManagedHSMId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/managedHSMs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ParseManagedHSMID parses a ManagedHSM ID into an ManagedHSMId struct
func ParseManagedHSMID(input string) (*ManagedHSMId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagedHSMId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("managedHSMs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseManagedHSMIDInsensitively parses an ManagedHSM ID into an ManagedHSMId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseManagedHSMID method should be used instead for validation etc.
func ParseManagedHSMIDInsensitively(input string) (*ManagedHSMId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagedHSMId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'managedHSMs' segment
	managedHSMsKey := "managedHSMs"
	for key := range id.Path {
		if strings.EqualFold(key, managedHSMsKey) {
			managedHSMsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(managedHSMsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
