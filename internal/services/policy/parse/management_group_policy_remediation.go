package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagementGroupPolicyRemediationId struct {
	ManagementGroupName string
	RemediationName     string
}

func NewManagementGroupPolicyRemediationID(managementGroupName, remediationName string) ManagementGroupPolicyRemediationId {
	return ManagementGroupPolicyRemediationId{
		ManagementGroupName: managementGroupName,
		RemediationName:     remediationName,
	}
}

func (id ManagementGroupPolicyRemediationId) String() string {
	segments := []string{
		fmt.Sprintf("Remediation Name %q", id.RemediationName),
		fmt.Sprintf("Management Group Name %q", id.ManagementGroupName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Management Group Policy Remediation", segmentsStr)
}

func (id ManagementGroupPolicyRemediationId) ID() string {
	fmtString := "/providers/namespace1/managementGroups/%s/providers/Microsoft.PolicyInsights/remediations/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupName, id.RemediationName)
}

// ManagementGroupPolicyRemediationID parses a ManagementGroupPolicyRemediation ID into an ManagementGroupPolicyRemediationId struct
func ManagementGroupPolicyRemediationID(input string) (*ManagementGroupPolicyRemediationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagementGroupPolicyRemediationId{}

	if resourceId.ManagementGroupName, err = id.PopSegment("managementGroups"); err != nil {
		return nil, err
	}
	if resourceId.RemediationName, err = id.PopSegment("remediations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
