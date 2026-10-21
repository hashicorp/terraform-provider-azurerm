// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/deletedconfigurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/replicas"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppConfigurationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	deletedConfigurationStoresClient := meta.(*clients.Client).AppConfiguration.DeletedConfigurationStoresClient

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resourceId := configurationstores.NewConfigurationStoreID(subscriptionId, resourceGroup, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, resourceId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", resourceId, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_app_configuration", resourceId.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))

	recoverSoftDeleted := false
	if meta.(*clients.Client).Features.AppConfiguration.RecoverSoftDeleted {
		deletedConfigurationStoresId := deletedconfigurationstores.NewDeletedConfigurationStoreID(subscriptionId, location, name)
		deleted, err := deletedConfigurationStoresClient.ConfigurationStoresGetDeleted(ctx, deletedConfigurationStoresId)
		if err != nil {
			if response.WasStatusCode(deleted.HttpResponse, http.StatusForbidden) {
				return errors.New(userIsMissingNecessaryPermission(name, location))
			}
			if !response.WasNotFound(deleted.HttpResponse) {
				return fmt.Errorf("checking for presence of deleted %s: %+v", deletedConfigurationStoresId, err)
			}
			// if the soft deleted is not found, skip the recovering
		} else {
			recoverSoftDeleted = true
		}
	}

	privLinkDelegation := configurationstores.PrivateLinkDelegationDisabled
	if d.Get("data_plane_proxy_private_link_delegation_enabled").(bool) {
		privLinkDelegation = configurationstores.PrivateLinkDelegationEnabled
	}

	parameters := configurationstores.ConfigurationStore{
		Location: location,
		Sku: configurationstores.Sku{
			Name: d.Get("sku").(string),
		},
		Properties: &configurationstores.ConfigurationStoreProperties{
			DataPlaneProxy: &configurationstores.DataPlaneProxyProperties{
				AuthenticationMode:    pointer.ToEnum[configurationstores.AuthenticationMode](d.Get("data_plane_proxy_authentication_mode").(string)),
				PrivateLinkDelegation: &privLinkDelegation,
			},
			EnablePurgeProtection: pointer.To(d.Get("purge_protection_enabled").(bool)),
			DisableLocalAuth:      pointer.To(!d.Get("local_auth_enabled").(bool)),
			Encryption:            expandAppConfigurationEncryption(d.Get("encryption").([]any)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.Get("soft_delete_retention_days").(int); ok && v != 7 {
		parameters.Properties.SoftDeleteRetentionInDays = pointer.To(int64(v))
	}

	if recoverSoftDeleted {
		parameters.Properties.CreateMode = pointer.To(configurationstores.CreateModeRecover)
	}

	publicNetworkAccessValue, publicNetworkAccessNotEmpty := d.GetOk("public_network_access")

	if publicNetworkAccessNotEmpty {
		parameters.Properties.PublicNetworkAccess = parsePublicNetworkAccess(publicNetworkAccessValue.(string))
	}

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	parameters.Identity = identity

	if err := client.CreateCallbackThenPoll(ctx, resourceId, parameters, sdk.SetIDCallback(meta, &resourceId, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())

	resp, err := client.Get(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", resourceId, err)
	}
	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Endpoint == nil {
		return fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", resourceId)
	}
	meta.(*clients.Client).AppConfiguration.AddToCache(resourceId, *resp.Model.Properties.Endpoint)

	expandedReplicas, err := expandAppConfigurationReplicas(d.Get("replica").(*pluginsdk.Set).List(), name, location)
	if err != nil {
		return fmt.Errorf("expanding `replica`: %+v", err)
	}

	replicaClient := meta.(*clients.Client).AppConfiguration.ReplicasClient
	for _, replica := range *expandedReplicas {
		replicaId := replicas.NewReplicaID(resourceId.SubscriptionId, resourceId.ResourceGroupName, resourceId.ConfigurationStoreName, *replica.Name)

		if err := replicaClient.CreateThenPoll(ctx, replicaId, replica); err != nil {
			return fmt.Errorf("creating %s: %+v", replicaId, err)
		}
	}

	return resourceAppConfigurationRead(d, meta)
}

func expandAppConfigurationEncryption(input []any) *configurationstores.EncryptionProperties {
	result := &configurationstores.EncryptionProperties{
		KeyVaultProperties: &configurationstores.KeyVaultProperties{},
	}

	if len(input) == 0 || input[0] == nil {
		return result
	}

	encryptionParam := input[0].(map[string]any)

	if v, ok := encryptionParam["identity_client_id"].(string); ok && v != "" {
		result.KeyVaultProperties.IdentityClientId = &v
	}
	if v, ok := encryptionParam["key_vault_key_identifier"].(string); ok && v != "" {
		result.KeyVaultProperties.KeyIdentifier = &v
	}
	return result
}

func expandAppConfigurationReplicas(input []any, configurationStoreName, configurationStoreLocation string) (*[]replicas.Replica, error) {
	result := make([]replicas.Replica, 0)

	// check if there are duplicated replica names or locations
	// location cannot be same as original configuration store and other replicas
	locationSet := make(map[string]string, 0)
	replicaNameSet := make(map[string]struct{}, 0)

	for _, v := range input {
		replica := v.(map[string]any)
		replicaName := replica["name"].(string)
		replicaLocation := location.Normalize(replica["location"].(string))
		if strings.EqualFold(replicaLocation, configurationStoreLocation) {
			return nil, fmt.Errorf("location (%q) of replica %q is duplicated with original configuration store %q", replicaLocation, replicaName, configurationStoreName)
		}

		if name, ok := locationSet[replicaLocation]; ok {
			return nil, fmt.Errorf("location (%q) of replica %q is duplicated with replica %q", replicaLocation, replicaName, name)
		}
		locationSet[replicaLocation] = replicaName

		normalizedReplicaName := strings.ToLower(replicaName)
		if _, ok := replicaNameSet[normalizedReplicaName]; ok {
			return nil, fmt.Errorf("replica name %q is duplicated", replicaName)
		}
		replicaNameSet[normalizedReplicaName] = struct{}{}

		if len(replicaName)+len(configurationStoreName) > 60 {
			return nil, fmt.Errorf("replica name %q is too long, the total length of replica name and configuration store name should be less or equal than 60", replicaName)
		}

		result = append(result, replicas.Replica{
			Name:     pointer.To(replicaName),
			Location: pointer.To(replicaLocation),
		})
	}

	return &result, nil
}
