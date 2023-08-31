// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ScriptId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	DatabaseName   string
	Name           string
}

func NewScriptID(subscriptionId, resourceGroup, clusterName, databaseName, name string) ScriptId {
	return ScriptId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
		DatabaseName:   databaseName,
		Name:           name,
	}
}

func (id ScriptId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Database Name %q", id.DatabaseName),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Script", segmentsStr)
}

func (id ScriptId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/databases/%s/scripts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
}

// ScriptID parses a Script ID into an ScriptId struct
func ScriptID(input string) (*ScriptId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Script ID: %+v", input, err)
	}

	resourceId := ScriptId{
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
	if resourceId.Name, err = id.PopSegment("scripts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ScriptIDInsensitively parses an Script ID into an ScriptId struct, insensitively
// This should only be used to parse an ID for rewriting, the ScriptID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ScriptIDInsensitively(input string) (*ScriptId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ScriptId{
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

	// find the correct casing for the 'scripts' segment
	scriptsKey := "scripts"
	for key := range id.Path {
		if strings.EqualFold(key, scriptsKey) {
			scriptsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(scriptsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
