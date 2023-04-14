package appconfiguration

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/deletedconfigurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppConfigurationCreate,
		Read:   resourceAppConfigurationRead,
		Update: resourceAppConfigurationUpdate,
		Delete: resourceAppConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := configurationstores.ParseConfigurationStoreID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ConfigurationStoreName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"encryption": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_identifier": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},
						"identity_client_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
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

			"purge_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "free",
				ValidateFunc: validation.StringInSlice([]string{
					"free",
					"standard",
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

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringInSlice(configurationstores.PossibleValuesForPublicNetworkAccess(), true),
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
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
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
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
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
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
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
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
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

	location := azure.NormalizeLocation(d.Get("location").(string))
	deletedConfigurationStoresId := deletedconfigurationstores.NewDeletedConfigurationStoreID(subscriptionId, location, name)
	deleted, err := deletedConfigurationStoresClient.ConfigurationStoresGetDeleted(ctx, deletedConfigurationStoresId)
	if err != nil {
		if !response.WasNotFound(deleted.HttpResponse) {
			return fmt.Errorf("checking for presence of deleted %s: %+v", deletedConfigurationStoresId, err)
		}
	}

	recoverSoftDeleted := false
	if !response.WasNotFound(deleted.HttpResponse) && !response.WasStatusCode(deleted.HttpResponse, http.StatusForbidden) {
		if !meta.(*clients.Client).Features.AppConfiguration.RecoverSoftDeleted {
			return fmt.Errorf(optedOutOfRecoveringSoftDeletedAppConfigurationErrorFmt(name, location))
		}
		recoverSoftDeleted = true
	}

	parameters := configurationstores.ConfigurationStore{
		Location: location,
		Sku: configurationstores.Sku{
			Name: d.Get("sku").(string),
		},
		Properties: &configurationstores.ConfigurationStoreProperties{
			EnablePurgeProtection: utils.Bool(d.Get("purge_protection_enabled").(bool)),
			DisableLocalAuth:      utils.Bool(!d.Get("local_auth_enabled").(bool)),
			Encryption:            expandAppConfigurationEncryption(d.Get("encryption").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.Get("soft_delete_retention_days").(int); ok && v != 7 {
		parameters.Properties.SoftDeleteRetentionInDays = utils.Int64(int64(v))
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
	// TODO: retry checkNameAvailability before creation when SDK is ready, see https://github.com/Azure/AppConfiguration/issues/677
	if err := client.CreateThenPoll(ctx, resourceId, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
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
		update.Properties.DisableLocalAuth = utils.Bool(!d.Get("local_auth_enabled").(bool))
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
		update.Properties.EnablePurgeProtection = utils.Bool(d.Get("purge_protection_enabled").(bool))
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
				return fmt.Errorf("retrieving the Deletion Details for %s: %+v", *id, err)
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
		if err := deletedConfigurationStoresClient.ConfigurationStoresPurgeDeletedThenPoll(ctx, deletedId); err != nil {
			return fmt.Errorf("purging %s: %+v", *id, err)
		}

		// TODO: retry checkNameAvailability after deletion when SDK is ready, see https://github.com/Azure/AppConfiguration/issues/677
		log.Printf("[DEBUG] Purged AppConfiguration %q.", id.ConfigurationStoreName)
	}

	return nil
}

type flattenedAccessKeys struct {
	primaryReadKey    []interface{}
	primaryWriteKey   []interface{}
	secondaryReadKey  []interface{}
	secondaryWriteKey []interface{}
}

func expandAppConfigurationEncryption(input []interface{}) *configurationstores.EncryptionProperties {
	if len(input) == 0 {
		return nil
	}

	encryptionParam := input[0].(map[string]interface{})
	result := &configurationstores.EncryptionProperties{
		KeyVaultProperties: &configurationstores.KeyVaultProperties{},
	}

	if v, ok := encryptionParam["identity_client_id"].(string); ok && v != "" {
		result.KeyVaultProperties.IdentityClientId = &v
	}
	if v, ok := encryptionParam["key_vault_key_identifier"].(string); ok && v != "" {
		result.KeyVaultProperties.KeyIdentifier = &v
	}
	return result
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

func optedOutOfRecoveringSoftDeletedAppConfigurationErrorFmt(name, location string) string {
	return fmt.Sprintf(`
An existing soft-deleted App Configuration exists with the Name %q in the location %q, however
automatically recovering this App Configuration has been disabled via the "features" block.

Terraform can automatically recover the soft-deleted App Configuration when this behaviour is
enabled within the "features" block (located within the "provider" block) - more information
can be found here:

https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#features

Alternatively you can manually recover this (e.g. using the Azure CLI) and then import
this into Terraform via "terraform import", or pick a different name/location.
`, name, location)
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
