// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
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
	resource := &pluginsdk.Resource{
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

			"premium_messaging_partitions": {
				Type:         pluginsdk.TypeInt,
				ForceNew:     true,
				Default:      0,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 4}),
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

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"minimum_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(namespaces.TlsVersionOnePointTwo),
				ValidateFunc: validation.StringInSlice([]string{
					string(namespaces.TlsVersionOnePointZero),
					string(namespaces.TlsVersionOnePointOne),
					string(namespaces.TlsVersionOnePointTwo),
				}, false),
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

			"network_rule_set": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(namespaces.DefaultActionAllow),
							ValidateFunc: validation.StringInSlice([]string{
								string(namespaces.DefaultActionAllow),
								string(namespaces.DefaultActionDeny),
							}, false),
						},

						"public_network_access_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"ip_rules": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"trusted_services_allowed": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"network_rules": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Set:      networkRuleHash,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"subnet_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: commonids.ValidateSubnetID,
										// The subnet ID returned from the service will have `resourceGroup/{resourceGroupName}` all in lower cases...
										DiffSuppressFunc: suppress.CaseDifference,
									},
									"ignore_missing_vnet_service_endpoint": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
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
			pluginsdk.CustomizeDiffShim(servicebusTLSVersionDiff),
		),
	}

	if !features.FourPointOhBeta() {
		resource.Schema["zone_redundant"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Optional:   true,
			Computed:   true,
			Deprecated: "The `zone_redundant` property has been deprecated and will be removed in v4.0 of the provider.",
			ForceNew:   true,
		}

		resource.Schema["minimum_tls_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(namespaces.TlsVersionOnePointZero),
				string(namespaces.TlsVersionOnePointOne),
				string(namespaces.TlsVersionOnePointTwo),
			}, false),
		}
	}

	return resource
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

	publicNetworkEnabled := namespaces.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkEnabled = namespaces.PublicNetworkAccessDisabled
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
			Encryption:          expandServiceBusNamespaceEncryption(d.Get("customer_managed_key").([]interface{})),
			DisableLocalAuth:    utils.Bool(!d.Get("local_auth_enabled").(bool)),
			PublicNetworkAccess: &publicNetworkEnabled,
		},
		Tags: expandTags(t),
	}

	if !features.FourPointOhBeta() {
		parameters.Properties.ZoneRedundant = utils.Bool(d.Get("zone_redundant").(bool))
	}

	if tlsValue := d.Get("minimum_tls_version").(string); tlsValue != "" {
		minimumTls := namespaces.TlsVersion(tlsValue)
		parameters.Properties.MinimumTlsVersion = &minimumTls
	}

	if capacity := d.Get("capacity"); capacity != nil {
		if !strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && capacity.(int) > 0 {
			return fmt.Errorf("service bus SKU %q only supports `capacity` of 0", sku)
		}
		if strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && capacity.(int) == 0 {
			return fmt.Errorf("service bus SKU %q only supports `capacity` of 1, 2, 4, 8 or 16", sku)
		}
		parameters.Sku.Capacity = utils.Int64(int64(capacity.(int)))
	}

	if premiumMessagingUnit := d.Get("premium_messaging_partitions"); premiumMessagingUnit != nil {
		if !strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && premiumMessagingUnit.(int) > 0 {
			return fmt.Errorf("premium messaging partition is not supported by service Bus SKU %q and it can only be set to 0", sku)
		}
		if strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && premiumMessagingUnit.(int) == 0 {
			return fmt.Errorf("service bus SKU %q only supports `premium_messaging_partitions` of 1, 2, 4", sku)
		}
		parameters.Properties.PremiumMessagingPartitions = utils.Int64(int64(premiumMessagingUnit.(int)))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if d.HasChange("network_rule_set") {
		oldNetworkRuleSet, newNetworkRuleSet := d.GetChange("network_rule_set")
		// if the network rule set has been removed from config, reset it instead as there is no way to remove a rule set
		if len(oldNetworkRuleSet.([]interface{})) == 1 && len(newNetworkRuleSet.([]interface{})) == 0 {
			log.Printf("[DEBUG] Resetting Network Rule Set associated with %s..", id)
			if err = resetNetworkRuleSetForNamespace(ctx, client, id); err != nil {
				return err
			}
			log.Printf("[DEBUG] Reset the Existing Network Rule Set associated with %s", id)
		} else {
			log.Printf("[DEBUG] Creating the Network Rule Set associated with %s..", id)
			if err = createNetworkRuleSetForNamespace(ctx, client, id, newNetworkRuleSet.([]interface{})); err != nil {
				return err
			}
			log.Printf("[DEBUG] Created the Network Rule Set associated with %s", id)
		}
	}

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

			if props := model.Properties; props != nil {
				d.Set("premium_messaging_partitions", props.PremiumMessagingPartitions)

				if !features.FourPointOhBeta() {
					d.Set("zone_redundant", props.ZoneRedundant)
				}

				if customerManagedKey, err := flattenServiceBusNamespaceEncryption(props.Encryption); err == nil {
					d.Set("customer_managed_key", customerManagedKey)
				}
				localAuthEnabled := true
				if props.DisableLocalAuth != nil {
					localAuthEnabled = !*props.DisableLocalAuth
				}
				d.Set("local_auth_enabled", localAuthEnabled)

				publicNetworkAccess := true
				if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == namespaces.PublicNetworkAccessDisabled {
					publicNetworkAccess = false
				}
				d.Set("public_network_access_enabled", publicNetworkAccess)

				if props.MinimumTlsVersion != nil {
					d.Set("minimum_tls_version", string(pointer.From(props.MinimumTlsVersion)))
				}

				d.Set("endpoint", props.ServiceBusEndpoint)
			}
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

	networkRuleSet, err := client.GetNetworkRuleSet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving network rule set %s: %+v", *id, err)
	}

	if model := networkRuleSet.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("network_rule_set", flattenServiceBusNamespaceNetworkRuleSet(*props))
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

	// need to wait the status to be ready before performing the deleting.
	if err := waitForNamespaceStatusToBeReady(ctx, meta, *id, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for serviceBus namespace %s state to be ready error: %+v", *id, err)
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
		keyVaultKeyId, err := keyVaultParse.NewNestedItemID(pointer.From(props.KeyVaultUri), keyVaultParse.NestedItemTypeKey, pointer.From(props.KeyName), pointer.From(props.KeyVersion))
		if err != nil {
			return nil, fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
		}
		keyId = keyVaultKeyId.ID()
		if props.Identity != nil && props.Identity.UserAssignedIdentity != nil {
			sbnUaiId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*props.Identity.UserAssignedIdentity)
			if err != nil {
				return nil, err
			}
			identityId = sbnUaiId.ID()
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

func servicebusTLSVersionDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) (err error) {
	old, new := d.GetChange("minimum_tls_version")
	if old != "" && new == "" {
		err = fmt.Errorf("`minimum_tls_version` has been set before, please set a valid value for this property ")
	}
	return
}

