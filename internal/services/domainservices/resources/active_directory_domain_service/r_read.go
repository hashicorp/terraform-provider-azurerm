// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceActiveDirectoryDomainServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceID(d.Id())
	if err != nil {
		return err
	}

	idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.Name)

	resp, err := client.Get(ctx, idsdk)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("resource_id", model.Id)
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}

		if props := model.Properties; props != nil {
			d.Set("deployment_id", props.DeploymentId)
			d.Set("domain_name", props.DomainName)
			d.Set("sync_owner", props.SyncOwner)
			d.Set("tenant_id", props.TenantId)
			d.Set("version", props.Version)
			d.Set("domain_configuration_type", props.DomainConfigurationType)

			d.Set("filtered_sync_enabled", false)
			if props.FilteredSync != nil && *props.FilteredSync == domainservices.FilteredSyncEnabled {
				d.Set("filtered_sync_enabled", true)
			}

			d.Set("sku", props.Sku)

			if err := d.Set("notifications", flattenDomainServiceNotifications(props.NotificationSettings)); err != nil {
				return fmt.Errorf("setting `notifications`: %+v", err)
			}

			var initialReplicaSet any
			replicaSets := helpers.FlattenDomainServiceReplicaSets(props.ReplicaSets)

			// Determine the initial replica set. This is why we need to include InitialReplicaSetId in the resource ID,
			// without it we would not be able to reliably support importing.
			for _, replicaSetRaw := range replicaSets {
				replicaSet := replicaSetRaw.(map[string]any)
				if replicaSet["id"].(string) == id.InitialReplicaSetIdName {
					initialReplicaSet = replicaSetRaw
					break
				}
			}
			if initialReplicaSet == nil {
				// It's safest to error out here, since we don't want to wipe the initial replica set from state if it was deleted manually
				return fmt.Errorf("reading %s: could not determine initial replica set from API response", id)
			}
			if err := d.Set("initial_replica_set", []any{initialReplicaSet}); err != nil {
				return fmt.Errorf("setting `initial_replica_set`: %+v", err)
			}

			if err := d.Set("secure_ldap", flattenDomainServiceLdaps(d, props.LdapsSettings, false)); err != nil {
				return fmt.Errorf("setting `secure_ldap`: %+v", err)
			}

			if err := d.Set("security", flattenDomainServiceSecurity(props.DomainSecuritySettings)); err != nil {
				return fmt.Errorf("setting `security`: %+v", err)
			}
		}
	}

	return nil
}
