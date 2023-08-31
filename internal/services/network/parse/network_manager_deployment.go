// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type ManagerDeploymentId struct {
	SubscriptionId     string
	ResourceGroup      string
	NetworkManagerName string
	Location           string
	ScopeAccess        string
}

func NewNetworkManagerDeploymentID(subscriptionId string, resourceGroup string, networkManagerName string, location string, scopeAccess string) *ManagerDeploymentId {
	return &ManagerDeploymentId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		NetworkManagerName: networkManagerName,
		Location:           azure.NormalizeLocation(location),
		ScopeAccess:        scopeAccess,
	}
}

func NetworkManagerDeploymentID(networkManagerDeploymentId string) (*ManagerDeploymentId, error) {
	v := strings.Split(networkManagerDeploymentId, "|")
	if len(v) != 3 {
		return nil, fmt.Errorf("expected the network manager deployment ID to be in format `{networkManagerId}/commit|{location}|{scopeAccess}`, but got %d segments ", len(v))
	}

	networkManagerId := strings.TrimSuffix(v[0], "/commit")
	managerId, err := networkmanagers.ParseNetworkManagerID(networkManagerId)
	if err != nil {
		return nil, err
	}

	if v[1] == "" {
		return nil, fmt.Errorf("expected location in network manager deployment ID with format `{networkManagerId}/commit|{location}|{scopeAccess}`, but got %s in %s", v[1], networkManagerDeploymentId)
	}
	normalizedLocation := azure.NormalizeLocation(v[1])

	if v[2] == "" {
		return nil, fmt.Errorf("expected scopeAccess in network manager deployment ID with format `{networkManagerId}/commit|{location}|{scopeAccess} to be one of the [Connectivity, SecurityAdmin]`, but got %s in %s", v[2], networkManagerDeploymentId)
	}
	scopeAccess := v[2]
	networkManagerDeployment := NewNetworkManagerDeploymentID(managerId.SubscriptionId, managerId.ResourceGroupName, managerId.NetworkManagerName, normalizedLocation, scopeAccess)
	return networkManagerDeployment, nil
}

func (id ManagerDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/commit|%s|%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.Location, id.ScopeAccess)
}

func (id ManagerDeploymentId) String() string {
	segments := []string{
		fmt.Sprintf("Scope Access %q", id.ScopeAccess),
		fmt.Sprintf("Location %q", id.Location),
		fmt.Sprintf("Network Manager %q", id.NetworkManagerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Manager Deployment", segmentsStr)
}
