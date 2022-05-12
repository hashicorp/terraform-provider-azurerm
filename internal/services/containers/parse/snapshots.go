package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SnapshotsId struct {
	SubscriptionId string
	ResourceGroup  string
	SnapshotName   string
}

func NewSnapshotsID(subscriptionId, resourceGroup, snapshotName string) SnapshotsId {
	return SnapshotsId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SnapshotName:   snapshotName,
	}
}

func (id SnapshotsId) String() string {
	segments := []string{
		fmt.Sprintf("Snapshot Name %q", id.SnapshotName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Snapshots", segmentsStr)
}

func (id SnapshotsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/snapshots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SnapshotName)
}

// SnapshotsID parses a Snapshots ID into an SnapshotsId struct
func SnapshotsID(input string) (*SnapshotsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SnapshotsId{
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
