// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ArcKubernetesProvisionedClusterInstanceId struct {
	SubscriptionId                 string
	ResourceGroup                  string
	ConnectedClusterName           string
	ProvisionedClusterInstanceName string
}

func NewArcKubernetesProvisionedClusterInstanceID(subscriptionId, resourceGroup, connectedClusterName, provisionedClusterInstanceName string) ArcKubernetesProvisionedClusterInstanceId {
	return ArcKubernetesProvisionedClusterInstanceId{
		SubscriptionId:                 subscriptionId,
		ResourceGroup:                  resourceGroup,
		ConnectedClusterName:           connectedClusterName,
		ProvisionedClusterInstanceName: provisionedClusterInstanceName,
	}
}

func (id ArcKubernetesProvisionedClusterInstanceId) String() string {
	segments := []string{
		fmt.Sprintf("Provisioned Cluster Instance Name %q", id.ProvisionedClusterInstanceName),
		fmt.Sprintf("Connected Cluster Name %q", id.ConnectedClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Arc Kubernetes Provisioned Cluster Instance", segmentsStr)
}

func (id ArcKubernetesProvisionedClusterInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kubernetes/connectedClusters/%s/providers/Microsoft.HybridContainerService/provisionedClusterInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ConnectedClusterName, id.ProvisionedClusterInstanceName)
}

// ArcKubernetesProvisionedClusterInstanceID parses a ArcKubernetesProvisionedClusterInstance ID into an ArcKubernetesProvisionedClusterInstanceId struct
func ArcKubernetesProvisionedClusterInstanceID(input string) (*ArcKubernetesProvisionedClusterInstanceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ArcKubernetesProvisionedClusterInstance ID: %+v", input, err)
	}

	resourceId := ArcKubernetesProvisionedClusterInstanceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ConnectedClusterName, err = id.PopSegment("connectedClusters"); err != nil {
		return nil, err
	}
	if resourceId.ProvisionedClusterInstanceName, err = id.PopSegment("provisionedClusterInstances"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ArcKubernetesProvisionedClusterInstanceIDInsensitively parses an ArcKubernetesProvisionedClusterInstance ID into an ArcKubernetesProvisionedClusterInstanceId struct, insensitively
// This should only be used to parse an ID for rewriting, the ArcKubernetesProvisionedClusterInstanceID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ArcKubernetesProvisionedClusterInstanceIDInsensitively(input string) (*ArcKubernetesProvisionedClusterInstanceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ArcKubernetesProvisionedClusterInstanceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'connectedClusters' segment
	connectedClustersKey := "connectedClusters"
	for key := range id.Path {
		if strings.EqualFold(key, connectedClustersKey) {
			connectedClustersKey = key
			break
		}
	}
	if resourceId.ConnectedClusterName, err = id.PopSegment(connectedClustersKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'provisionedClusterInstances' segment
	provisionedClusterInstancesKey := "provisionedClusterInstances"
	for key := range id.Path {
		if strings.EqualFold(key, provisionedClusterInstancesKey) {
			provisionedClusterInstancesKey = key
			break
		}
	}
	if resourceId.ProvisionedClusterInstanceName, err = id.PopSegment(provisionedClusterInstancesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
