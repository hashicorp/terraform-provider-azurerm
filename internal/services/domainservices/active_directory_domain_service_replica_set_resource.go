// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package domainservices

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceActiveDirectoryDomainServiceReplicaSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceActiveDirectoryDomainServiceReplicaSetCreate,
		Read:   resourceActiveDirectoryDomainServiceReplicaSetRead,
		Delete: resourceActiveDirectoryDomainServiceReplicaSetDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(3 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(2 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DomainServiceReplicaSetID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"domain_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DomainServiceID,
			},

			"location": commonschema.Location(),

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"domain_controller_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"external_access_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceActiveDirectoryDomainServiceReplicaSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	for _, r := range replicaSets {
		if r.ReplicaSetId == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing ReplicaSetID", domainServiceId)
		}
		if r.SubnetId == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing SubnetID", domainServiceId)
		}

		// We assume that two replica sets cannot coexist in the same subnet
		if strings.EqualFold(subnetId, *r.SubnetId) {
			// Generate an ID here since we only know it once we know the ReplicaSetID
			id := parse.NewDomainServiceReplicaSetID(domainServiceId.SubscriptionId, domainServiceId.ResourceGroup, domainServiceId.Name, *r.ReplicaSetId)
			return tf.ImportAsExistsError("azurerm_active_directory_domain_service_replica_set", id.ID())
		}
	}

	loc := location.Normalize(d.Get("location").(string))
	replicaSets = append(replicaSets, domainservices.ReplicaSet{
		Location: utils.String(loc),
		SubnetId: utils.String(subnetId),
	})

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

	// Wait for all replica sets to become available with two domain controllers each before proceeding
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      domainServiceControllerRefreshFunc(ctx, client, *domainServiceId, false),
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
		Timeout:      time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for both domain controllers to become available in all replica sets for %s: %+v", domainServiceId, err)
	}

	d.SetId(id.ID())

	return resourceActiveDirectoryDomainServiceReplicaSetRead(d, meta)
}

func resourceActiveDirectoryDomainServiceReplicaSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceReplicaSetID(d.Id())
	if err != nil {
		return err
	}

	idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

	domainService, err := client.Get(ctx, idsdk)
	if err != nil {
		if response.WasNotFound(domainService.HttpResponse) {
			d.SetId("")
			return nil
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

	var (
		domainControllerIPAddresses []string
		externalAccessIPAddress     string
		loc                         string
		serviceStatus               string
		subnetId                    string
	)

	replicaSets := *model.Properties.ReplicaSets

	for _, r := range replicaSets {
		if r.ReplicaSetId == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing ReplicaSetID", id)
		}

		// ReplicaSetName in the ID struct is really the replica set ID
		if *r.ReplicaSetId == id.ReplicaSetName {
			if r.DomainControllerIPAddress != nil {
				domainControllerIPAddresses = *r.DomainControllerIPAddress
			}
			if r.ExternalAccessIPAddress != nil {
				externalAccessIPAddress = *r.ExternalAccessIPAddress
			}
			if r.Location != nil {
				loc = location.NormalizeNilable(r.Location)
			}
			if r.ServiceStatus != nil {
				serviceStatus = *r.ServiceStatus
			}
			if r.SubnetId != nil {
				subnetId = *r.SubnetId
			}
		}
	}

	d.Set("domain_controller_ip_addresses", domainControllerIPAddresses)
	d.Set("external_access_ip_address", externalAccessIPAddress)
	d.Set("location", loc)
	d.Set("service_status", serviceStatus)
	d.Set("subnet_id", subnetId)

	return nil
}

func resourceActiveDirectoryDomainServiceReplicaSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
		Refresh:      domainServiceControllerRefreshFunc(ctx, client, domainServiceId, true),
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
		Timeout:      time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for replica sets to finish updating for %s: %+v", domainServiceId, err)
	}

	return nil
}
