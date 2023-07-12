// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type CassandraDatacenterId struct {
	SubscriptionId       string
	ResourceGroup        string
	CassandraClusterName string
	DataCenterName       string
}

func NewCassandraDatacenterID(subscriptionId, resourceGroup, cassandraClusterName, dataCenterName string) CassandraDatacenterId {
	return CassandraDatacenterId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		CassandraClusterName: cassandraClusterName,
		DataCenterName:       dataCenterName,
	}
}

func (id CassandraDatacenterId) String() string {
	segments := []string{
		fmt.Sprintf("Data Center Name %q", id.DataCenterName),
		fmt.Sprintf("Cassandra Cluster Name %q", id.CassandraClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cassandra Datacenter", segmentsStr)
}

func (id CassandraDatacenterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/cassandraClusters/%s/dataCenters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CassandraClusterName, id.DataCenterName)
}

// CassandraDatacenterID parses a CassandraDatacenter ID into an CassandraDatacenterId struct
func CassandraDatacenterID(input string) (*CassandraDatacenterId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an CassandraDatacenter ID: %+v", input, err)
	}

	resourceId := CassandraDatacenterId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CassandraClusterName, err = id.PopSegment("cassandraClusters"); err != nil {
		return nil, err
	}
	if resourceId.DataCenterName, err = id.PopSegment("dataCenters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
