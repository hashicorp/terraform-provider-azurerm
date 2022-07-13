package servicebus

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var (
	serviceBusNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"
	serviceBusNamespaceResourceName             = "azurerm_servicebus_namespace"
)

func resourceServiceBusNamespace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusNamespaceCreateUpdate,
		Read:   resourceServiceBusNamespaceRead,
		Update: resourceServiceBusNamespaceCreateUpdate,
		Delete: resourceServiceBusNamespaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := namespaces.ParseNamespaceID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NamespaceV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NamespaceName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(namespaces.SkuNameBasic),
					string(namespaces.SkuNameStandard),
					string(namespaces.SkuNamePremium),
				}, false),
			},

			"capacity": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 4, 8, 16}),
			},

			"customer_managed_key": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
						},

						"identity_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},

						"infrastructure_encryption_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"default_primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			oldCustomerManagedKey, newCustomerManagedKey := diff.GetChange("customer_managed_key")
			if len(oldCustomerManagedKey.([]interface{})) != 0 && len(newCustomerManagedKey.([]interface{})) == 0 {
				diff.ForceNew("customer_managed_key")
			}

			oldSku, newSku := diff.GetChange("sku")
			if diff.HasChange("sku") {
				if strings.EqualFold(newSku.(string), string(namespaces.SkuNamePremium)) || strings.EqualFold(oldSku.(string), string(namespaces.SkuNamePremium)) {
					log.Printf("[DEBUG] cannot migrate a namespace from or to Premium SKU")
					diff.ForceNew("sku")
				}
			}
			return nil
		}),
	}
}

func resourceServiceBusNamespaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for ServiceBus Namespace create/update.")

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	t := d.Get("tags").(map[string]interface{})

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace", id.ID())
		}
	}

	identity, err := expandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	s := namespaces.SkuTier(sku)
	parameters := namespaces.SBNamespace{
		Location: location,
		Identity: identity,
		Sku: &namespaces.SBSku{
			Name: namespaces.SkuName(sku),
			Tier: &s,
		},
		Properties: &namespaces.SBNamespaceProperties{
			ZoneRedundant:    utils.Bool(d.Get("zone_redundant").(bool)),
			Encryption:       expandServiceBusNamespaceEncryption(d.Get("customer_managed_key").([]interface{})),
			DisableLocalAuth: utils.Bool(!d.Get("local_auth_enabled").(bool)),
		},
		Tags: expandTags(t),
	}

	if capacity := d.Get("capacity"); capacity != nil {
		if !strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && capacity.(int) > 0 {
			return fmt.Errorf("Service Bus SKU %q only supports `capacity` of 0", sku)
		}
		if strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && capacity.(int) == 0 {
			return fmt.Errorf("Service Bus SKU %q only supports `capacity` of 1, 2, 4, 8 or 16", sku)
		}
		parameters.Sku.Capacity = utils.Int64(int64(capacity.(int)))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceServiceBusNamespaceRead(d, meta)
}

func resourceServiceBusNamespaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	namespaceAuthClient := meta.(*clients.Client).ServiceBus.NamespacesAuthClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		d.Set("tags", flattenTags(model.Tags))

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if sku := model.Sku; sku != nil {
			skuName := ""
			// the Azure API is inconsistent here, so rewrite this into the casing we expect
			for _, v := range namespaces.PossibleValuesForSkuName() {
				if strings.EqualFold(v, string(sku.Name)) {
					skuName = v
				}
			}
			d.Set("sku", skuName)
			d.Set("capacity", sku.Capacity)
		}

		if props := model.Properties; model != nil {
			d.Set("zone_redundant", props.ZoneRedundant)
			if customerManagedKey, err := flattenServiceBusNamespaceEncryption(props.Encryption); err == nil {
				d.Set("customer_managed_key", customerManagedKey)
			}
			localAuthEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !*props.DisableLocalAuth
			}
			d.Set("local_auth_enabled", localAuthEnabled)
		}
	}

	authRuleId := namespacesauthorizationrule.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, serviceBusNamespaceDefaultAuthorizationRule)

	keys, err := namespaceAuthClient.NamespacesListKeys(ctx, authRuleId)
	if err != nil {
		log.Printf("[WARN] listing default keys for %s: %+v", id, err)
	} else {
		if keysModel := keys.Model; keysModel != nil {
			d.Set("default_primary_connection_string", keysModel.PrimaryConnectionString)
			d.Set("default_secondary_connection_string", keysModel.SecondaryConnectionString)
			d.Set("default_primary_key", keysModel.PrimaryKey)
			d.Set("default_secondary_key", keysModel.SecondaryKey)
		}
	}
	return nil
}

func resourceServiceBusNamespaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandServiceBusNamespaceEncryption(input []interface{}) *namespaces.Encryption {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	keyId, _ := keyVaultParse.ParseOptionallyVersionedNestedItemID(v["key_vault_key_id"].(string))
	keySource := namespaces.KeySourceMicrosoftPointKeyVault
	return &namespaces.Encryption{
		KeyVaultProperties: &[]namespaces.KeyVaultProperties{
			{
				KeyName:     utils.String(keyId.Name),
				KeyVersion:  utils.String(keyId.Version),
				KeyVaultUri: utils.String(keyId.KeyVaultBaseUrl),
				Identity: &namespaces.UserAssignedIdentityProperties{
					UserAssignedIdentity: utils.String(v["identity_id"].(string)),
				},
			},
		},
		KeySource:                       &keySource,
		RequireInfrastructureEncryption: utils.Bool(v["infrastructure_encryption_enabled"].(bool)),
	}
}

func flattenServiceBusNamespaceEncryption(encryption *namespaces.Encryption) ([]interface{}, error) {
	if encryption == nil {
		return []interface{}{}, nil
	}

	var keyId string
	var identityId string
	if keyVaultProperties := encryption.KeyVaultProperties; keyVaultProperties != nil && len(*keyVaultProperties) != 0 {
		props := (*keyVaultProperties)[0]
		keyVaultKeyId, err := keyVaultParse.NewNestedItemID(*props.KeyVaultUri, "keys", *props.KeyName, *props.KeyVersion)
		if err != nil {
			return nil, fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
		}
		keyId = keyVaultKeyId.ID()
		if props.Identity != nil && props.Identity.UserAssignedIdentity != nil {
			identityId = *props.Identity.UserAssignedIdentity
		}
	}

	return []interface{}{
		map[string]interface{}{
			"infrastructure_encryption_enabled": encryption.RequireInfrastructureEncryption,
			"key_vault_key_id":                  keyId,
			"identity_id":                       identityId,
		},
	}, nil
}

func expandSystemAndUserAssignedMap(input []interface{}) (*identity.SystemAndUserAssignedMap, error) {
	identityType := identity.TypeNone
	identityIds := make(map[string]identity.UserAssignedIdentityDetails, 0)

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})
		typeRaw := raw["type"].(string)
		if typeRaw == string(identity.TypeSystemAssigned) {
			identityType = identity.TypeSystemAssigned
		}
		if typeRaw == string(identity.TypeSystemAssignedUserAssigned) {
			identityType = identity.TypeSystemAssignedUserAssigned
		}
		if typeRaw == string(identity.TypeUserAssigned) {
			identityType = identity.TypeUserAssigned
		}

		identityIdsRaw := raw["identity_ids"].(*schema.Set).List()
		for _, v := range identityIdsRaw {
			identityIds[v.(string)] = identity.UserAssignedIdentityDetails{
				// intentionally empty since the expand shouldn't send these values
			}
		}
	}

	if len(identityIds) > 0 && (identityType != identity.TypeSystemAssignedUserAssigned && identityType != identity.TypeUserAssigned) {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q or %q", string(identity.TypeSystemAssignedUserAssigned), string(identity.TypeUserAssigned))
	}

	if len(identityIds) == 0 {
		return &identity.SystemAndUserAssignedMap{
			Type: identityType,
		}, nil
	}

	return &identity.SystemAndUserAssignedMap{
		Type:        identityType,
		IdentityIds: identityIds,
	}, nil
}
