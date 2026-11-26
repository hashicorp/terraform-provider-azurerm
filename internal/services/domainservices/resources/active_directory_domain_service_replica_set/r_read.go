// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_replica_set

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceActiveDirectoryDomainServiceReplicaSetRead(d *pluginsdk.ResourceData, meta any) error {
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
