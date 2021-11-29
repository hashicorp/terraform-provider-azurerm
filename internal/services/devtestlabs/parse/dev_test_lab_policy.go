package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
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
	id, err := azure.ParseAzureResourceID(input)
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
