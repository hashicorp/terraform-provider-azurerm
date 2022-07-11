package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DevTestLabPolicyId struct {
	SubscriptionId string
	ResourceGroup  string
	LabName        string
	PolicySetName  string
	PolicyName     string
}

func NewDevTestLabPolicyID(subscriptionId, resourceGroup, labName, policySetName, policyName string) DevTestLabPolicyId {
	return DevTestLabPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		LabName:        labName,
		PolicySetName:  policySetName,
		PolicyName:     policyName,
	}
}

func (id DevTestLabPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Policy Name %q", id.PolicyName),
		fmt.Sprintf("Policy Set Name %q", id.PolicySetName),
		fmt.Sprintf("Lab Name %q", id.LabName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Dev Test Lab Policy", segmentsStr)
}

func (id DevTestLabPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevTestLab/labs/%s/policySets/%s/policies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LabName, id.PolicySetName, id.PolicyName)
}

// DevTestLabPolicyID parses a DevTestLabPolicy ID into an DevTestLabPolicyId struct
func DevTestLabPolicyID(input string) (*DevTestLabPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DevTestLabPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LabName, err = id.PopSegment("labs"); err != nil {
		return nil, err
	}
	if resourceId.PolicySetName, err = id.PopSegment("policySets"); err != nil {
		return nil, err
	}
	if resourceId.PolicyName, err = id.PopSegment("policies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// DevTestLabPolicyIDInsensitively parses an DevTestLabPolicy ID into an DevTestLabPolicyId struct, insensitively
// This should only be used to parse an ID for rewriting, the DevTestLabPolicyID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func DevTestLabPolicyIDInsensitively(input string) (*DevTestLabPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DevTestLabPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'labs' segment
	labsKey := "labs"
	for key := range id.Path {
		if strings.EqualFold(key, labsKey) {
			labsKey = key
			break
		}
	}
	if resourceId.LabName, err = id.PopSegment(labsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'policySets' segment
	policySetsKey := "policySets"
	for key := range id.Path {
		if strings.EqualFold(key, policySetsKey) {
			policySetsKey = key
			break
		}
	}
	if resourceId.PolicySetName, err = id.PopSegment(policySetsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'policies' segment
	policiesKey := "policies"
	for key := range id.Path {
		if strings.EqualFold(key, policiesKey) {
			policiesKey = key
			break
		}
	}
	if resourceId.PolicyName, err = id.PopSegment(policiesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
