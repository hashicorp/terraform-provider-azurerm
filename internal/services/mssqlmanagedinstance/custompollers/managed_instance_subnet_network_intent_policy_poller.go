// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/subnets"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = managedInstanceSubnetNetworkIntentPolicyPoller{}

type SubnetsClient interface {
	Get(ctx context.Context, id commonids.SubnetId, options subnets.GetOperationOptions) (subnets.GetOperationResponse, error)
}

type managedInstanceSubnetNetworkIntentPolicyPoller struct {
	client SubnetsClient
	id     commonids.SubnetId
}

func NewManagedInstanceSubnetNetworkIntentPolicyPoller(client SubnetsClient, id commonids.SubnetId) managedInstanceSubnetNetworkIntentPolicyPoller {
	return managedInstanceSubnetNetworkIntentPolicyPoller{
		client: client,
		id:     id,
	}
}

func (p managedInstanceSubnetNetworkIntentPolicyPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &managedInstanceSubnetNetworkIntentPolicyDeleted, nil
		}
		return nil, fmt.Errorf("retrieving subnet %s: %+v", p.id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.ServiceAssociationLinks == nil {
		return &managedInstanceSubnetNetworkIntentPolicyDeleted, nil
	}

	for _, serviceAssociationLink := range *resp.Model.Properties.ServiceAssociationLinks {
		if isSqlManagedInstanceServiceAssociationLink(serviceAssociationLink) {
			log.Printf("[DEBUG] Waiting for SQL Managed Instance Service Association Link %q to be removed from subnet %s", pointer.From(serviceAssociationLink.Name), p.id.ID())
			return &managedInstanceSubnetNetworkIntentPolicyExists, nil
		}
	}

	log.Printf("[DEBUG] No SQL Managed Instance Service Association Links found on subnet %s", p.id.ID())
	return &managedInstanceSubnetNetworkIntentPolicyDeleted, nil
}

func isSqlManagedInstanceServiceAssociationLink(input subnets.ServiceAssociationLink) bool {
	if input.Properties == nil {
		return false
	}

	linkedResourceType := strings.ToLower(pointer.From(input.Properties.LinkedResourceType))
	if linkedResourceType == "microsoft.sql/managedinstances" || linkedResourceType == "microsoft.sql/virtualclusters" {
		return true
	}

	link := strings.ToLower(pointer.From(input.Properties.Link))
	return strings.Contains(link, "/providers/microsoft.sql/managedinstances/") ||
		strings.Contains(link, "/providers/microsoft.sql/virtualclusters/")
}

var (
	managedInstanceSubnetNetworkIntentPolicyDeleted = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}

	managedInstanceSubnetNetworkIntentPolicyExists = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)
