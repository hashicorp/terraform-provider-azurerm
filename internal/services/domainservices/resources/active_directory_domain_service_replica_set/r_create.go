// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_replica_set

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceActiveDirectoryDomainServiceReplicaSetCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	domainServiceId, err := parse.DomainServiceID(d.Get("domain_service_id").(string))
	if err != nil {
		return err
	}
	if domainServiceId == nil {
		return fmt.Errorf("parsing ID for Domain Service Replica Set")
	}

	idsdk := domainservices.NewDomainServiceID(domainServiceId.SubscriptionId, domainServiceId.ResourceGroup, domainServiceId.Name)

	locks.ByName(domainServiceId.Name, DomainServiceResourceName)
	defer locks.UnlockByName(domainServiceId.Name, DomainServiceResourceName)

	domainService, err := client.Get(ctx, idsdk)
	if err != nil {
		if response.WasNotFound(domainService.HttpResponse) {
			return fmt.Errorf("could not find %s: %s", domainServiceId, err)
		}
		return fmt.Errorf("reading %s: %s", domainServiceId, err)
	}

	model := domainService.Model
	if model == nil {
		return fmt.Errorf("reading %s: returned with null model", domainServiceId)
	}

	if model.Properties == nil || model.Properties.ReplicaSets == nil || len(*model.Properties.ReplicaSets) == 0 {
		return fmt.Errorf("reading %s: returned with missing replica set information, expected at least 1 replica set", domainServiceId)
	}

	subnetId := d.Get("subnet_id").(string)
	replicaSets := *model.Properties.ReplicaSets

	exists := false
	for _, r := range replicaSets {
		if r.ReplicaSetId == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing ReplicaSetID", domainServiceId)
		}
		if r.SubnetId == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing SubnetID", domainServiceId)
		}

		// We assume that two replica sets cannot coexist in the same subnet
		if strings.EqualFold(subnetId, *r.SubnetId) {
			exists = true
			if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				// Generate an ID here since we only know it once we know the ReplicaSetID
				id := parse.NewDomainServiceReplicaSetID(domainServiceId.SubscriptionId, domainServiceId.ResourceGroup, domainServiceId.Name, *r.ReplicaSetId)
				return tf.ImportAsExistsError("azurerm_active_directory_domain_service_replica_set", id.ID())
			}
		}
	}

	loc := location.Normalize(d.Get("location").(string))
	if !exists {
		replicaSets = append(replicaSets, domainservices.ReplicaSet{
			Location: pointer.To(loc),
			SubnetId: pointer.To(subnetId),
		})
	}

	model.Properties.ReplicaSets = &replicaSets

	if err := client.CreateOrUpdateThenPoll(ctx, idsdk, *model); err != nil {
		return fmt.Errorf("creating/updating Replica Sets for %s: %+v", domainServiceId, err)
	}

	// We need to retrieve the domain service again to find out the new replica set ID
	domainService, err = client.Get(ctx, idsdk)
	if err != nil {
		if response.WasNotFound(domainService.HttpResponse) {
			return fmt.Errorf("could not find %s: %s", domainServiceId, err)
		}
		return fmt.Errorf("reading %s: %s", domainServiceId, err)
	}

	model = domainService.Model
	if model == nil {
		return fmt.Errorf("reading %s: returned with null model", domainServiceId)
	}

	if model.Properties == nil || model.Properties.ReplicaSets == nil || len(*model.Properties.ReplicaSets) == 0 {
		return fmt.Errorf("reading %s: returned with missing replica set information, expected at least 1 replica set", domainServiceId)
	}

	var id parse.DomainServiceReplicaSetId
	// Assuming that two replica sets cannot coexist in the same subnet, we identify our new replica set by its SubnetID
	for _, r := range *model.Properties.ReplicaSets {
		if r.ReplicaSetId == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing ReplicaSetID", domainServiceId)
		}
		if r.SubnetId == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing SubnetID", domainServiceId)
		}

		if strings.EqualFold(subnetId, *r.SubnetId) {
			// We found it!
			id = parse.NewDomainServiceReplicaSetID(domainServiceId.SubscriptionId, domainServiceId.ResourceGroup, domainServiceId.Name, *r.ReplicaSetId)
		}
	}

	if id.ReplicaSetName == "" {
		return fmt.Errorf("reading %s: the new replica set was not returned", domainServiceId)
	}

	d.SetId(id.ID())

	// Wait for all replica sets to become available with two domain controllers each before proceeding
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      helpers.DomainServiceControllerRefreshFunc(ctx, client, *domainServiceId, false),
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
		Timeout:      time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for both domain controllers to become available in all replica sets for %s: %+v", domainServiceId, err)
	}

	return resourceActiveDirectoryDomainServiceReplicaSetRead(d, meta)
}
