package vaults

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DeletedVaultId struct {
	SubscriptionId string
	LocationName   string
	Name           string
}

func NewDeletedVaultID(subscriptionId, locationName, name string) DeletedVaultId {
	return DeletedVaultId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		Name:           name,
	}
}

func (id DeletedVaultId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Location Name %q", id.LocationName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Deleted Vault", segmentsStr)
}

func (id DeletedVaultId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.KeyVault/locations/%s/deletedVaults/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.Name)
}

// ParseDeletedVaultID parses a DeletedVault ID into an DeletedVaultId struct
func ParseDeletedVaultID(input string) (*DeletedVaultId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DeletedVaultId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.LocationName, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("deletedVaults"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseDeletedVaultIDInsensitively parses an DeletedVault ID into an DeletedVaultId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseDeletedVaultID method should be used instead for validation etc.
func ParseDeletedVaultIDInsensitively(input string) (*DeletedVaultId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DeletedVaultId{
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

	// find the correct casing for the 'deletedVaults' segment
	deletedVaultsKey := "deletedVaults"
	for key := range id.Path {
		if strings.EqualFold(key, deletedVaultsKey) {
			deletedVaultsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(deletedVaultsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
