package vaults

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type OperationKindId struct {
	SubscriptionId   string
	ResourceGroup    string
	VaultName        string
	AccessPolicyName string
}

func NewOperationKindID(subscriptionId, resourceGroup, vaultName, accessPolicyName string) OperationKindId {
	return OperationKindId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		VaultName:        vaultName,
		AccessPolicyName: accessPolicyName,
	}
}

func (id OperationKindId) String() string {
	segments := []string{
		fmt.Sprintf("Access Policy Name %q", id.AccessPolicyName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Operation Kind", segmentsStr)
}

func (id OperationKindId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/accessPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.AccessPolicyName)
}

// ParseOperationKindID parses a OperationKind ID into an OperationKindId struct
func ParseOperationKindID(input string) (*OperationKindId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := OperationKindId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VaultName, err = id.PopSegment("vaults"); err != nil {
		return nil, err
	}
	if resourceId.AccessPolicyName, err = id.PopSegment("accessPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseOperationKindIDInsensitively parses an OperationKind ID into an OperationKindId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseOperationKindID method should be used instead for validation etc.
func ParseOperationKindIDInsensitively(input string) (*OperationKindId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := OperationKindId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'vaults' segment
	vaultsKey := "vaults"
	for key := range id.Path {
		if strings.EqualFold(key, vaultsKey) {
			vaultsKey = key
			break
		}
	}
	if resourceId.VaultName, err = id.PopSegment(vaultsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'accessPolicies' segment
	accessPoliciesKey := "accessPolicies"
	for key := range id.Path {
		if strings.EqualFold(key, accessPoliciesKey) {
			accessPoliciesKey = key
			break
		}
	}
	if resourceId.AccessPolicyName, err = id.PopSegment(accessPoliciesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
