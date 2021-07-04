package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagementGroupAssignmentId struct {
	ManagementGroupName  string
	PolicyAssignmentName string
}

func NewManagementGroupAssignmentID(managementGroupName, policyAssignmentName string) ManagementGroupAssignmentId {
	return ManagementGroupAssignmentId{
		ManagementGroupName:  managementGroupName,
		PolicyAssignmentName: policyAssignmentName,
	}
}

func (id ManagementGroupAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Policy Assignment Name %q", id.PolicyAssignmentName),
		fmt.Sprintf("Management Group Name %q", id.ManagementGroupName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Management Group Assignment", segmentsStr)
}

func (id ManagementGroupAssignmentId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Authorization/policyAssignments/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupName, id.PolicyAssignmentName)
}

// ManagementGroupAssignmentID parses a ManagementGroupAssignment ID into an ManagementGroupAssignmentId struct
func ManagementGroupAssignmentID(input string) (*ManagementGroupAssignmentId, error) {
	// TODO: the generator should support outputting this method too
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagementGroupAssignmentId{}

	if resourceId.ManagementGroupName, err = id.PopSegment("managementGroups"); err != nil {
		return nil, err
	}
	if resourceId.PolicyAssignmentName, err = id.PopSegment("policyAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