func createNetworkRuleSetForNamespace(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId, input []interface{}) error {
	if len(input) < 1 || input[0] == nil {
		return nil
	}
	item := input[0].(map[string]interface{})

	defaultAction := namespaces.DefaultAction(item["default_action"].(string))
	vnetRule := expandServiceBusNamespaceVirtualNetworkRules(item["network_rules"].(*pluginsdk.Set).List())
	ipRule := expandServiceBusNamespaceIPRules(item["ip_rules"].(*pluginsdk.Set).List())
	publicNetworkAcc := "Disabled"
	if item["public_network_access_enabled"].(bool) {
		publicNetworkAcc = "Enabled"
	}

	// API doesn't accept "Deny" to be set for "default_action" if no "ip_rules" or "network_rules" is defined and returns no error message to the user
	if defaultAction == namespaces.DefaultActionDeny && vnetRule == nil && ipRule == nil {
		return fmt.Errorf(" The `default_action` of `network_rule_set` can only be set to `Allow` if no `ip_rules` or `network_rules` is set")
	}

	publicNetworkAccess := namespaces.PublicNetworkAccessFlag(publicNetworkAcc)

	parameters := namespaces.NetworkRuleSet{
		Properties: &namespaces.NetworkRuleSetProperties{
			DefaultAction:               &defaultAction,
			VirtualNetworkRules:         vnetRule,
			IPRules:                     ipRule,
			PublicNetworkAccess:         &publicNetworkAccess,
			TrustedServiceAccessEnabled: utils.Bool(item["trusted_services_allowed"].(bool)),
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	return nil
}

func resetNetworkRuleSetForNamespace(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId) error {
	defaultAction := namespaces.DefaultActionAllow
	parameters := namespaces.NetworkRuleSet{
		Properties: &namespaces.NetworkRuleSetProperties{
			DefaultAction: &defaultAction,
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, id, parameters); err != nil {
		return fmt.Errorf("resetting %s: %+v", id, err)
	}

	return nil
}

func flattenServiceBusNamespaceNetworkRuleSet(networkRuleSet namespaces.NetworkRuleSetProperties) []interface{} {
	defaultAction := ""
	if v := networkRuleSet.DefaultAction; v != nil {
		defaultAction = string(*v)
	}
	publicNetworkAccess := namespaces.PublicNetworkAccessFlagEnabled
	if v := networkRuleSet.PublicNetworkAccess; v != nil {
		publicNetworkAccess = *v
	}

	trustedServiceEnabled := false
	if networkRuleSet.TrustedServiceAccessEnabled != nil {
		trustedServiceEnabled = *networkRuleSet.TrustedServiceAccessEnabled
	}

	networkRules := flattenServiceBusNamespaceVirtualNetworkRules(networkRuleSet.VirtualNetworkRules)
	ipRules := flattenServiceBusNamespaceIPRules(networkRuleSet.IPRules)

	// only set network rule set if the values are different than what they are defaulted to during namespace creation
	// this has to wait until 4.0 due to `azurerm_servicebus_namespace_network_rule_set` which forces `network_rule_set` to be Optional/Computed
	if features.FourPointOhBeta() {
		if defaultAction == string(namespaces.DefaultActionAllow) &&
			publicNetworkAccess == namespaces.PublicNetworkAccessFlagEnabled &&
			!trustedServiceEnabled &&
			len(networkRules) == 0 &&
			len(ipRules) == 0 {

			return []interface{}{}
		}
	}

	return []interface{}{map[string]interface{}{
		"default_action":                defaultAction,
		"trusted_services_allowed":      trustedServiceEnabled,
		"public_network_access_enabled": publicNetworkAccess == namespaces.PublicNetworkAccessFlagEnabled,
		"network_rules":                 pluginsdk.NewSet(networkRuleHash, networkRules),
		"ip_rules":                      ipRules,
	}}
}
