// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AttachedDatabaseConfigurationId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	Name           string
}

func NewAttachedDatabaseConfigurationID(subscriptionId, resourceGroup, clusterName, name string) AttachedDatabaseConfigurationId {
	return AttachedDatabaseConfigurationId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
		Name:           name,
	}
}

func (id AttachedDatabaseConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Attached Database Configuration", segmentsStr)
}

func (id AttachedDatabaseConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/attachedDatabaseConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.Name)
}

// AttachedDatabaseConfigurationID parses a AttachedDatabaseConfiguration ID into an AttachedDatabaseConfigurationId struct
func AttachedDatabaseConfigurationID(input string) (*AttachedDatabaseConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AttachedDatabaseConfiguration ID: %+v", input, err)
	}

	resourceId := AttachedDatabaseConfigurationId{
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
	if resourceId.Name, err = id.PopSegment("attachedDatabaseConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// AttachedDatabaseConfigurationIDInsensitively parses an AttachedDatabaseConfiguration ID into an AttachedDatabaseConfigurationId struct, insensitively
// This should only be used to parse an ID for rewriting, the AttachedDatabaseConfigurationID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func AttachedDatabaseConfigurationIDInsensitively(input string) (*AttachedDatabaseConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AttachedDatabaseConfigurationId{
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

	// find the correct casing for the 'attachedDatabaseConfigurations' segment
	attachedDatabaseConfigurationsKey := "attachedDatabaseConfigurations"
	for key := range id.Path {
		if strings.EqualFold(key, attachedDatabaseConfigurationsKey) {
			attachedDatabaseConfigurationsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(attachedDatabaseConfigurationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
