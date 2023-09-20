package parse

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
)

type ManagementGroupRemediationId struct {
	ManagementGroupId commonids.ManagementGroupId
	RemediationName   string
}

func NewManagementGroupRemediationID(managementGroupId commonids.ManagementGroupId, remediationName string) ManagementGroupRemediationId {
	return ManagementGroupRemediationId{
		ManagementGroupId: managementGroupId,
		RemediationName:   remediationName,
	}
}

func ParseManagementGroupRemediationID(input string) (*ManagementGroupRemediationId, error) {
	parsed, err := remediations.ParseScopedRemediationID(input)
	if err != nil {
		return nil, err
	}

	managementGroupId, err := commonids.ParseManagementGroupID(parsed.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Management Group ID: %+v", parsed.ResourceId, err)
	}

	return &ManagementGroupRemediationId{
		ManagementGroupId: *managementGroupId,
		RemediationName:   parsed.RemediationName,
	}, nil
}

func ParseManagementGroupRemediationIDInsensitively(input string) (*ManagementGroupRemediationId, error) {
	parsed, err := remediations.ParseScopedRemediationIDInsensitively(input)
	if err != nil {
		return nil, err
	}

	managementGroupId, err := commonids.ParseManagementGroupIDInsensitively(parsed.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Management Group ID: %+v", parsed.ResourceId, err)
	}

	return &ManagementGroupRemediationId{
		ManagementGroupId: *managementGroupId,
		RemediationName:   parsed.RemediationName,
	}, nil
}

func (id ManagementGroupRemediationId) ToRemediationID() remediations.ScopedRemediationId {
	return remediations.ScopedRemediationId{
		ResourceId:      id.ManagementGroupId.ID(),
		RemediationName: id.RemediationName,
	}
}
