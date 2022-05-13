package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NodePoolSnapshotId struct {
	SubscriptionId string
	ResourceGroup  string
	SnapshotName   string
}

func NewNodePoolSnapshotID(subscriptionId, resourceGroup, snapshotName string) NodePoolSnapshotId {
	return NodePoolSnapshotId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SnapshotName:   snapshotName,
	}
}

func (id NodePoolSnapshotId) String() string {
	segments := []string{
		fmt.Sprintf("Snapshot Name %q", id.SnapshotName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Node Pool Snapshot", segmentsStr)
}

func (id NodePoolSnapshotId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/snapshots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SnapshotName)
}

// NodePoolSnapshotID parses a NodePoolSnapshot ID into an NodePoolSnapshotId struct
func NodePoolSnapshotID(input string) (*NodePoolSnapshotId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NodePoolSnapshotId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SnapshotName, err = id.PopSegment("snapshots"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
