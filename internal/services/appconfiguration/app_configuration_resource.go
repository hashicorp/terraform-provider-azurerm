// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/deletedconfigurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/operations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/replicas"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppConfigurationCreate,
		Read:   resourceAppConfigurationRead,
		Update: resourceAppConfigurationUpdate,
		Delete: resourceAppConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := configurationstores.ParseConfigurationStoreID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// sku cannot be downgraded from a production tier (`premium` or `standard`) to a non-production tier (`developer` or `free`), or downgraded from `developer` to `free`
			// https://learn.microsoft.com/azure/azure-app-configuration/faq#can-i-upgrade-or-downgrade-an-app-configuration-store
			pluginsdk.ForceNewIfChange("sku", func(ctx context.Context, old, new, meta interface{}) bool {
				return ((old == "premium" || old == "standard") && new == "developer") || new == "free"
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
				authMode := d.Get("data_plane_proxy_authentication_mode").(string)
				privLinkDelegation := d.Get("data_plane_proxy_private_link_delegation_enabled").(bool)

				if authMode == string(configurationstores.AuthenticationModeLocal) && privLinkDelegation {
					return errors.New("`data_plane_proxy_private_link_delegation_enabled` cannot be set to `true` when `data_plane_proxy_authentication_mode` is `Local`")
				}

				return nil
			}),
		),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ConfigurationStoreName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"data_plane_proxy_authentication_mode": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(configurationstores.AuthenticationModeLocal),
				ValidateFunc: validation.StringInSlice(configurationstores.PossibleValuesForAuthenticationMode(), false),
			},

			"data_plane_proxy_private_link_delegation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"encryption": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"identity_client_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
						"key_vault_key_identifier": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},
					},
				},
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"public_network_access": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringInSlice(configurationstores.PossibleValuesForPublicNetworkAccess(), true),
			},

			"purge_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"replica": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MinItems: 1,
				Set:      resourceConfigurationStoreReplicaHash,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ConfigurationStoreReplicaName,
						},
						"location": commonschema.LocationWithoutForceNew(),
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"endpoint": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			// `sku` is not enum, https://github.com/Azure/azure-rest-api-specs/issues/23902
			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "free",
				ValidateFunc: validation.StringInSlice([]string{
					"free",
					"developer",
					"standard",
					"premium",
				}, false),
			},

			"soft_delete_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      7,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 7),
				DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
					return old == "0"
				},
			},

			"tags": commonschema.Tags(),

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_read_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"primary_write_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"secondary_read_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"secondary_write_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func resourceAppConfigurationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	deletedConfigurationStoresClient := meta.(*clients.Client).AppConfiguration.DeletedConfigurationStoresClient

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resourceId := configurationstores.NewConfigurationStoreID(subscriptionId, resourceGroup, name)
	existing, err := client.Get(ctx, resourceId)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", resourceId, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_app_configuration", resourceId.ID())
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
			log.Printf("[DEBUG] Soft Deleted App Configuration exists, marked for recover")
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
				AuthenticationMode:    pointer.To(configurationstores.AuthenticationMode(d.Get("data_plane_proxy_authentication_mode").(string))),
				PrivateLinkDelegation: &privLinkDelegation,
			},
			EnablePurgeProtection: pointer.To(d.Get("purge_protection_enabled").(bool)),
			DisableLocalAuth:      pointer.To(!d.Get("local_auth_enabled").(bool)),
			Encryption:            expandAppConfigurationEncryption(d.Get("encryption").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.Get("soft_delete_retention_days").(int); ok && v != 7 {
		parameters.Properties.SoftDeleteRetentionInDays = pointer.To(int64(v))
	}

	if recoverSoftDeleted {
		t := configurationstores.CreateModeRecover
		parameters.Properties.CreateMode = &t
	}

	publicNetworkAccessValue, publicNetworkAccessNotEmpty := d.GetOk("public_network_access")

	if publicNetworkAccessNotEmpty {
		parameters.Properties.PublicNetworkAccess = parsePublicNetworkAccess(publicNetworkAccessValue.(string))
	}

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	parameters.Identity = identity

	if err := client.CreateThenPoll(ctx, resourceId, parameters); err != nil {
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

func resourceAppConfigurationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration update.")
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
		t := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(t)
	}

	if d.HasChange("identity") {
		identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
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
		props.DataPlaneProxy.AuthenticationMode = pointer.To(configurationstores.AuthenticationMode(d.Get("data_plane_proxy_authentication_mode").(string)))
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
		update.Properties.Encryption = expandAppConfigurationEncryption(d.Get("encryption").([]interface{}))
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
		oldValue := false
		if existing.Model.Properties.EnablePurgeProtection != nil {
			oldValue = *existing.Model.Properties.EnablePurgeProtection
		}

		if oldValue && !newValue {
			return fmt.Errorf("updating %s: once Purge Protection has been Enabled it's not possible to disable it", *id)
		}
		update.Properties.EnablePurgeProtection = pointer.To(d.Get("purge_protection_enabled").(bool))
	}

	if d.HasChange("public_network_enabled") {
		v := d.GetRawConfig().AsValueMap()["public_network_access_enabled"]
		if v.IsNull() && existing.Model.Properties.SoftDeleteRetentionInDays != nil {
			return fmt.Errorf("updating %s: once Public Network Access has been explicitly Enabled or Disabled it's not possible to unset it to which means Automatic", *id)
		}

		if update.Properties == nil {
			update.Properties = &configurationstores.ConfigurationStorePropertiesUpdateParameters{}
		}

		publicNetworkAccess := configurationstores.PublicNetworkAccessEnabled
		if v.False() {
			publicNetworkAccess = configurationstores.PublicNetworkAccessDisabled
		}
		update.Properties.PublicNetworkAccess = &publicNetworkAccess
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
			oldReplicaMap := oldReplica.(map[string]interface{})

			for _, newReplica := range newReplicas.(*pluginsdk.Set).List() {
				newReplicaMap := newReplica.(map[string]interface{})

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

func resourceAppConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationstores.ParseConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	resultPage, err := client.ListKeysComplete(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving access keys for %s: %+v", *id, err)
	}

	d.Set("name", id.ConfigurationStoreName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("sku", model.Sku.Name)

		if props := model.Properties; props != nil {
			if dataPlaneProxy := props.DataPlaneProxy; dataPlaneProxy != nil {
				d.Set("data_plane_proxy_authentication_mode", string(pointer.From(dataPlaneProxy.AuthenticationMode)))
				d.Set("data_plane_proxy_private_link_delegation_enabled", pointer.From(dataPlaneProxy.PrivateLinkDelegation) == configurationstores.PrivateLinkDelegationEnabled)
			}

			d.Set("endpoint", props.Endpoint)
			d.Set("encryption", flattenAppConfigurationEncryption(props.Encryption))
			d.Set("public_network_access", string(pointer.From(props.PublicNetworkAccess)))

			localAuthEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !(*props.DisableLocalAuth)
			}

			d.Set("local_auth_enabled", localAuthEnabled)

			purgeProtectionEnabled := false
			if props.EnablePurgeProtection != nil {
				purgeProtectionEnabled = *props.EnablePurgeProtection
			}
			d.Set("purge_protection_enabled", purgeProtectionEnabled)

			softDeleteRetentionDays := 0
			if props.SoftDeleteRetentionInDays != nil {
				softDeleteRetentionDays = int(*props.SoftDeleteRetentionInDays)
			}
			d.Set("soft_delete_retention_days", softDeleteRetentionDays)
		}

		accessKeys := flattenAppConfigurationAccessKeys(resultPage.Items)
		d.Set("primary_read_key", accessKeys.primaryReadKey)
		d.Set("primary_write_key", accessKeys.primaryWriteKey)
		d.Set("secondary_read_key", accessKeys.secondaryReadKey)
		d.Set("secondary_write_key", accessKeys.secondaryWriteKey)

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		replicasClient := meta.(*clients.Client).AppConfiguration.ReplicasClient
		resp, err := replicasClient.ListByConfigurationStoreComplete(ctx, replicas.NewConfigurationStoreID(id.SubscriptionId, id.ResourceGroupName, id.ConfigurationStoreName))
		if err != nil {
			return fmt.Errorf("retrieving replicas for %s: %+v", *id, err)
		}

		replica, err := flattenAppConfigurationReplicas(resp.Items)
		if err != nil {
			return fmt.Errorf("flattening replicas for %s: %+v", *id, err)
		}
		d.Set("replica", replica)

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceAppConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	deletedConfigurationStoresClient := meta.(*clients.Client).AppConfiguration.DeletedConfigurationStoresClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id, err := configurationstores.ParseConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %q: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %q: `properties` was nil", *id)
	}

	purgeProtectionEnabled := false
	if ppe := existing.Model.Properties.EnablePurgeProtection; ppe != nil {
		purgeProtectionEnabled = *ppe
	}
	softDeleteEnabled := false
	if sde := existing.Model.Properties.SoftDeleteRetentionInDays; sde != nil && *sde > 0 {
		softDeleteEnabled = true
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if meta.(*clients.Client).Features.AppConfiguration.PurgeSoftDeleteOnDestroy && softDeleteEnabled {
		deletedId := deletedconfigurationstores.NewDeletedConfigurationStoreID(subscriptionId, existing.Model.Location, id.ConfigurationStoreName)

		// AppConfiguration with Purge Protection Enabled cannot be deleted unless done by Azure
		if purgeProtectionEnabled {
			deletedInfo, err := deletedConfigurationStoresClient.ConfigurationStoresGetDeleted(ctx, deletedId)
			if err != nil {
				return fmt.Errorf("while purging the soft-deleted, retrieving the Deletion Details for %s: %+v", *id, err)
			}

			if deletedInfo.Model != nil && deletedInfo.Model.Properties != nil && deletedInfo.Model.Properties.DeletionDate != nil && deletedInfo.Model.Properties.ScheduledPurgeDate != nil {
				log.Printf("[DEBUG] The App Configuration %q has Purge Protection Enabled and was deleted on %q. Azure will purge this on %q",
					id.ConfigurationStoreName, *deletedInfo.Model.Properties.DeletionDate, *deletedInfo.Model.Properties.ScheduledPurgeDate)
			} else {
				log.Printf("[DEBUG] The App Configuration %q has Purge Protection Enabled and will be purged automatically by Azure", id.ConfigurationStoreName)
			}
			return nil
		}

		log.Printf("[DEBUG]  %q marked for purge - executing purge", id.ConfigurationStoreName)

		if _, err := deletedConfigurationStoresClient.ConfigurationStoresPurgeDeleted(ctx, deletedId); err != nil {
			return fmt.Errorf("purging %s: %+v", *id, err)
		}

		// The PurgeDeleted API is a POST which returns a 200 with no body and nothing to poll on, so we'll need
		// a custom poller to poll until the LRO returns a 404
		pollerType := &purgeDeletedPoller{
			client: deletedConfigurationStoresClient,
			id:     deletedId,
		}
		poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling after purging for %s: %+v", *id, err)
		}

		// retry checkNameAvailability until the name is released by purged app configuration, see https://github.com/Azure/AppConfiguration/issues/677
		operationsClient := meta.(*clients.Client).AppConfiguration.OperationsClient
		if err = resourceConfigurationStoreWaitForNameAvailable(ctx, operationsClient, *id); err != nil {
			return err
		}
		log.Printf("[DEBUG] Purged AppConfiguration %q.", id.ConfigurationStoreName)
	}

	meta.(*clients.Client).AppConfiguration.RemoveFromCache(*id)

	return nil
}

var _ pollers.PollerType = &purgeDeletedPoller{}

type purgeDeletedPoller struct {
	client *deletedconfigurationstores.DeletedConfigurationStoresClient
	id     deletedconfigurationstores.DeletedConfigurationStoreId
}

func (p *purgeDeletedPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.ConfigurationStoresGetDeleted(ctx, p.id)

	status := "dropped connection"
	if resp.HttpResponse != nil {
		status = fmt.Sprintf("%d", resp.HttpResponse.StatusCode)
	}

	if response.WasNotFound(resp.HttpResponse) {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			Status: pollers.PollingStatusSucceeded,
		}, nil
	}

	if response.WasStatusCode(resp.HttpResponse, http.StatusOK) {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			Status: pollers.PollingStatusInProgress,
		}, nil
	}

	return nil, fmt.Errorf("unexpected status %q: %+v", status, err)
}

