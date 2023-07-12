// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type InstanceFailoverGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	LocationName   string
	Name           string
}

func NewInstanceFailoverGroupID(subscriptionId, resourceGroup, locationName, name string) InstanceFailoverGroupId {
	return InstanceFailoverGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		LocationName:   locationName,
		Name:           name,
	}
}

func (id InstanceFailoverGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Location Name %q", id.LocationName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Instance Failover Group", segmentsStr)
}

func (id InstanceFailoverGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/locations/%s/instanceFailoverGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LocationName, id.Name)
}

// InstanceFailoverGroupID parses a InstanceFailoverGroup ID into an InstanceFailoverGroupId struct
func InstanceFailoverGroupID(input string) (*InstanceFailoverGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an InstanceFailoverGroup ID: %+v", input, err)
	}

	resourceId := InstanceFailoverGroupId{
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
	if resourceId.Name, err = id.PopSegment("instanceFailoverGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
