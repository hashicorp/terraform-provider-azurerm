package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type UserAssignedIdentityId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewUserAssignedIdentityID(subscriptionId, resourceGroup, name string) UserAssignedIdentityId {
	return UserAssignedIdentityId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id UserAssignedIdentityId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "User Assigned Identity", segmentsStr)
}

func (id UserAssignedIdentityId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// UserAssignedIdentityID parses a UserAssignedIdentity ID into an UserAssignedIdentityId struct
func UserAssignedIdentityID(input string) (*UserAssignedIdentityId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := UserAssignedIdentityId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("userAssignedIdentities"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// UserAssignedIdentityIDInsensitively parses an UserAssignedIdentity ID into an UserAssignedIdentityId struct, insensitively
// This should only be used to parse an ID for rewriting, the UserAssignedIdentityID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func UserAssignedIdentityIDInsensitively(input string) (*UserAssignedIdentityId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := UserAssignedIdentityId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'userAssignedIdentities' segment
	userAssignedIdentitiesKey := "userAssignedIdentities"
	for key := range id.Path {
		if strings.EqualFold(key, userAssignedIdentitiesKey) {
			userAssignedIdentitiesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(userAssignedIdentitiesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
