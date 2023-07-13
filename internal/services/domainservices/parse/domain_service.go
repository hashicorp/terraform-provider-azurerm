// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DomainServiceId struct {
	SubscriptionId          string
	ResourceGroup           string
	Name                    string
	InitialReplicaSetIdName string
}

func NewDomainServiceID(subscriptionId, resourceGroup, name, initialReplicaSetIdName string) DomainServiceId {
	return DomainServiceId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		Name:                    name,
		InitialReplicaSetIdName: initialReplicaSetIdName,
	}
}

func (id DomainServiceId) String() string {
	segments := []string{
		fmt.Sprintf("Initial Replica Set Id Name %q", id.InitialReplicaSetIdName),
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Domain Service", segmentsStr)
}

func (id DomainServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AAD/domainServices/%s/initialReplicaSetId/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name, id.InitialReplicaSetIdName)
}

// DomainServiceID parses a DomainService ID into an DomainServiceId struct
func DomainServiceID(input string) (*DomainServiceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an DomainService ID: %+v", input, err)
	}

	resourceId := DomainServiceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("domainServices"); err != nil {
		return nil, err
	}
	if resourceId.InitialReplicaSetIdName, err = id.PopSegment("initialReplicaSetId"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
