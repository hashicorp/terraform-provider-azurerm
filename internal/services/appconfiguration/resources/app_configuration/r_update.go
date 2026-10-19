// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/replicas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppConfigurationUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationstores.ParseConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	update := configurationstores.ConfigurationStoreUpdateParameters{}

	if d.HasChange("sku") {
		update.Sku = &configurationstores.Sku{
			Name: d.Get("sku").(string),
		}
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]any)
		update.Tags = tags.Expand(t)
	}

	if d.HasChange("identity") {
		identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		update.Identity = identity
	}

	if d.HasChange("data_plane_proxy_authentication_mode") {
		if update.Properties == nil {
			update.Properties = &configurationstores.ConfigurationStorePropertiesUpdateParameters{}
		}

		props := update.Properties
		if props.DataPlaneProxy == nil {
			props.DataPlaneProxy = &configurationstores.DataPlaneProxyProperties{}
		}
		props.DataPlaneProxy.AuthenticationMode = pointer.ToEnum[configurationstores.AuthenticationMode](d.Get("data_plane_proxy_authentication_mode").(string))
	}

	if d.HasChange("data_plane_proxy_private_link_delegation_enabled") {
		if update.Properties == nil {
			update.Properties = &configurationstores.ConfigurationStorePropertiesUpdateParameters{}
		}

		props := update.Properties
		if props.DataPlaneProxy == nil {
			props.DataPlaneProxy = &configurationstores.DataPlaneProxyProperties{}
		}

		privLinkDelegation := configurationstores.PrivateLinkDelegationDisabled
		if d.Get("data_plane_proxy_private_link_delegation_enabled").(bool) {
			privLinkDelegation = configurationstores.PrivateLinkDelegationEnabled
		}
		props.DataPlaneProxy.PrivateLinkDelegation = &privLinkDelegation
	}

	if d.HasChange("encryption") {
		if update.Properties == nil {
			update.Properties = &configurationstores.ConfigurationStorePropertiesUpdateParameters{}
		}
		update.Properties.Encryption = expandAppConfigurationEncryption(d.Get("encryption").([]any))
	}

	if d.HasChange("local_auth_enabled") {
		if update.Properties == nil {
			update.Properties = &configurationstores.ConfigurationStorePropertiesUpdateParameters{}
		}
		update.Properties.DisableLocalAuth = pointer.To(!d.Get("local_auth_enabled").(bool))
	}

	if d.HasChange("public_network_access") {
		if update.Properties == nil {
			update.Properties = &configurationstores.ConfigurationStorePropertiesUpdateParameters{}
		}

		publicNetworkAccessValue, publicNetworkAccessNotEmpty := d.GetOk("public_network_access")
		if publicNetworkAccessNotEmpty {
			update.Properties.PublicNetworkAccess = parsePublicNetworkAccess(publicNetworkAccessValue.(string))
		}
	}

	if d.HasChange("purge_protection_enabled") {
		if update.Properties == nil {
			update.Properties = &configurationstores.ConfigurationStorePropertiesUpdateParameters{}
		}

		newValue := d.Get("purge_protection_enabled").(bool)
		oldValue := pointer.From(existing.Model.Properties.EnablePurgeProtection)

		if oldValue && !newValue {
			return fmt.Errorf("updating %s: once Purge Protection has been Enabled it's not possible to disable it", *id)
		}
		update.Properties.EnablePurgeProtection = pointer.To(d.Get("purge_protection_enabled").(bool))
	}

	if err := client.UpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if d.HasChange("replica") {
		replicaClient := meta.(*clients.Client).AppConfiguration.ReplicasClient
		operationsClient := meta.(*clients.Client).AppConfiguration.OperationsClient

		// check if a replica has been removed from config and if so, delete it
		deleteReplicaIds := make([]replicas.ReplicaId, 0)
		unchangedReplicaNames := make(map[string]struct{}, 0)
		oldReplicas, newReplicas := d.GetChange("replica")
		for _, oldReplica := range oldReplicas.(*pluginsdk.Set).List() {
			isRemoved := true
			oldReplicaMap := oldReplica.(map[string]any)

			for _, newReplica := range newReplicas.(*pluginsdk.Set).List() {
				newReplicaMap := newReplica.(map[string]any)

				if strings.EqualFold(oldReplicaMap["name"].(string), newReplicaMap["name"].(string)) && strings.EqualFold(location.Normalize(oldReplicaMap["location"].(string)), location.Normalize(newReplicaMap["location"].(string))) {
					unchangedReplicaNames[oldReplicaMap["name"].(string)] = struct{}{}
					isRemoved = false
					break
				}
			}

			if isRemoved {
				deleteReplicaIds = append(deleteReplicaIds, replicas.NewReplicaID(id.SubscriptionId, id.ResourceGroupName, id.ConfigurationStoreName, oldReplicaMap["name"].(string)))
			}
		}

		if err := deleteReplicas(ctx, replicaClient, operationsClient, deleteReplicaIds); err != nil {
			return err
		}

		expandedReplicas, err := expandAppConfigurationReplicas(d.Get("replica").(*pluginsdk.Set).List(), id.ConfigurationStoreName, location.Normalize(existing.Model.Location))
		if err != nil {
			return fmt.Errorf("expanding `replica`: %+v", err)
		}

		// check if a replica has been added or an existing one changed its location, (re)create it
		for _, replica := range *expandedReplicas {
			if _, isUnchanged := unchangedReplicaNames[*replica.Name]; isUnchanged {
				continue
			}

			replicaId := replicas.NewReplicaID(id.SubscriptionId, id.ResourceGroupName, id.ConfigurationStoreName, *replica.Name)

			existingReplica, err := replicaClient.Get(ctx, replicaId)
			if err != nil {
				if !response.WasNotFound(existingReplica.HttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", replicaId, err)
				}
			}

			if !response.WasNotFound(existingReplica.HttpResponse) {
				return fmt.Errorf("updating %s: replica %s already exists", *id, replicaId)
			}

			if err = replicaClient.CreateThenPoll(ctx, replicaId, replica); err != nil {
				return fmt.Errorf("creating %s: %+v", replicaId, err)
			}
		}
	}

	return resourceAppConfigurationRead(d, meta)
}
