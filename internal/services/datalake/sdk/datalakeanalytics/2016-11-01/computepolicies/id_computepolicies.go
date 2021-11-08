package computepolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ComputePoliciesId struct {
	SubscriptionId    string
	ResourceGroup     string
	AccountName       string
	ComputePolicyName string
}

func NewComputePoliciesID(subscriptionId, resourceGroup, accountName, computePolicyName string) ComputePoliciesId {
	return ComputePoliciesId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		AccountName:       accountName,
		ComputePolicyName: computePolicyName,
	}
}

func (id ComputePoliciesId) String() string {
	segments := []string{
		fmt.Sprintf("Compute Policy Name %q", id.ComputePolicyName),
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Compute Policies", segmentsStr)
}

func (id ComputePoliciesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeAnalytics/accounts/%s/computePolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName, id.ComputePolicyName)
}

// ParseComputePoliciesID parses a ComputePolicies ID into an ComputePoliciesId struct
func ParseComputePoliciesID(input string) (*ComputePoliciesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ComputePoliciesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AccountName, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}
	if resourceId.ComputePolicyName, err = id.PopSegment("computePolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseComputePoliciesIDInsensitively parses an ComputePolicies ID into an ComputePoliciesId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseComputePoliciesID method should be used instead for validation etc.
func ParseComputePoliciesIDInsensitively(input string) (*ComputePoliciesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ComputePoliciesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'accounts' segment
	accountsKey := "accounts"
	for key := range id.Path {
		if strings.EqualFold(key, accountsKey) {
			accountsKey = key
			break
		}
	}
	if resourceId.AccountName, err = id.PopSegment(accountsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'computePolicies' segment
	computePoliciesKey := "computePolicies"
	for key := range id.Path {
		if strings.EqualFold(key, computePoliciesKey) {
			computePoliciesKey = key
			break
		}
	}
	if resourceId.ComputePolicyName, err = id.PopSegment(computePoliciesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