type flattenedAccessKeys struct {
	primaryReadKey    []interface{}
	primaryWriteKey   []interface{}
	secondaryReadKey  []interface{}
	secondaryWriteKey []interface{}
}

func expandAppConfigurationEncryption(input []interface{}) *configurationstores.EncryptionProperties {
	result := &configurationstores.EncryptionProperties{
		KeyVaultProperties: &configurationstores.KeyVaultProperties{},
	}

	if len(input) == 0 || input[0] == nil {
		return result
	}

	encryptionParam := input[0].(map[string]interface{})

	if v, ok := encryptionParam["identity_client_id"].(string); ok && v != "" {
		result.KeyVaultProperties.IdentityClientId = &v
	}
	if v, ok := encryptionParam["key_vault_key_identifier"].(string); ok && v != "" {
		result.KeyVaultProperties.KeyIdentifier = &v
	}
	return result
}

func expandAppConfigurationReplicas(input []interface{}, configurationStoreName, configurationStoreLocation string) (*[]replicas.Replica, error) {
	result := make([]replicas.Replica, 0)

	// check if there are duplicated replica names or locations
	// location cannot be same as original configuration store and other replicas
	locationSet := make(map[string]string, 0)
	replicaNameSet := make(map[string]struct{}, 0)

	for _, v := range input {
		replica := v.(map[string]interface{})
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

func flattenAppConfigurationAccessKeys(values []configurationstores.ApiKey) flattenedAccessKeys {
	result := flattenedAccessKeys{
		primaryReadKey:    make([]interface{}, 0),
		primaryWriteKey:   make([]interface{}, 0),
		secondaryReadKey:  make([]interface{}, 0),
		secondaryWriteKey: make([]interface{}, 0),
	}

	for _, value := range values {
		if value.Name == nil || value.ReadOnly == nil {
			continue
		}

		accessKey := flattenAppConfigurationAccessKey(value)
		name := *value.Name
		readOnly := *value.ReadOnly

		if strings.HasPrefix(strings.ToLower(name), "primary") {
			if readOnly {
				result.primaryReadKey = accessKey
			} else {
				result.primaryWriteKey = accessKey
			}
		}

		if strings.HasPrefix(strings.ToLower(name), "secondary") {
			if readOnly {
				result.secondaryReadKey = accessKey
			} else {
				result.secondaryWriteKey = accessKey
			}
		}
	}

	return result
}

func flattenAppConfigurationAccessKey(input configurationstores.ApiKey) []interface{} {
	connectionString := ""

	if input.ConnectionString != nil {
		connectionString = *input.ConnectionString
	}

	id := ""
	if input.Id != nil {
		id = *input.Id
	}

	secret := ""
	if input.Value != nil {
		secret = *input.Value
	}

	return []interface{}{
		map[string]interface{}{
			"connection_string": connectionString,
			"id":                id,
			"secret":            secret,
		},
	}
}

func parsePublicNetworkAccess(input string) *configurationstores.PublicNetworkAccess {
	vals := map[string]configurationstores.PublicNetworkAccess{
		"disabled": configurationstores.PublicNetworkAccessDisabled,
		"enabled":  configurationstores.PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v
	}

	// otherwise presume it's an undefined value and best-effort it
	out := configurationstores.PublicNetworkAccess(input)
	return &out
}

func userIsMissingNecessaryPermission(name, location string) string {
	return fmt.Sprintf(`
An existing soft-deleted App Configuration exists with the Name %q in the location %q, however
the credentials Terraform is using has insufficient permissions to check for an existing soft-deleted App Configuration.
You can opt out of this behaviour by using the "features" block (located within the "provider" block) - more information
can be found here:
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block
`, name, location)
}

func resourceConfigurationStoreWaitForNameAvailable(ctx context.Context, client *operations.OperationsClient, configurationStoreId configurationstores.ConfigurationStoreId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal error: context had no deadline")
	}
	state := &pluginsdk.StateChangeConf{
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 2,
		Pending:                   []string{"Unavailable"},
		Target:                    []string{"Available"},
		Refresh:                   resourceConfigurationStoreNameAvailabilityRefreshFunc(ctx, client, configurationStoreId),
		Timeout:                   time.Until(deadline),
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the name from %s to become available: %+v", configurationStoreId, err)
	}

	return nil
}

func resourceConfigurationStoreNameAvailabilityRefreshFunc(ctx context.Context, client *operations.OperationsClient, configurationStoreId configurationstores.ConfigurationStoreId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if the name for %s is available ..", configurationStoreId)

		subscriptionId := commonids.NewSubscriptionID(configurationStoreId.SubscriptionId)

		parameters := operations.CheckNameAvailabilityParameters{
			Name: configurationStoreId.ConfigurationStoreName,
			Type: operations.ConfigurationResourceTypeMicrosoftPointAppConfigurationConfigurationStores,
		}

		resp, err := client.CheckNameAvailability(ctx, subscriptionId, parameters)
		if err != nil {
			return resp, "Error", fmt.Errorf("retrieving Deployment: %+v", err)
		}

		if resp.Model == nil {
			return resp, "Error", fmt.Errorf("unexpected null model of %s", configurationStoreId)
		}

		if resp.Model.NameAvailable == nil {
			return resp, "Error", fmt.Errorf("unexpected null NameAvailable property of %s", configurationStoreId)
		}

		if !*resp.Model.NameAvailable {
			return resp, "Unavailable", nil
		}
		return resp, "Available", nil
	}
}

