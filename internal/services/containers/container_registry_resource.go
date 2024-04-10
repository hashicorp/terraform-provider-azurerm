// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/operation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/replications"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/migration"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceContainerRegistry() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerRegistryCreate,
		Read:   resourceContainerRegistryRead,
		Update: resourceContainerRegistryUpdate,
		Delete: resourceContainerRegistryDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.RegistryV0ToV1{},
			1: migration.RegistryV1ToV2{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := registries.ParseRegistryID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceContainerRegistrySchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			sku := d.Get("sku").(string)

			geoReplications := d.Get("georeplications").([]interface{})
			// if locations have been specified for geo-replication then, the SKU has to be Premium
			if len(geoReplications) > 0 && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return fmt.Errorf("ACR geo-replication can only be applied when using the Premium Sku.")
			}

			// ensure location is different than any location of the geo-replication
			var geoReplicationLocations []string
			for _, v := range geoReplications {
				v := v.(map[string]interface{})
				geoReplicationLocations = append(geoReplicationLocations, azure.NormalizeLocation(v["location"]))
			}
			location := location.Normalize(d.Get("location").(string))
			for _, loc := range geoReplicationLocations {
				if loc == location {
					return fmt.Errorf("The `georeplications` list cannot contain the location where the Container Registry exists.")
				}
			}

			quarantinePolicyEnabled := d.Get("quarantine_policy_enabled").(bool)
			if quarantinePolicyEnabled && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return fmt.Errorf("ACR quarantine policy can only be applied when using the Premium Sku. If you are downgrading from a Premium SKU please unset quarantine_policy_enabled")
			}

			retentionPolicyEnabled, ok := d.GetOk("retention_policy.0.enabled")
			if ok && retentionPolicyEnabled.(bool) && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return fmt.Errorf("ACR retention policy can only be applied when using the Premium Sku. If you are downgrading from a Premium SKU please set retention_policy {}")
			}

			trustPolicyEnabled, ok := d.GetOk("trust_policy.0.enabled")
			if ok && trustPolicyEnabled.(bool) && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return fmt.Errorf("ACR trust policy can only be applied when using the Premium Sku. If you are downgrading from a Premium SKU please set trust_policy {}")
			}

			exportPolicyEnabled := d.Get("export_policy_enabled").(bool)
			if !exportPolicyEnabled {
				if !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
					return fmt.Errorf("ACR export policy can only be disabled when using the Premium Sku. If you are downgrading from a Premium SKU please unset `export_policy_enabled` or set `export_policy_enabled = true`")
				}
				if d.Get("public_network_access_enabled").(bool) {
					return fmt.Errorf("To disable export of artifacts, `public_network_access_enabled` must also be `false`")
				}
			}

			encryptionEnabled, ok := d.GetOk("encryption.0.enabled")
			if ok && encryptionEnabled.(bool) && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return fmt.Errorf("ACR encryption can only be applied when using the Premium Sku.")
			}

			// zone redundancy is only available for Premium Sku.
			zoneRedundancyEnabled, ok := d.GetOk("zone_redundancy_enabled")
			if ok && zoneRedundancyEnabled.(bool) && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return fmt.Errorf("ACR zone redundancy can only be applied when using the Premium Sku")
			}
			for _, loc := range geoReplications {
				loc := loc.(map[string]interface{})
				zoneRedundancyEnabled, ok := loc["zone_redundancy_enabled"]
				if ok && zoneRedundancyEnabled.(bool) && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
					return fmt.Errorf("ACR zone redundancy can only be applied when using the Premium Sku")
				}
			}

			// anonymous pull is only available for Standard/Premium Sku.
			if d.Get("anonymous_pull_enabled").(bool) && (!strings.EqualFold(sku, string(registries.SkuNameStandard)) && !strings.EqualFold(sku, string(registries.SkuNamePremium))) {
				return fmt.Errorf("`anonymous_pull_enabled` can only be applied when using the Standard/Premium Sku")
			}

			// data endpoint is only available for Premium Sku.
			if d.Get("data_endpoint_enabled").(bool) && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return fmt.Errorf("`data_endpoint_enabled` can only be applied when using the Premium Sku")
			}

			return nil
		}),
	}
}

func resourceContainerRegistryCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Registries
	operationClient := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Operation
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for  Container Registry creation.")

	id := registries.NewRegistryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry", id.ID())
		}
	}

	sId := commonids.NewSubscriptionID(subscriptionId)
	availabilityRequest := operation.RegistryNameCheckRequest{
		Name: id.RegistryName,
		Type: "Microsoft.ContainerRegistry/registries",
	}
	resp, err := operationClient.RegistriesCheckNameAvailability(ctx, sId, availabilityRequest)
	if err != nil {
		return fmt.Errorf("checking if the name %q was available: %+v", id.RegistryName, err)
	}

	if resp.Model == nil && resp.Model.NameAvailable == nil {
		return fmt.Errorf("checking name availability for %s: model was nil", id)
	}

	if available := *resp.Model.NameAvailable; !available {
		return fmt.Errorf("the name %q used for the Container Registry needs to be globally unique and isn't available: %s", id.RegistryName, *resp.Model.Message)
	}

	sku := d.Get("sku").(string)

	networkRuleSet := expandNetworkRuleSet(d.Get("network_rule_set").([]interface{}))
	if networkRuleSet != nil && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
		return fmt.Errorf("`network_rule_set_set` can only be specified for a Premium Sku. If you are reverting from a Premium to Basic SKU plese set network_rule_set = []")
	}

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	publicNetworkAccess := registries.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = registries.PublicNetworkAccessDisabled
	}

	zoneRedundancy := registries.ZoneRedundancyDisabled
	if d.Get("zone_redundancy_enabled").(bool) {
		zoneRedundancy = registries.ZoneRedundancyEnabled
	}

	parameters := registries.Registry{
		Location: location.Normalize(d.Get("location").(string)),
		Sku: registries.Sku{
			Name: registries.SkuName(sku),
			Tier: pointer.To(registries.SkuTier(sku)),
		},
		Identity: identity,
		Properties: &registries.RegistryProperties{
			AdminUserEnabled: utils.Bool(d.Get("admin_enabled").(bool)),
			Encryption:       expandEncryption(d.Get("encryption").([]interface{})),
			NetworkRuleSet:   networkRuleSet,
			Policies: &registries.Policies{
				QuarantinePolicy: expandQuarantinePolicy(d.Get("quarantine_policy_enabled").(bool)),
				RetentionPolicy:  expandRetentionPolicy(d.Get("retention_policy").([]interface{})),
				TrustPolicy:      expandTrustPolicy(d.Get("trust_policy").([]interface{})),
				ExportPolicy:     expandExportPolicy(d.Get("export_policy_enabled").(bool)),
			},
			PublicNetworkAccess:      &publicNetworkAccess,
			ZoneRedundancy:           &zoneRedundancy,
			AnonymousPullEnabled:     utils.Bool(d.Get("anonymous_pull_enabled").(bool)),
			DataEndpointEnabled:      utils.Bool(d.Get("data_endpoint_enabled").(bool)),
			NetworkRuleBypassOptions: pointer.To(registries.NetworkRuleBypassOptions(d.Get("network_rule_bypass_option").(string))),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// the ACR is being created so no previous geo-replication locations
	var oldGeoReplicationLocations, newGeoReplicationLocations []replications.Replication
	newGeoReplicationLocations = expandReplications(d.Get("georeplications").([]interface{}))
	// geo replications have been specified
	if len(newGeoReplicationLocations) > 0 {
		err = applyGeoReplicationLocations(ctx, meta, id, oldGeoReplicationLocations, newGeoReplicationLocations)
		if err != nil {
			return fmt.Errorf("applying geo replications for %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceContainerRegistryRead(d, meta)
}

func resourceContainerRegistryUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Registries
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for  Container Registry update.")

	id, err := registries.ParseRegistryID(d.Id())
	if err != nil {
		return err
	}

	sku := d.Get("sku").(string)
	skuChange := d.HasChange("sku")
	isBasicSku := strings.EqualFold(sku, string(registries.SkuNameBasic))
	isPremiumSku := strings.EqualFold(sku, string(registries.SkuNamePremium))
	isStandardSku := strings.EqualFold(sku, string(registries.SkuNameStandard))

	oldReplicationsRaw, newReplicationsRaw := d.GetChange("georeplications")
	hasGeoReplicationsChanges := d.HasChange("georeplications")
	oldReplications := oldReplicationsRaw.([]interface{})
	newReplications := newReplicationsRaw.([]interface{})

	// handle upgrade to Premium SKU first
	if skuChange && isPremiumSku {
		if err := applyContainerRegistrySku(d, meta, sku, *id); err != nil {
			return fmt.Errorf("applying sku %q for %s: %+v", sku, id, err)
		}
	}

	networkRuleSet := expandNetworkRuleSet(d.Get("network_rule_set").([]interface{}))
	if networkRuleSet != nil && isBasicSku {
		return fmt.Errorf("`network_rule_set_set` can only be specified for a Premium Sku. If you are reverting from a Premium to Basic SKU plese set network_rule_set = []")
	}

	publicNetworkAccess := registries.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = registries.PublicNetworkAccessDisabled
	}

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := registries.RegistryUpdateParameters{
		Properties: &registries.RegistryPropertiesUpdateParameters{
			AdminUserEnabled: utils.Bool(d.Get("admin_enabled").(bool)),
			NetworkRuleSet:   networkRuleSet,
			Policies: &registries.Policies{
				QuarantinePolicy: expandQuarantinePolicy(d.Get("quarantine_policy_enabled").(bool)),
				RetentionPolicy:  expandRetentionPolicy(d.Get("retention_policy").([]interface{})),
				TrustPolicy:      expandTrustPolicy(d.Get("trust_policy").([]interface{})),
				ExportPolicy:     expandExportPolicy(d.Get("export_policy_enabled").(bool)),
			},
			PublicNetworkAccess:      &publicNetworkAccess,
			Encryption:               expandEncryption(d.Get("encryption").([]interface{})),
			AnonymousPullEnabled:     utils.Bool(d.Get("anonymous_pull_enabled").(bool)),
			DataEndpointEnabled:      utils.Bool(d.Get("data_endpoint_enabled").(bool)),
			NetworkRuleBypassOptions: pointer.To(registries.NetworkRuleBypassOptions(d.Get("network_rule_bypass_option").(string))),
		},
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	// geo replication is only supported by Premium Sku
	if len(newReplications) > 0 && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
		return fmt.Errorf("ACR geo-replication can only be applied when using the Premium Sku.")
	}

	if hasGeoReplicationsChanges {
		err := applyGeoReplicationLocations(ctx, meta, *id, expandReplications(oldReplications), expandReplications(newReplications))
		if err != nil {
			return fmt.Errorf("applying geo replications for %s: %+v", id, err)
		}
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// downgrade to Basic or Standard SKU
	if skuChange && (isBasicSku || isStandardSku) {
		if err := applyContainerRegistrySku(d, meta, sku, *id); err != nil {
			return fmt.Errorf("applying sku %q for %s: %+v", sku, id, err)
		}
	}

	d.SetId(id.ID())

	return resourceContainerRegistryRead(d, meta)
}

func applyContainerRegistrySku(d *pluginsdk.ResourceData, meta interface{}, sku string, id registries.RegistryId) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Registries
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parameters := registries.RegistryUpdateParameters{
		Sku: &registries.Sku{
			Name: registries.SkuName(sku),
			Tier: pointer.To(registries.SkuTier(sku)),
		},
	}

	if err := client.UpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return nil
}

func applyGeoReplicationLocations(ctx context.Context, meta interface{}, registryId registries.RegistryId, oldGeoReplications []replications.Replication, newGeoReplications []replications.Replication) error {
	replicationClient := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Replications
	log.Printf("[INFO] preparing to apply geo-replications for Container Registry.")

	oldReplications := map[string]replications.Replication{}
	for _, replication := range oldGeoReplications {
		loc := location.Normalize(replication.Location)
		oldReplications[loc] = replication
	}

	newReplications := map[string]replications.Replication{}
	for _, replication := range newGeoReplications {
		loc := location.Normalize(replication.Location)
		newReplications[loc] = replication
	}

	// Delete replications that only appear in the old locations.
	for loc := range oldReplications {
		if _, ok := newReplications[loc]; ok {
			continue
		}
		id := replications.NewReplicationID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, loc)
		if err := replicationClient.DeleteThenPoll(ctx, id); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	// Create replications that only exists in the new locations.
	for loc, repl := range newReplications {
		if _, ok := oldReplications[loc]; ok {
			continue
		}
		id := replications.NewReplicationID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, loc)
		if err := replicationClient.CreateThenPoll(ctx, id, repl); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
	}

	// Update (potentially replace) replications that exists at both side.
	for loc, newRepl := range newReplications {
		oldRepl, ok := oldReplications[loc]
		if !ok {
			continue
		}
		// Compare old and new replication parameters to see whether it has updated.
		// If there no update, then skip it. Otherwise, need to check whether the update
		// can happen in place, or need a recreation.

		var (
			needUpdate  bool
			needReplace bool
		)
		// Since the replications here are all derived from expand function, where we guaranteed
		// each properties are non-nil. Whilst we are still doing nil check here in case.
		if oprop, nprop := oldRepl.Properties, newRepl.Properties; oprop != nil && nprop != nil {
			// zoneRedundency can't be updated in place
			if ov, nv := oprop.ZoneRedundancy, nprop.ZoneRedundancy; ov != nil && nv != nil && *ov != *nv {
				needUpdate = true
				needReplace = true
			}
			if ov, nv := oprop.RegionEndpointEnabled, nprop.RegionEndpointEnabled; ov != nil && nv != nil && *ov != *nv {
				needUpdate = true
			}
		}
		otag, ntag := *oldRepl.Tags, *newRepl.Tags
		if len(otag) != len(ntag) {
			needUpdate = true
		} else {
			for k, ov := range otag {
				nv, ok := ntag[k]
				if !ok {
					needUpdate = true
					break
				}
				if ov != nv {
					needUpdate = true
					break
				}
			}
		}

		if !needUpdate {
			continue
		}

		if needReplace {
			id := replications.NewReplicationID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, loc)
			if err := replicationClient.DeleteThenPoll(ctx, id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			// Following can be removed once https://github.com/Azure/azure-rest-api-specs/issues/18934 is resolved. Otherwise, the create right after delete will always fail.
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("context is missing a timeout")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"InProgress"},
				Target:  []string{"NotFound"},
				Refresh: func() (interface{}, string, error) {
					resp, err := replicationClient.Get(ctx, id)
					if err != nil {
						if response.WasNotFound(resp.HttpResponse) {
							return resp, "NotFound", nil
						}

						return nil, "Error", err
					}

					return resp, "InProgress", nil
				},
				ContinuousTargetOccurence: 5,
				PollInterval:              5 * time.Second,
				Timeout:                   time.Until(deadline),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("additional waiting for deletion of %s: %+v", id, err)
			}
		}

		id := replications.NewReplicationID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, loc)
		if err := replicationClient.CreateThenPoll(ctx, id, newRepl); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", id, err)
		}
	}

	return nil
}

func resourceContainerRegistryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Registries
	replicationClient := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Replications
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := registries.ParseRegistryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.RegistryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	loc := ""

	if model := resp.Model; model != nil {
		loc = location.Normalize(model.Location)
		d.Set("location", loc)

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("sku", string(pointer.From(model.Sku.Tier)))

		if props := model.Properties; props != nil {
			d.Set("admin_enabled", props.AdminUserEnabled)
			d.Set("login_server", props.LoginServer)
			d.Set("public_network_access_enabled", *props.PublicNetworkAccess == registries.PublicNetworkAccessEnabled)

			networkRuleSet := flattenNetworkRuleSet(props.NetworkRuleSet)
			if err := d.Set("network_rule_set", networkRuleSet); err != nil {
				return fmt.Errorf("setting `network_rule_set`: %+v", err)
			}

			d.Set("quarantine_policy_enabled", flattenQuarantinePolicy(props.Policies))
			if err := d.Set("retention_policy", flattenRetentionPolicy(props.Policies)); err != nil {
				return fmt.Errorf("setting `retention_policy`: %+v", err)
			}
			if err := d.Set("trust_policy", flattenTrustPolicy(props.Policies)); err != nil {
				return fmt.Errorf("setting `trust_policy`: %+v", err)
			}
			d.Set("export_policy_enabled", flattenExportPolicy(props.Policies))
			if err := d.Set("encryption", flattenEncryption(props.Encryption)); err != nil {
				return fmt.Errorf("setting `encryption`: %+v", err)
			}
			d.Set("zone_redundancy_enabled", *props.ZoneRedundancy == registries.ZoneRedundancyEnabled)
			d.Set("anonymous_pull_enabled", props.AnonymousPullEnabled)
			d.Set("data_endpoint_enabled", props.DataEndpointEnabled)
			d.Set("network_rule_bypass_option", string(pointer.From(props.NetworkRuleBypassOptions)))

			if *props.AdminUserEnabled {
				credsResp, errList := client.ListCredentials(ctx, *id)
				if errList != nil {
					return fmt.Errorf("retrieving credentials for %s: %s", *id, errList)
				}

				if credsModel := credsResp.Model; credsModel != nil {
					d.Set("admin_username", credsModel.Username)
					for _, v := range *credsModel.Passwords {
						d.Set("admin_password", v.Value)
						break
					}
				}
			} else {
				d.Set("admin_username", "")
				d.Set("admin_password", "")
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
	}

	rId, err := replications.ParseRegistryID(id.ID())
	if err != nil {
		return err
	}
	replicationsResp, err := replicationClient.List(ctx, *rId)
	if err != nil {
		return fmt.Errorf("retrieving replications for %s: %s", *id, err)
	}

	geoReplications := make([]interface{}, 0)
	if replicationsModel := replicationsResp.Model; replicationsModel != nil {
		for _, value := range *replicationsModel {
			valueLocation := location.Normalize(value.Location)
			if valueLocation != loc {
				replication := make(map[string]interface{})
				replication["location"] = valueLocation
				replication["tags"] = tags.Flatten(value.Tags)
				replication["zone_redundancy_enabled"] = *value.Properties.ZoneRedundancy == replications.ZoneRedundancyEnabled
				replication["regional_endpoint_enabled"] = value.Properties.RegionEndpointEnabled != nil && *value.Properties.RegionEndpointEnabled
				geoReplications = append(geoReplications, replication)
			}
		}
	}

	// The order of the georeplications returned from the list API is not consistent. We simply order it alphabetically to be consistent.
	sort.Slice(geoReplications, func(i, j int) bool {
		return geoReplications[i].(map[string]interface{})["location"].(string) < geoReplications[j].(map[string]interface{})["location"].(string)
	})

	d.Set("georeplications", geoReplications)

	return nil
}

func resourceContainerRegistryDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Registries
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := registries.ParseRegistryID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandNetworkRuleSet(profiles []interface{}) *registries.NetworkRuleSet {
	if len(profiles) == 0 {
		return nil
	}

	profile := profiles[0].(map[string]interface{})

	ipRuleConfigs := profile["ip_rule"].(*pluginsdk.Set).List()
	ipRules := make([]registries.IPRule, 0)
	for _, ipRuleInterface := range ipRuleConfigs {
		config := ipRuleInterface.(map[string]interface{})
		newIpRule := registries.IPRule{
			Action: pointer.To(registries.Action(config["action"].(string))),
			Value:  config["ip_range"].(string),
		}
		ipRules = append(ipRules, newIpRule)
	}

	networkRuleConfigs := profile["virtual_network"].(*pluginsdk.Set).List()
	virtualNetworkRules := make([]registries.VirtualNetworkRule, 0)
	for _, networkRuleInterface := range networkRuleConfigs {
		config := networkRuleInterface.(map[string]interface{})
		newVirtualNetworkRule := registries.VirtualNetworkRule{
			Action: pointer.To(registries.Action(config["action"].(string))),
			Id:     config["subnet_id"].(string),
		}
		virtualNetworkRules = append(virtualNetworkRules, newVirtualNetworkRule)
	}

	networkRuleSet := registries.NetworkRuleSet{
		DefaultAction:       registries.DefaultAction(profile["default_action"].(string)),
		IPRules:             &ipRules,
		VirtualNetworkRules: &virtualNetworkRules,
	}
	return &networkRuleSet
}

func expandQuarantinePolicy(enabled bool) *registries.QuarantinePolicy {
	quarantinePolicy := registries.QuarantinePolicy{
		Status: pointer.To(registries.PolicyStatusDisabled),
	}

	if enabled {
		quarantinePolicy.Status = pointer.To(registries.PolicyStatusEnabled)
	}

	return &quarantinePolicy
}

func expandRetentionPolicy(p []interface{}) *registries.RetentionPolicy {
	retentionPolicy := registries.RetentionPolicy{
		Status: pointer.To(registries.PolicyStatusDisabled),
	}

	if len(p) > 0 {
		v := p[0].(map[string]interface{})
		days := int32(v["days"].(int))
		enabled := v["enabled"].(bool)
		if enabled {
			retentionPolicy.Status = pointer.To(registries.PolicyStatusEnabled)
		}
		retentionPolicy.Days = utils.Int64(int64(days))
	}

	return &retentionPolicy
}

func expandTrustPolicy(p []interface{}) *registries.TrustPolicy {
	trustPolicy := registries.TrustPolicy{
		Status: pointer.To(registries.PolicyStatusDisabled),
	}

	if len(p) > 0 {
		v := p[0].(map[string]interface{})
		enabled := v["enabled"].(bool)
		if enabled {
			trustPolicy.Status = pointer.To(registries.PolicyStatusEnabled)
		}
		trustPolicy.Type = pointer.To(registries.TrustPolicyTypeNotary)
	}

	return &trustPolicy
}

func expandExportPolicy(enabled bool) *registries.ExportPolicy {
	exportPolicy := registries.ExportPolicy{
		Status: pointer.To(registries.ExportPolicyStatusDisabled),
	}

	if enabled {
		exportPolicy.Status = pointer.To(registries.ExportPolicyStatusEnabled)
	}

	return &exportPolicy
}

func expandReplications(p []interface{}) []replications.Replication {
	reps := make([]replications.Replication, 0)
	if p == nil {
		return reps
	}
	for _, v := range p {
		value := v.(map[string]interface{})
		location := azure.NormalizeLocation(value["location"])
		tags := tags.Expand(value["tags"].(map[string]interface{}))
		zoneRedundancy := replications.ZoneRedundancyDisabled
		if value["zone_redundancy_enabled"].(bool) {
			zoneRedundancy = replications.ZoneRedundancyEnabled
		}
		reps = append(reps, replications.Replication{
			Location: location,
			Name:     &location,
			Tags:     tags,
			Properties: &replications.ReplicationProperties{
				ZoneRedundancy:        &zoneRedundancy,
				RegionEndpointEnabled: utils.Bool(value["regional_endpoint_enabled"].(bool)),
			},
		})
	}
	return reps
}

func expandEncryption(e []interface{}) *registries.EncryptionProperty {
	encryptionProperty := registries.EncryptionProperty{
		Status: pointer.To(registries.EncryptionStatusDisabled),
	}
	if len(e) > 0 {
		v := e[0].(map[string]interface{})
		enabled := v["enabled"].(bool)
		if enabled {
			encryptionProperty.Status = pointer.To(registries.EncryptionStatusEnabled)
			keyId := v["key_vault_key_id"].(string)
			identityClientId := v["identity_client_id"].(string)
			encryptionProperty.KeyVaultProperties = &registries.KeyVaultProperties{
				KeyIdentifier: &keyId,
				Identity:      &identityClientId,
			}
		}
	}

	return &encryptionProperty
}

func flattenEncryption(encryptionProperty *registries.EncryptionProperty) []interface{} {
	if encryptionProperty == nil {
		return nil
	}
	encryption := make(map[string]interface{})
	encryption["enabled"] = strings.EqualFold(string(*encryptionProperty.Status), string(registries.EncryptionStatusEnabled))
	if encryptionProperty.KeyVaultProperties != nil {
		encryption["key_vault_key_id"] = encryptionProperty.KeyVaultProperties.KeyIdentifier
		encryption["identity_client_id"] = encryptionProperty.KeyVaultProperties.Identity
	}

	return []interface{}{encryption}
}

func flattenNetworkRuleSet(networkRuleSet *registries.NetworkRuleSet) []interface{} {
	if networkRuleSet == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	values["default_action"] = string(networkRuleSet.DefaultAction)

	ipRules := make([]interface{}, 0)
	for _, ipRule := range *networkRuleSet.IPRules {
		value := make(map[string]interface{})
		value["action"] = string(*ipRule.Action)

		// When a /32 CIDR is passed as an ip rule, Azure will drop the /32 leading to the resource wanting to be re-created next run
		if !strings.Contains(ipRule.Value, "/") {
			ipRule.Value += "/32"
		}

		value["ip_range"] = ipRule.Value
		ipRules = append(ipRules, value)
	}

	values["ip_rule"] = ipRules

	virtualNetworkRules := make([]interface{}, 0)

	if networkRuleSet.VirtualNetworkRules != nil {
		for _, virtualNetworkRule := range *networkRuleSet.VirtualNetworkRules {
			value := make(map[string]interface{})
			value["action"] = string(*virtualNetworkRule.Action)

			value["subnet_id"] = virtualNetworkRule.Id
			virtualNetworkRules = append(virtualNetworkRules, value)
		}
	}

	values["virtual_network"] = virtualNetworkRules

	return []interface{}{values}
}

func flattenQuarantinePolicy(p *registries.Policies) bool {
	if p == nil || p.QuarantinePolicy == nil {
		return false
	}

	return *p.QuarantinePolicy.Status == registries.PolicyStatusEnabled
}

func flattenRetentionPolicy(p *registries.Policies) []interface{} {
	if p == nil || p.RetentionPolicy == nil {
		return nil
	}

	r := *p.RetentionPolicy
	retentionPolicy := make(map[string]interface{})
	retentionPolicy["days"] = r.Days
	enabled := strings.EqualFold(string(*r.Status), string(registries.PolicyStatusEnabled))
	retentionPolicy["enabled"] = utils.Bool(enabled)
	return []interface{}{retentionPolicy}
}

func flattenTrustPolicy(p *registries.Policies) []interface{} {
	if p == nil || p.TrustPolicy == nil {
		return nil
	}

	t := *p.TrustPolicy
	trustPolicy := make(map[string]interface{})
	enabled := strings.EqualFold(string(*t.Status), string(registries.PolicyStatusEnabled))
	trustPolicy["enabled"] = utils.Bool(enabled)
	return []interface{}{trustPolicy}
}

func flattenExportPolicy(p *registries.Policies) bool {
	if p == nil || p.ExportPolicy == nil {
		return false
	}

	return *p.ExportPolicy.Status == registries.ExportPolicyStatusEnabled
}

func resourceContainerRegistrySchema() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerValidate.ContainerRegistryName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(registries.SkuNameBasic),
				string(registries.SkuNameStandard),
				string(registries.SkuNamePremium),
			}, false),
		},

		"admin_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"georeplications": {
			// Don't make this a TypeSet since TypeSet has bugs when there is a nested property using `StateFunc`.
			// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
			Type:       pluginsdk.TypeList,
			Optional:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAuto,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": commonschema.LocationWithoutForceNew(),

					"zone_redundancy_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"regional_endpoint_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"tags": commonschema.Tags(),
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"login_server": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_password": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"encryption": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true,
			MaxItems:   1,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"identity_client_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
					},
				},
			},
		},

		"network_rule_set": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true,
			MaxItems:   1,
			ConfigMode: pluginsdk.SchemaConfigModeAttr, // make sure we can set this to an empty array for Premium -> Basic
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  registries.DefaultActionAllow,
						ValidateFunc: validation.StringInSlice([]string{
							string(registries.DefaultActionAllow),
							string(registries.DefaultActionDeny),
						}, false),
					},

					"ip_rule": {
						Type:       pluginsdk.TypeSet,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(registries.ActionAllow),
									}, false),
								},
								"ip_range": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.CIDR,
								},
							},
						},
					},

					"virtual_network": {
						Deprecated: " This is only used exclusively for service endpoints (which is a feature being deprecated). Users are expected to use Private Endpoints instead",
						Type:       pluginsdk.TypeSet,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(registries.ActionAllow),
									}, false),
								},
								"subnet_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateSubnetID,
								},
							},
						},
					},
				},
			},
		},

		"quarantine_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"retention_policy": {
			Type:       pluginsdk.TypeList,
			MaxItems:   1,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"days": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  7,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"trust_policy": {
			Type:       pluginsdk.TypeList,
			MaxItems:   1,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"export_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"zone_redundancy_enabled": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
			Default:  false,
		},

		"anonymous_pull_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"data_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"network_rule_bypass_option": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(registries.NetworkRuleBypassOptionsAzureServices),
				string(registries.NetworkRuleBypassOptionsNone),
			}, false),
			Default: string(registries.NetworkRuleBypassOptionsAzureServices),
		},

		"tags": commonschema.Tags(),
	}

	if features.FourPointOhBeta() {
		delete(schema["network_rule_set"].Elem.(*pluginsdk.Resource).Schema, "virtual_network")
	}

	return schema
}
