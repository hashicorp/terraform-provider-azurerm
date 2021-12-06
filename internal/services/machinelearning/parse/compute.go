package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ComputeId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkspaceName  string
	Name           string
}

func NewComputeID(subscriptionId, resourceGroup, workspaceName, name string) ComputeId {
	return ComputeId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkspaceName:  workspaceName,
		Name:           name,
	}
}

func (id ComputeId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Compute", segmentsStr)
}

func (id ComputeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/workspaces/%s/computes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.Name)
}

// ComputeID parses a Compute ID into an ComputeId struct
func ComputeID(input string) (*ComputeId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ComputeId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("computes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