func deleteReplicas(ctx context.Context, replicaClient *replicas.ReplicasClient, operationClient *operations.OperationsClient, configurationStoreReplicaIds []replicas.ReplicaId) error {
	for _, configurationStoreReplicaId := range configurationStoreReplicaIds {
		log.Printf("[DEBUG] Deleting Replica %q", configurationStoreReplicaId)
		if err := replicaClient.DeleteThenPoll(ctx, configurationStoreReplicaId); err != nil {
			return fmt.Errorf("deleting replica %q: %+v", configurationStoreReplicaId, err)
		}
	}

	for _, configurationStoreReplicaId := range configurationStoreReplicaIds {
		if err := resourceConfigurationStoreReplicaWaitForNameAvailable(ctx, operationClient, configurationStoreReplicaId); err != nil {
			return fmt.Errorf("waiting for replica %q name to be released: %+v", configurationStoreReplicaId, err)
		}
	}

	return nil
}

func resourceConfigurationStoreReplicaWaitForNameAvailable(ctx context.Context, client *operations.OperationsClient, configurationStoreReplicaId replicas.ReplicaId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal error: context had no deadline")
	}
	state := &pluginsdk.StateChangeConf{
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 2,
		Pending:                   []string{"Unavailable"},
		Target:                    []string{"Available"},
		Refresh:                   resourceConfigurationStoreReplicaNameAvailabilityRefreshFunc(ctx, client, configurationStoreReplicaId),
		Timeout:                   time.Until(deadline),
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the name from %s to become available: %+v", configurationStoreReplicaId, err)
	}

	return nil
}

