package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type CapacityReservationId struct {
	SubscriptionId               string
	ResourceGroup                string
	CapacityReservationGroupName string
	Name                         string
}

func NewCapacityReservationID(subscriptionId, resourceGroup, capacityReservationGroupName, name string) CapacityReservationId {
	return CapacityReservationId{
		SubscriptionId:               subscriptionId,
		ResourceGroup:                resourceGroup,
		CapacityReservationGroupName: capacityReservationGroupName,
		Name:                         name,
	}
}

func (id CapacityReservationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Capacity Reservation Group Name %q", id.CapacityReservationGroupName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Capacity Reservation", segmentsStr)
}

func (id CapacityReservationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/capacityReservationGroups/%s/capacityReservations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CapacityReservationGroupName, id.Name)
}

// CapacityReservationID parses a CapacityReservation ID into an CapacityReservationId struct
func CapacityReservationID(input string) (*CapacityReservationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CapacityReservationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CapacityReservationGroupName, err = id.PopSegment("capacityReservationGroups"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("capacityReservations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
