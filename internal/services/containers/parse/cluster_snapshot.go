package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ClusterSnapshotId struct {
	SubscriptionId             string
	ResourceGroup              string
	ManagedclustersnapshotName string
}

func NewClusterSnapshotID(subscriptionId, resourceGroup, managedclustersnapshotName string) ClusterSnapshotId {
	return ClusterSnapshotId{
		SubscriptionId:             subscriptionId,
		ResourceGroup:              resourceGroup,
		ManagedclustersnapshotName: managedclustersnapshotName,
	}
}

func (id ClusterSnapshotId) String() string {
	segments := []string{
		fmt.Sprintf("Managedclustersnapshot Name %q", id.ManagedclustersnapshotName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cluster Snapshot", segmentsStr)
}

func (id ClusterSnapshotId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedclustersnapshots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedclustersnapshotName)
}

// ClusterSnapshotID parses a ClusterSnapshot ID into an ClusterSnapshotId struct
func ClusterSnapshotID(input string) (*ClusterSnapshotId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ClusterSnapshotId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedclustersnapshotName, err = id.PopSegment("managedclustersnapshots"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
