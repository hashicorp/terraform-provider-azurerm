package containers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2021-08-01-preview/containerregistry"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
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
			_, err := parse.RegistryID(id)
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

			hasGeoReplications := false
			if !features.ThreePointOhBeta() {
				geoReplicationLocations := d.Get("georeplication_locations").(*pluginsdk.Set)
				hasGeoReplications = geoReplicationLocations.Len() > 0
			}
			geoReplications := d.Get("georeplications").([]interface{})
			hasGeoReplicationsApplied := hasGeoReplications || len(geoReplications) > 0
			// if locations have been specified for geo-replication then, the SKU has to be Premium
			if hasGeoReplicationsApplied && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
				return fmt.Errorf("ACR geo-replication can only be applied when using the Premium Sku.")
			}

			quarantinePolicyEnabled := d.Get("quarantine_policy_enabled").(bool)
			if quarantinePolicyEnabled && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
				return fmt.Errorf("ACR quarantine policy can only be applied when using the Premium Sku. If you are downgrading from a Premium SKU please set quarantine_policy {}")
			}

			retentionPolicyEnabled, ok := d.GetOk("retention_policy.0.enabled")
			if ok && retentionPolicyEnabled.(bool) && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
				return fmt.Errorf("ACR retention policy can only be applied when using the Premium Sku. If you are downgrading from a Premium SKU please set retention_policy {}")
			}

			trustPolicyEnabled, ok := d.GetOk("trust_policy.0.enabled")
			if ok && trustPolicyEnabled.(bool) && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
				return fmt.Errorf("ACR trust policy can only be applied when using the Premium Sku. If you are downgrading from a Premium SKU please set trust_policy {}")
			}

			encryptionEnabled, ok := d.GetOk("encryption.0.enabled")
			if ok && encryptionEnabled.(bool) && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
				return fmt.Errorf("ACR encryption can only be applied when using the Premium Sku.")
			}

			// zone redundancy is only available for Premium Sku.
			zoneRedundancyEnabled, ok := d.GetOk("zone_redundancy_enabled")
			if ok && zoneRedundancyEnabled.(bool) && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
				return fmt.Errorf("ACR zone redundancy can only be applied when using the Premium Sku")
			}
			for _, loc := range geoReplications {
				loc := loc.(map[string]interface{})
				zoneRedundancyEnabled, ok := loc["zone_redundancy_enabled"]
				if ok && zoneRedundancyEnabled.(bool) && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
					return fmt.Errorf("ACR zone redundancy can only be applied when using the Premium Sku")
				}
			}

			// anonymous pull is only available for Standard/Premium Sku.
			if d.Get("anonymous_pull_enabled").(bool) && (!strings.EqualFold(sku, string(containerregistry.SkuNameStandard)) && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium))) {
				return fmt.Errorf("`anonymous_pull_enabled` can only be applied when using the Standard/Premium Sku")
			}

			// data endpoint is only available for Premium Sku.
			if d.Get("data_endpoint_enabled").(bool) && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
				return fmt.Errorf("`data_endpoint_enabled` can only be applied when using the Premium Sku")
			}

			return nil
		}),
	}
}

func resourceContainerRegistryCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.RegistriesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for  Container Registry creation.")

	id := parse.NewRegistryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_container_registry", id.ID())
		}
	}

	availabilityRequest := containerregistry.RegistryNameCheckRequest{
		Name: utils.String(id.Name),
		Type: utils.String("Microsoft.ContainerRegistry/registries"),
	}
	available, err := client.CheckNameAvailability(ctx, availabilityRequest)
	if err != nil {
		return fmt.Errorf("checking if the name %q was available: %+v", id.Name, err)
	}

	if !*available.NameAvailable {
		return fmt.Errorf("the name %q used for the Container Registry needs to be globally unique and isn't available: %s", id.Name, *available.Message)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	adminUserEnabled := d.Get("admin_enabled").(bool)
	t := d.Get("tags").(map[string]interface{})
	geoReplications := d.Get("georeplications").([]interface{})

	networkRuleSet := expandNetworkRuleSet(d.Get("network_rule_set").([]interface{}))
	if networkRuleSet != nil && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
		return fmt.Errorf("`network_rule_set_set` can only be specified for a Premium Sku. If you are reverting from a Premium to Basic SKU plese set network_rule_set = []")
	}

	quarantinePolicy := expandQuarantinePolicy(d.Get("quarantine_policy_enabled").(bool))

	retentionPolicyRaw := d.Get("retention_policy").([]interface{})
	retentionPolicy := expandRetentionPolicy(retentionPolicyRaw)

	trustPolicyRaw := d.Get("trust_policy").([]interface{})
	trustPolicy := expandTrustPolicy(trustPolicyRaw)

	encryptionRaw := d.Get("encryption").([]interface{})
	encryption := expandEncryption(encryptionRaw)

	identity, err := expandRegistryIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	publicNetworkAccess := containerregistry.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = containerregistry.PublicNetworkAccessDisabled
	}

	zoneRedundancy := containerregistry.ZoneRedundancyDisabled
	if d.Get("zone_redundancy_enabled").(bool) {
		zoneRedundancy = containerregistry.ZoneRedundancyEnabled
	}

	parameters := containerregistry.Registry{
		Location: &location,
		Sku: &containerregistry.Sku{
			Name: containerregistry.SkuName(sku),
			Tier: containerregistry.SkuTier(sku),
		},
		Identity: identity,
		RegistryProperties: &containerregistry.RegistryProperties{
			AdminUserEnabled: utils.Bool(adminUserEnabled),
			Encryption:       encryption,
			NetworkRuleSet:   networkRuleSet,
			Policies: &containerregistry.Policies{
				QuarantinePolicy: quarantinePolicy,
				RetentionPolicy:  retentionPolicy,
				TrustPolicy:      trustPolicy,
			},
			PublicNetworkAccess:      publicNetworkAccess,
			ZoneRedundancy:           zoneRedundancy,
			AnonymousPullEnabled:     utils.Bool(d.Get("anonymous_pull_enabled").(bool)),
			DataEndpointEnabled:      utils.Bool(d.Get("data_endpoint_enabled").(bool)),
			NetworkRuleBypassOptions: containerregistry.NetworkRuleBypassOptions(d.Get("network_rule_bypass_option").(string)),
		},

		Tags: tags.Expand(t),
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	// the ACR is being created so no previous geo-replication locations
	var oldGeoReplicationLocations, newGeoReplicationLocations []containerregistry.Replication
	newGeoReplicationLocations = expandReplications(geoReplications)
	if !features.ThreePointOhBeta() {
		geoReplicationLocations := d.Get("georeplication_locations").(*pluginsdk.Set)
		if geoReplicationLocations != nil && geoReplicationLocations.Len() > 0 {
			newGeoReplicationLocations = expandReplicationsFromLocations(geoReplicationLocations.List())
		}
	}
	// geo replications have been specified
	if len(newGeoReplicationLocations) > 0 {
		err = applyGeoReplicationLocations(d, meta, id.ResourceGroup, id.Name, oldGeoReplicationLocations, newGeoReplicationLocations)
		if err != nil {
			return fmt.Errorf("applying geo replications for %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceContainerRegistryRead(d, meta)
}

func resourceContainerRegistryUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.RegistriesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for  Container Registry update.")

	id, err := parse.RegistryID(d.Id())
	if err != nil {
		return err
	}

	sku := d.Get("sku").(string)
	skuChange := d.HasChange("sku")
	isBasicSku := strings.EqualFold(sku, string(containerregistry.SkuNameBasic))
	isPremiumSku := strings.EqualFold(sku, string(containerregistry.SkuNamePremium))
	isStandardSku := strings.EqualFold(sku, string(containerregistry.SkuNameStandard))

	adminUserEnabled := d.Get("admin_enabled").(bool)
	t := d.Get("tags").(map[string]interface{})

	oldReplicationsRaw, newReplicationsRaw := d.GetChange("georeplications")
	hasGeoReplicationsChanges := d.HasChange("georeplications")
	oldReplications := oldReplicationsRaw.([]interface{})
	newReplications := newReplicationsRaw.([]interface{})

	// handle upgrade to Premium SKU first
	if skuChange && isPremiumSku {
		if err := applyContainerRegistrySku(d, meta, sku, id.ResourceGroup, id.Name); err != nil {
			return fmt.Errorf("applying sku %q for %s: %+v", sku, id, err)
		}
	}

	networkRuleSet := expandNetworkRuleSet(d.Get("network_rule_set").([]interface{}))
	if networkRuleSet != nil && isBasicSku {
		return fmt.Errorf("`network_rule_set_set` can only be specified for a Premium Sku. If you are reverting from a Premium to Basic SKU plese set network_rule_set = []")
	}

	quarantinePolicy := expandQuarantinePolicy(d.Get("quarantine_policy_enabled").(bool))
	retentionPolicy := expandRetentionPolicy(d.Get("retention_policy").([]interface{}))
	trustPolicy := expandTrustPolicy(d.Get("trust_policy").([]interface{}))

	publicNetworkAccess := containerregistry.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = containerregistry.PublicNetworkAccessDisabled
	}

	identity, err := expandRegistryIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	encryptionRaw := d.Get("encryption").([]interface{})
	encryption := expandEncryption(encryptionRaw)

	parameters := containerregistry.RegistryUpdateParameters{
		RegistryPropertiesUpdateParameters: &containerregistry.RegistryPropertiesUpdateParameters{
			AdminUserEnabled: utils.Bool(adminUserEnabled),
			NetworkRuleSet:   networkRuleSet,
			Policies: &containerregistry.Policies{
				QuarantinePolicy: quarantinePolicy,
				RetentionPolicy:  retentionPolicy,
				TrustPolicy:      trustPolicy,
			},
			PublicNetworkAccess:      publicNetworkAccess,
			Encryption:               encryption,
			AnonymousPullEnabled:     utils.Bool(d.Get("anonymous_pull_enabled").(bool)),
			DataEndpointEnabled:      utils.Bool(d.Get("data_endpoint_enabled").(bool)),
			NetworkRuleBypassOptions: containerregistry.NetworkRuleBypassOptions(d.Get("network_rule_bypass_option").(string)),
		},
		Identity: identity,
		Tags:     tags.Expand(t),
	}

	var hasGeoReplicationLocationsChanges bool
	var hasNewGeoReplicationLocations bool
	oldGeoReplicationLocations := make([]interface{}, 0)
	newGeoReplicationLocations := make([]interface{}, 0)

	if !features.ThreePointOhBeta() {
		hasGeoReplicationLocationsChanges = d.HasChange("georeplication_locations")
		old, new := d.GetChange("georeplication_locations")
		oldGeoReplicationLocations = old.(*pluginsdk.Set).List()
		newGeoReplicationLocations := new.(*pluginsdk.Set).List()
		hasNewGeoReplicationLocations = len(newGeoReplicationLocations) > 0
	}

	// geo replication is only supported by Premium Sku
	hasGeoReplicationsApplied := hasNewGeoReplicationLocations || len(newReplications) > 0
	if hasGeoReplicationsApplied && !strings.EqualFold(sku, string(containerregistry.SkuNamePremium)) {
		return fmt.Errorf("ACR geo-replication can only be applied when using the Premium Sku.")
	}

	if hasGeoReplicationsChanges {
		err := applyGeoReplicationLocations(d, meta, id.ResourceGroup, id.Name, expandReplications(oldReplications), expandReplications(newReplications))
		if err != nil {
			return fmt.Errorf("applying geo replications for %s: %+v", id, err)
		}
	} else if hasGeoReplicationLocationsChanges {
		err := applyGeoReplicationLocations(d, meta, id.ResourceGroup, id.Name, expandReplicationsFromLocations(oldGeoReplicationLocations), expandReplicationsFromLocations(newGeoReplicationLocations))
		if err != nil {
			return fmt.Errorf("applying geo replications for %s: %+v", id, err)
		}
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	// downgrade to Basic or Standard SKU
	if skuChange && (isBasicSku || isStandardSku) {
		if err := applyContainerRegistrySku(d, meta, sku, id.ResourceGroup, id.Name); err != nil {
			return fmt.Errorf("applying sku %q for %s: %+v", sku, id, err)
		}
	}

	d.SetId(id.ID())

	return resourceContainerRegistryRead(d, meta)
}

func applyContainerRegistrySku(d *pluginsdk.ResourceData, meta interface{}, sku string, resourceGroup string, name string) error {
	client := meta.(*clients.Client).Containers.RegistriesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parameters := containerregistry.RegistryUpdateParameters{
		Sku: &containerregistry.Sku{
			Name: containerregistry.SkuName(sku),
			Tier: containerregistry.SkuTier(sku),
		},
	}

	future, err := client.Update(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("updating Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func applyGeoReplicationLocations(d *pluginsdk.ResourceData, meta interface{}, resourceGroup string, name string, oldGeoReplications []containerregistry.Replication, newGeoReplications []containerregistry.Replication) error {
	replicationClient := meta.(*clients.Client).Containers.ReplicationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing to apply geo-replications for  Container Registry.")

	// delete previously deployed locations
	for _, replication := range oldGeoReplications {
		if replication.Location == nil {
			continue
		}
		oldLocation := azure.NormalizeLocation(*replication.Location)

		future, err := replicationClient.Delete(ctx, resourceGroup, name, oldLocation)
		if err != nil {
			return fmt.Errorf("deleting Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, oldLocation, err)
		}
		if err = future.WaitForCompletionRef(ctx, replicationClient.Client); err != nil {
			return fmt.Errorf("waiting for deletion of Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, oldLocation, err)
		}
	}

	// create new geo-replication locations
	for _, replication := range newGeoReplications {
		if replication.Location == nil {
			continue
		}
		locationToCreate := azure.NormalizeLocation(*replication.Location)
		future, err := replicationClient.Create(ctx, resourceGroup, name, locationToCreate, replication)
		if err != nil {
			return fmt.Errorf("creating Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, locationToCreate, err)
		}

		if err = future.WaitForCompletionRef(ctx, replicationClient.Client); err != nil {
			return fmt.Errorf("waiting for creation of Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, locationToCreate, err)
		}
	}

	return nil
}

func resourceContainerRegistryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.RegistriesClient
	replicationClient := meta.(*clients.Client).Containers.ReplicationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RegistryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Registry %q was not found in Resource Group %q", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure Container Registry %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	location := resp.Location
	if location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("admin_enabled", resp.AdminUserEnabled)
	d.Set("login_server", resp.LoginServer)
	d.Set("public_network_access_enabled", resp.PublicNetworkAccess == containerregistry.PublicNetworkAccessEnabled)

	networkRuleSet := flattenNetworkRuleSet(resp.NetworkRuleSet)
	if err := d.Set("network_rule_set", networkRuleSet); err != nil {
		return fmt.Errorf("setting `network_rule_set`: %+v", err)
	}

	identity, _ := flattenRegistryIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if properties := resp.RegistryProperties; properties != nil {
		if err := d.Set("quarantine_policy_enabled", flattenQuarantinePolicy(properties.Policies)); err != nil {
			return fmt.Errorf("setting `quarantine_policy`: %+v", err)
		}
		if err := d.Set("retention_policy", flattenRetentionPolicy(properties.Policies)); err != nil {
			return fmt.Errorf("setting `retention_policy`: %+v", err)
		}
		if err := d.Set("trust_policy", flattenTrustPolicy(properties.Policies)); err != nil {
			return fmt.Errorf("setting `trust_policy`: %+v", err)
		}
		if err := d.Set("encryption", flattenEncryption(properties.Encryption)); err != nil {
			return fmt.Errorf("setting `encryption`: %+v", err)
		}
		d.Set("zone_redundancy_enabled", properties.ZoneRedundancy == containerregistry.ZoneRedundancyEnabled)
		d.Set("anonymous_pull_enabled", properties.AnonymousPullEnabled)
		d.Set("data_endpoint_enabled", properties.DataEndpointEnabled)
		d.Set("network_rule_bypass_option", string(properties.NetworkRuleBypassOptions))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Tier))
	}

	if *resp.AdminUserEnabled {
		credsResp, errList := client.ListCredentials(ctx, id.ResourceGroup, id.Name)
		if errList != nil {
			return fmt.Errorf("making Read request on Azure Container Registry %s for Credentials: %s", id.Name, errList)
		}

		d.Set("admin_username", credsResp.Username)
		for _, v := range *credsResp.Passwords {
			d.Set("admin_password", v.Value)
			break
		}
	} else {
		d.Set("admin_username", "")
		d.Set("admin_password", "")
	}

	replications, err := replicationClient.List(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("making Read request on Azure Container Registry %s for replications: %s", id.Name, err)
	}

	geoReplicationLocations := make([]interface{}, 0)
	geoReplications := make([]interface{}, 0)
	for _, value := range replications.Values() {
		if value.Location != nil {
			valueLocation := azure.NormalizeLocation(*value.Location)
			if location != nil && valueLocation != azure.NormalizeLocation(*location) {
				geoReplicationLocations = append(geoReplicationLocations, *value.Location)
				replication := make(map[string]interface{})
				replication["location"] = valueLocation
				replication["tags"] = tags.Flatten(value.Tags)
				replication["zone_redundancy_enabled"] = value.ZoneRedundancy == containerregistry.ZoneRedundancyEnabled
				replication["regional_endpoint_enabled"] = value.RegionEndpointEnabled != nil && *value.RegionEndpointEnabled
				geoReplications = append(geoReplications, replication)
			}
		}
	}

	if !features.ThreePointOhBeta() {
		d.Set("georeplication_locations", geoReplicationLocations)
		d.Set("storage_account_id", "")
	}

	d.Set("georeplications", geoReplications)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceContainerRegistryDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.RegistriesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RegistryID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Container Registry %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Container Registry %q (Resource Group %q) to be deleted: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandNetworkRuleSet(profiles []interface{}) *containerregistry.NetworkRuleSet {
	if len(profiles) == 0 {
		return nil
	}

	profile := profiles[0].(map[string]interface{})

	ipRuleConfigs := profile["ip_rule"].(*pluginsdk.Set).List()
	ipRules := make([]containerregistry.IPRule, 0)
	for _, ipRuleInterface := range ipRuleConfigs {
		config := ipRuleInterface.(map[string]interface{})
		newIpRule := containerregistry.IPRule{
			Action:           containerregistry.Action(config["action"].(string)),
			IPAddressOrRange: utils.String(config["ip_range"].(string)),
		}
		ipRules = append(ipRules, newIpRule)
	}

	networkRuleConfigs := profile["virtual_network"].(*pluginsdk.Set).List()
	virtualNetworkRules := make([]containerregistry.VirtualNetworkRule, 0)
	for _, networkRuleInterface := range networkRuleConfigs {
		config := networkRuleInterface.(map[string]interface{})
		newVirtualNetworkRule := containerregistry.VirtualNetworkRule{
			Action:                   containerregistry.Action(config["action"].(string)),
			VirtualNetworkResourceID: utils.String(config["subnet_id"].(string)),
		}
		virtualNetworkRules = append(virtualNetworkRules, newVirtualNetworkRule)
	}

	networkRuleSet := containerregistry.NetworkRuleSet{
		DefaultAction:       containerregistry.DefaultAction(profile["default_action"].(string)),
		IPRules:             &ipRules,
		VirtualNetworkRules: &virtualNetworkRules,
	}
	return &networkRuleSet
}

func expandQuarantinePolicy(enabled bool) *containerregistry.QuarantinePolicy {
	quarantinePolicy := containerregistry.QuarantinePolicy{
		Status: containerregistry.PolicyStatusDisabled,
	}

	if enabled {
		quarantinePolicy.Status = containerregistry.PolicyStatusEnabled
	}

	return &quarantinePolicy
}

func expandRetentionPolicy(p []interface{}) *containerregistry.RetentionPolicy {
	retentionPolicy := containerregistry.RetentionPolicy{
		Status: containerregistry.PolicyStatusDisabled,
	}

	if len(p) > 0 {
		v := p[0].(map[string]interface{})
		days := int32(v["days"].(int))
		enabled := v["enabled"].(bool)
		if enabled {
			retentionPolicy.Status = containerregistry.PolicyStatusEnabled
		}
		retentionPolicy.Days = utils.Int32(days)
	}

	return &retentionPolicy
}

func expandTrustPolicy(p []interface{}) *containerregistry.TrustPolicy {
	trustPolicy := containerregistry.TrustPolicy{
		Status: containerregistry.PolicyStatusDisabled,
	}

	if len(p) > 0 {
		v := p[0].(map[string]interface{})
		enabled := v["enabled"].(bool)
		if enabled {
			trustPolicy.Status = containerregistry.PolicyStatusEnabled
		}
		trustPolicy.Type = containerregistry.TrustPolicyTypeNotary
	}

	return &trustPolicy
}

func expandReplicationsFromLocations(p []interface{}) []containerregistry.Replication {
	replications := make([]containerregistry.Replication, 0)
	for _, value := range p {
		location := azure.NormalizeLocation(value)
		replications = append(replications, containerregistry.Replication{
			Location: &location,
			Name:     &location,
		})
	}
	return replications
}

func expandReplications(p []interface{}) []containerregistry.Replication {
	replications := make([]containerregistry.Replication, 0)
	if p == nil {
		return replications
	}
	for _, v := range p {
		value := v.(map[string]interface{})
		location := azure.NormalizeLocation(value["location"])
		tags := tags.Expand(value["tags"].(map[string]interface{}))
		zoneRedundancy := containerregistry.ZoneRedundancyDisabled
		if value["zone_redundancy_enabled"].(bool) {
			zoneRedundancy = containerregistry.ZoneRedundancyEnabled
		}
		replications = append(replications, containerregistry.Replication{
			Location: &location,
			Name:     &location,
			Tags:     tags,
			ReplicationProperties: &containerregistry.ReplicationProperties{
				ZoneRedundancy:        zoneRedundancy,
				RegionEndpointEnabled: utils.Bool(value["regional_endpoint_enabled"].(bool)),
			},
		})
	}
	return replications
}

func expandRegistryIdentity(input []interface{}) (*containerregistry.IdentityProperties, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := containerregistry.IdentityProperties{
		Type: containerregistry.ResourceIdentityType(string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*containerregistry.UserIdentityProperties)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &containerregistry.UserIdentityProperties{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func expandEncryption(e []interface{}) *containerregistry.EncryptionProperty {
	encryptionProperty := containerregistry.EncryptionProperty{
		Status: containerregistry.EncryptionStatusDisabled,
	}
	if len(e) > 0 {
		v := e[0].(map[string]interface{})
		enabled := v["enabled"].(bool)
		if enabled {
			encryptionProperty.Status = containerregistry.EncryptionStatusEnabled
			keyId := v["key_vault_key_id"].(string)
			identityClientId := v["identity_client_id"].(string)
			encryptionProperty.KeyVaultProperties = &containerregistry.KeyVaultProperties{
				KeyIdentifier: &keyId,
				Identity:      &identityClientId,
			}
		}
	}

	return &encryptionProperty
}

func flattenEncryption(encryptionProperty *containerregistry.EncryptionProperty) []interface{} {
	if encryptionProperty == nil {
		return nil
	}
	encryption := make(map[string]interface{})
	encryption["enabled"] = strings.EqualFold(string(encryptionProperty.Status), string(containerregistry.EncryptionStatusEnabled))
	if encryptionProperty.KeyVaultProperties != nil {
		encryption["key_vault_key_id"] = encryptionProperty.KeyVaultProperties.KeyIdentifier
		encryption["identity_client_id"] = encryptionProperty.KeyVaultProperties.Identity
	}

	return []interface{}{encryption}
}

func flattenRegistryIdentity(input *containerregistry.IdentityProperties) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func flattenNetworkRuleSet(networkRuleSet *containerregistry.NetworkRuleSet) []interface{} {
	if networkRuleSet == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	values["default_action"] = string(networkRuleSet.DefaultAction)

	ipRules := make([]interface{}, 0)
	for _, ipRule := range *networkRuleSet.IPRules {
		value := make(map[string]interface{})
		value["action"] = string(ipRule.Action)

		// When a /32 CIDR is passed as an ip rule, Azure will drop the /32 leading to the resource wanting to be re-created next run
		if !strings.Contains(*ipRule.IPAddressOrRange, "/") {
			*ipRule.IPAddressOrRange += "/32"
		}

		value["ip_range"] = ipRule.IPAddressOrRange
		ipRules = append(ipRules, value)
	}

	values["ip_rule"] = ipRules

	virtualNetworkRules := make([]interface{}, 0)

	if networkRuleSet.VirtualNetworkRules != nil {
		for _, virtualNetworkRule := range *networkRuleSet.VirtualNetworkRules {
			value := make(map[string]interface{})
			value["action"] = string(virtualNetworkRule.Action)

			value["subnet_id"] = virtualNetworkRule.VirtualNetworkResourceID
			virtualNetworkRules = append(virtualNetworkRules, value)
		}
	}

	values["virtual_network"] = virtualNetworkRules

	return []interface{}{values}
}

func flattenQuarantinePolicy(p *containerregistry.Policies) bool {
	if p == nil || p.QuarantinePolicy == nil {
		return false
	}

	return p.QuarantinePolicy.Status == containerregistry.PolicyStatusEnabled
}

func flattenRetentionPolicy(p *containerregistry.Policies) []interface{} {
	if p == nil || p.RetentionPolicy == nil {
		return nil
	}

	r := *p.RetentionPolicy
	retentionPolicy := make(map[string]interface{})
	retentionPolicy["days"] = r.Days
	enabled := strings.EqualFold(string(r.Status), string(containerregistry.PolicyStatusEnabled))
	retentionPolicy["enabled"] = utils.Bool(enabled)
	return []interface{}{retentionPolicy}
}

func flattenTrustPolicy(p *containerregistry.Policies) []interface{} {
	if p == nil || p.TrustPolicy == nil {
		return nil
	}

	t := *p.TrustPolicy
	trustPolicy := make(map[string]interface{})
	enabled := strings.EqualFold(string(t.Status), string(containerregistry.PolicyStatusEnabled))
	trustPolicy["enabled"] = utils.Bool(enabled)
	return []interface{}{trustPolicy}
}

func resourceContainerRegistrySchema() map[string]*pluginsdk.Schema {
	out := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerValidate.ContainerRegistryName,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"sku": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Default:          string(containerregistry.SkuNameClassic),
			DiffSuppressFunc: suppress.CaseDifferenceV2Only,
			ValidateFunc: validation.StringInSlice([]string{
				string(containerregistry.SkuNameClassic),
				string(containerregistry.SkuNameBasic),
				string(containerregistry.SkuNameStandard),
				string(containerregistry.SkuNamePremium),
			}, !features.ThreePointOhBeta()),
		},

		"admin_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"georeplications": {
			// Don't make this a TypeSet since TypeSet has bugs when there is a nested property using `StateFunc`.
			// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: !features.ThreePointOhBeta(),
			ConflictsWith: func() []string {
				if !features.ThreePointOhBeta() {
					return []string{"georeplication_locations"}
				}
				return []string{}
			}(),
			ConfigMode: func() schema.SchemaConfigMode {
				if !features.ThreePointOhBeta() {
					return pluginsdk.SchemaConfigModeAttr
				}
				return pluginsdk.SchemaConfigModeAuto
			}(),
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

					"tags": tags.Schema(),
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
						Default:  containerregistry.DefaultActionAllow,
						ValidateFunc: validation.StringInSlice([]string{
							string(containerregistry.DefaultActionAllow),
							string(containerregistry.DefaultActionDeny),
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
										string(containerregistry.ActionAllow),
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
						Type:       pluginsdk.TypeSet,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(containerregistry.ActionAllow),
									}, false),
								},
								"subnet_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: azure.ValidateResourceID,
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
				string(containerregistry.NetworkRuleBypassOptionsAzureServices),
				string(containerregistry.NetworkRuleBypassOptionsNone),
			}, false),
			Default: string(containerregistry.NetworkRuleBypassOptionsAzureServices),
		},

		"tags": tags.Schema(),
	}

	if !features.ThreePointOhBeta() {

		out["storage_account_id"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Optional:   true,
			Computed:   true,
			ForceNew:   true,
			Deprecated: "this attribute is no longer recognized by the API and is not functional anymore, thus this property will be removed in v3.0",
		}

		out["georeplication_locations"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeSet,
			Optional:      true,
			Deprecated:    "Deprecated in favour of `georeplications`",
			Computed:      true,
			ConflictsWith: []string{"georeplications"},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			Set: location.HashCode,
		}
	}

	return out
}