func resourceConfigurationStoreReplicaNameAvailabilityRefreshFunc(ctx context.Context, client *operations.OperationsClient, configurationStoreReplicaId replicas.ReplicaId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if the name for %s is available ..", configurationStoreReplicaId)

		subscriptionId := commonids.NewSubscriptionID(configurationStoreReplicaId.SubscriptionId)

		parameters := operations.CheckNameAvailabilityParameters{
			Name: fmt.Sprintf("%s-%s", configurationStoreReplicaId.ConfigurationStoreName, configurationStoreReplicaId.ReplicaName),
			Type: operations.ConfigurationResourceTypeMicrosoftPointAppConfigurationConfigurationStores,
		}

		resp, err := client.CheckNameAvailability(ctx, subscriptionId, parameters)
		if err != nil {
			return resp, "Error", fmt.Errorf("retrieving Deployment: %+v", err)
		}

		if resp.Model == nil {
			return resp, "Error", fmt.Errorf("unexpected null model of %s", configurationStoreReplicaId)
		}

		if resp.Model.NameAvailable == nil {
			return resp, "Error", fmt.Errorf("unexpected null NameAvailable property of %s", configurationStoreReplicaId)
		}

		if !*resp.Model.NameAvailable {
			return resp, "Unavailable", nil
		}
		return resp, "Available", nil
	}
}
