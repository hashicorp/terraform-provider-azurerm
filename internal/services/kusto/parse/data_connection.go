// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DataConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	DatabaseName   string
	Name           string
}

func NewDataConnectionID(subscriptionId, resourceGroup, clusterName, databaseName, name string) DataConnectionId {
	return DataConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
		DatabaseName:   databaseName,
		Name:           name,
	}
}

func (id DataConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Database Name %q", id.DatabaseName),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Data Connection", segmentsStr)
}

func (id DataConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/databases/%s/dataConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
}

// DataConnectionID parses a DataConnection ID into an DataConnectionId struct
func DataConnectionID(input string) (*DataConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an DataConnection ID: %+v", input, err)
	}

	resourceId := DataConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("clusters"); err != nil {
		return nil, err
	}
	if resourceId.DatabaseName, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("dataConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// DataConnectionIDInsensitively parses an DataConnection ID into an DataConnectionId struct, insensitively
// This should only be used to parse an ID for rewriting, the DataConnectionID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func DataConnectionIDInsensitively(input string) (*DataConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DataConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'clusters' segment
	clustersKey := "clusters"
	for key := range id.Path {
		if strings.EqualFold(key, clustersKey) {
			clustersKey = key
			break
		}
	}
	if resourceId.ClusterName, err = id.PopSegment(clustersKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'databases' segment
	databasesKey := "databases"
	for key := range id.Path {
		if strings.EqualFold(key, databasesKey) {
			databasesKey = key
			break
		}
	}
	if resourceId.DatabaseName, err = id.PopSegment(databasesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'dataConnections' segment
	dataConnectionsKey := "dataConnections"
	for key := range id.Path {
		if strings.EqualFold(key, dataConnectionsKey) {
			dataConnectionsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(dataConnectionsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
