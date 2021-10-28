package managedhsms

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DeletedManagedHSMId struct {
	SubscriptionId string
	LocationName   string
	Name           string
}

func NewDeletedManagedHSMID(subscriptionId, locationName, name string) DeletedManagedHSMId {
	return DeletedManagedHSMId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		Name:           name,
	}
}

func (id DeletedManagedHSMId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Location Name %q", id.LocationName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Deleted Managed H S M", segmentsStr)
}

func (id DeletedManagedHSMId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.KeyVault/locations/%s/deletedManagedHSMs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.Name)
}

// ParseDeletedManagedHSMID parses a DeletedManagedHSM ID into an DeletedManagedHSMId struct
func ParseDeletedManagedHSMID(input string) (*DeletedManagedHSMId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DeletedManagedHSMId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.LocationName, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("deletedManagedHSMs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseDeletedManagedHSMIDInsensitively parses an DeletedManagedHSM ID into an DeletedManagedHSMId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseDeletedManagedHSMID method should be used instead for validation etc.
func ParseDeletedManagedHSMIDInsensitively(input string) (*DeletedManagedHSMId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DeletedManagedHSMId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	// find the correct casing for the 'locations' segment
	locationsKey := "locations"
	for key := range id.Path {
		if strings.EqualFold(key, locationsKey) {
			locationsKey = key
			break
		}
	}
	if resourceId.LocationName, err = id.PopSegment(locationsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'deletedManagedHSMs' segment
	deletedManagedHSMsKey := "deletedManagedHSMs"
	for key := range id.Path {
		if strings.EqualFold(key, deletedManagedHSMsKey) {
			deletedManagedHSMsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(deletedManagedHSMsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
