// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceActiveDirectoryDomainServiceCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resourceErrorName := fmt.Sprintf("Domain Service (Name: %q, Resource Group: %q)", name, resourceGroup)

	locks.ByName(name, DomainServiceResourceName)
	defer locks.UnlockByName(name, DomainServiceResourceName)

	// If this is a new resource, we cannot determine the resource ID until after it has been created since we need to
	// know the ID of the first replica set.
	var id *parse.DomainServiceId

	idsdk := domainservices.NewDomainServiceID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, idsdk)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", resourceErrorName, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				// Parse the replica sets and assume the first one returned to be the initial replica set
				// This is a best effort and the user can choose any replica set if they structure their config accordingly
				model := existing.Model
				if model == nil {
					return fmt.Errorf("checking for presence of existing %s: API response contained nil or missing model", resourceErrorName)
				}
				props := model.Properties
				if props == nil {
					return fmt.Errorf("checking for presence of existing %s: API response contained nil or missing properties", resourceErrorName)
				}
				replicaSets := helpers.FlattenDomainServiceReplicaSets(props.ReplicaSets)
				if len(replicaSets) == 0 {
					return fmt.Errorf("checking for presence of existing %s: API response contained nil or missing replica set details", resourceErrorName)
				}
				initialReplicaSetId := replicaSets[0].(map[string]any)["id"].(string)
				id := parse.NewDomainServiceID(subscriptionId, resourceGroup, name, initialReplicaSetId)

				return tf.ImportAsExistsError(DomainServiceResourceName, id.ID())
			}
		}
	} else {
		var err error
		id, err = parse.DomainServiceID(d.Id())
		if err != nil {
			return fmt.Errorf("preparing update for %s: %+v", resourceErrorName, err)
		}
		if id == nil {
			return fmt.Errorf("preparing update for %s: resource ID could not be parsed", resourceErrorName)
		}
	}

	loc := location.Normalize(d.Get("location").(string))
	filteredSync := domainservices.FilteredSyncDisabled
	if d.Get("filtered_sync_enabled").(bool) {
		filteredSync = domainservices.FilteredSyncEnabled
	}

	domainService := domainservices.DomainService{
		Properties: &domainservices.DomainServiceProperties{
			DomainName:             pointer.To(d.Get("domain_name").(string)),
			DomainSecuritySettings: expandDomainServiceSecurity(d.Get("security").([]any)),
			FilteredSync:           &filteredSync,
			LdapsSettings:          expandDomainServiceLdaps(d.Get("secure_ldap").([]any)),
			NotificationSettings:   expandDomainServiceNotifications(d.Get("notifications").([]any)),
			Sku:                    pointer.To(d.Get("sku").(string)),
		},
		Location: pointer.To(loc),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v := d.Get("domain_configuration_type").(string); v != "" {
		domainService.Properties.DomainConfigurationType = &v
	}

	if d.IsNewResource() {
		// On resource creation, specify the initial replica set.
		// No provision is made for changing the initial replica set, it should remain intact for the resource to function properly
		replicaSets := []domainservices.ReplicaSet{
			{
				Location: pointer.To(loc),
				SubnetId: pointer.To(d.Get("initial_replica_set.0.subnet_id").(string)),
			},
		}
		domainService.Properties.ReplicaSets = &replicaSets
	}

	// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
	if err := client.CreateOrUpdateThenPoll(ctx, idsdk, domainService); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceErrorName, err)
	}

	// Retrieve the domain service to discover the unique ID for the initial replica set, which should not subsequently change
	if d.IsNewResource() {
		resp, err := client.Get(ctx, idsdk)
		if err != nil {
			return fmt.Errorf("retrieving %s after creating: %+v", resourceErrorName, err)
		}
		model := resp.Model
		if model == nil {
			return fmt.Errorf("%s returned with no model", resourceErrorName)
		}
		props := model.Properties
		if props == nil {
			return fmt.Errorf("%s returned with no properties", resourceErrorName)
		}
		if props.ReplicaSets == nil {
			return fmt.Errorf("%s returned with no replica set details", resourceErrorName)
		}

		replicaSets := helpers.FlattenDomainServiceReplicaSets(props.ReplicaSets)
		if replicaSetCount := len(replicaSets); replicaSetCount != 1 {
			return fmt.Errorf("unexpected number of replica sets for %s: expected 1, saw %d", resourceErrorName, replicaSetCount)
		}

		// Once we know the initial replica set ID, we can build a resource ID
		initialReplicaSetId := replicaSets[0].(map[string]any)["id"].(string)
		newId := parse.NewDomainServiceID(subscriptionId, resourceGroup, name, initialReplicaSetId)
		id = &newId
		d.SetId(id.ID())

		if err := d.Set("initial_replica_set", []any{replicaSets[0]}); err != nil {
			return fmt.Errorf("setting `initial_replica_set` after creating resource: %+v", err)
		}
	}

	if id == nil {
		return fmt.Errorf("after creating/updating %s: id was unexpectedly nil", resourceErrorName)
	}

	// A fully deployed domain service has 2 domain controllers per replica set, but the create operation completes early before the DCs are online.
	// The domain service is still provisioning and further operations are blocked until both DCs are up and ready.
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      helpers.DomainServiceControllerRefreshFunc(ctx, client, *id, false),
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
		Timeout:      time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for both domain controllers to become available in initial replica set for %s: %+v", id, err)
	}

	return resourceActiveDirectoryDomainServiceRead(d, meta)
}
