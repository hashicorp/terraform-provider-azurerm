// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedInstanceFailoverGroupId struct {
	SubscriptionId            string
	ResourceGroup             string
	LocationName              string
	InstanceFailoverGroupName string
}

func NewManagedInstanceFailoverGroupID(subscriptionId, resourceGroup, locationName, instanceFailoverGroupName string) ManagedInstanceFailoverGroupId {
	return ManagedInstanceFailoverGroupId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		LocationName:              locationName,
		InstanceFailoverGroupName: instanceFailoverGroupName,
	}
}

func (id ManagedInstanceFailoverGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Instance Failover Group Name %q", id.InstanceFailoverGroupName),
		fmt.Sprintf("Location Name %q", id.LocationName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed Instance Failover Group", segmentsStr)
}

func (id ManagedInstanceFailoverGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/locations/%s/instanceFailoverGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LocationName, id.InstanceFailoverGroupName)
}

// ManagedInstanceFailoverGroupID parses a ManagedInstanceFailoverGroup ID into an ManagedInstanceFailoverGroupId struct
func ManagedInstanceFailoverGroupID(input string) (*ManagedInstanceFailoverGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ManagedInstanceFailoverGroup ID: %+v", input, err)
	}

	resourceId := ManagedInstanceFailoverGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LocationName, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}
	if resourceId.InstanceFailoverGroupName, err = id.PopSegment("instanceFailoverGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
