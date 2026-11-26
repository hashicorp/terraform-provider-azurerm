// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_replica_set

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceActiveDirectoryDomainServiceReplicaSetDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceReplicaSetID(d.Id())
	if err != nil {
		return err
	}

	idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

	domainService, err := client.Get(ctx, idsdk)
	if err != nil {
		if response.WasNotFound(domainService.HttpResponse) {
			return fmt.Errorf("deleting %s: domain service was not found: %s", id, err)
		}
		return err
	}

	model := domainService.Model
	if model == nil {
		return fmt.Errorf("reading %s: returned with null model", id)
	}

	if model.Properties == nil || model.Properties.ReplicaSets == nil || len(*model.Properties.ReplicaSets) == 0 {
		return fmt.Errorf("reading %s: returned with missing replica set information, expected at least 1 replica set", id)
	}

	replicaSets := *model.Properties.ReplicaSets

	newReplicaSets := make([]domainservices.ReplicaSet, 0)
	for _, r := range replicaSets {
		if r.ReplicaSetId == nil {
			return fmt.Errorf("deleting %s: a replica set was returned with a missing ReplicaSetID", id)
		}

		if *r.ReplicaSetId == id.ReplicaSetName {
			continue
		}

		newReplicaSets = append(newReplicaSets, r)
	}

	if len(replicaSets) == len(newReplicaSets) {
		return fmt.Errorf("deleting %s: could not determine which replica set to remove", id)
	}

	model.Properties.ReplicaSets = &newReplicaSets

	if err := client.CreateOrUpdateThenPoll(ctx, idsdk, *model); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	// Wait for all replica sets to become available with two domain controllers each before proceeding
	// Generate a partial DomainServiceId since we don't need to know the initial replica set ID here
	domainServiceId := parse.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName, "")
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      helpers.DomainServiceControllerRefreshFunc(ctx, client, domainServiceId, true),
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
		Timeout:      time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for replica sets to finish updating for %s: %+v", domainServiceId, err)
	}

	return nil
}
