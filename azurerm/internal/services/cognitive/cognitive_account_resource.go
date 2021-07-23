package cognitive

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2021-04-30/cognitiveservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	commonValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive/validate"
	msiparse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	msiValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	networkParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCognitiveAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCognitiveAccountCreate,
		Read:   resourceCognitiveAccountRead,
		Update: resourceCognitiveAccountUpdate,
		Delete: resourceCognitiveAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CognitiveServicesAccountName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Academic",
					"AnomalyDetector",
					"Bing.Autosuggest",
					"Bing.Autosuggest.v7",
					"Bing.CustomSearch",
					"Bing.Search",
					"Bing.Search.v7",
					"Bing.Speech",
					"Bing.SpellCheck",
					"Bing.SpellCheck.v7",
					"CognitiveServices",
					"ComputerVision",
					"ContentModerator",
					"CustomSpeech",
					"CustomVision.Prediction",
					"CustomVision.Training",
					"Emotion",
					"Face",
					"FormRecognizer",
					"ImmersiveReader",
					"LUIS",
					"LUIS.Authoring",
					"MetricsAdvisor",
					"Personalizer",
					"QnAMaker",
					"Recommendations",
					"SpeakerRecognition",
					"Speech",
					"SpeechServices",
					"SpeechTranslation",
					"TextAnalytics",
					"TextTranslation",
					"WebLM",
				}, false),
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"F0", "F1", "S0", "S", "S1", "S2", "S3", "S4", "S5", "S6", "P0", "P1", "P2",
				}, false),
			},

			"custom_subdomain_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"fqdns": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cognitiveservices.ResourceIdentityTypeNone),
							ValidateFunc: validation.StringInSlice([]string{
								string(cognitiveservices.ResourceIdentityTypeNone),
								string(cognitiveservices.ResourceIdentityTypeSystemAssigned),
								string(cognitiveservices.ResourceIdentityTypeUserAssigned),
								string(cognitiveservices.ResourceIdentityTypeSystemAssignedUserAssigned),
							}, false),
						},

						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"identity_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: msiValidate.UserAssignedIdentityID,
							},
						},
					},
				},
			},

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"metrics_advisor_aad_client_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"metrics_advisor_aad_tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"metrics_advisor_super_user_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"metrics_advisor_website_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"network_acls": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MaxItems:     1,
				RequiredWith: []string{"custom_subdomain_name"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cognitiveservices.NetworkRuleActionAllow),
								string(cognitiveservices.NetworkRuleActionDeny),
							}, false),
						},
						"ip_rules": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.Any(
									commonValidate.IPv4Address,
									commonValidate.CIDR,
								),
							},
							Set: set.HashIPv4AddressOrCIDR,
						},
						// TODO 3.0 - Remove below property
						"virtual_network_subnet_ids": {
							Type:          pluginsdk.TypeSet,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"network_acls.0.virtual_network_rules"},
							Deprecated:    "Deprecated in favour of `virtual_network_rules`",
							Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},

						"virtual_network_rules": {
							Type:          pluginsdk.TypeSet,
							Optional:      true,
							Computed:      true, // TODO -- remove this when deprecation resolves
							ConflictsWith: []string{"network_acls.0.virtual_network_subnet_ids"},
							ConfigMode:    pluginsdk.SchemaConfigModeAttr, // TODO -- remove in 3.0, because this property is optional and computed, it has to be declared as empty array to remove existed values
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"subnet_id": {
										Type:     pluginsdk.TypeString,
										Required: true,
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

			"outbound_network_access_restrited": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"qna_runtime_endpoint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			"storage": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"storage_account_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: storageValidate.StorageAccountID,
						},

						"identity_client_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"tags": tags.Schema(),

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceCognitiveAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	kind := d.Get("kind").(string)

	id := parse.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cognitive_account", id.ID())
		}
	}

	sku, err := expandAccountSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku_name for %s: %v", id, err)
	}

	networkAcls, subnetIds := expandCognitiveAccountNetworkAcls(d)

	// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
	virtualNetworkNames := make([]string, 0)
	for _, v := range subnetIds {
		id, err := networkParse.SubnetIDInsensitively(v)
		if err != nil {
			return err
		}
		if !utils.SliceContainsValue(virtualNetworkNames, id.VirtualNetworkName) {
			virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	publicNetworkAccess := cognitiveservices.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = cognitiveservices.PublicNetworkAccessDisabled
	}

	apiProps, err := expandCognitiveAccountAPIProperties(d)
	if err != nil {
		return err
	}

	props := cognitiveservices.Account{
		Kind:     utils.String(kind),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Sku:      sku,
		Properties: &cognitiveservices.AccountProperties{
			APIProperties:                 apiProps,
			NetworkAcls:                   networkAcls,
			CustomSubDomainName:           utils.String(d.Get("custom_subdomain_name").(string)),
			AllowedFqdnList:               utils.ExpandStringSlice(d.Get("fqdns").([]interface{})),
			PublicNetworkAccess:           publicNetworkAccess,
			UserOwnedStorage:              expandCognitiveAccountStorage(d.Get("storage").([]interface{})),
			RestrictOutboundNetworkAccess: utils.Bool(d.Get("outbound_network_access_restrited").(bool)),
			DisableLocalAuth:              utils.Bool(!d.Get("local_auth_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	identityRaw := d.Get("identity").([]interface{})
	identity, err := expandCognitiveAccountIdentity(identityRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `identity`: %+v", err)
	}
	props.Identity = identity

	if _, err := client.Create(ctx, id.ResourceGroup, id.Name, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    cognitiveAccountStateRefreshFunc(ctx, client, id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCognitiveAccountRead(d, meta)
}

func resourceCognitiveAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	sku, err := expandAccountSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("error expanding sku_name for %s: %+v", *id, err)
	}

	networkAcls, subnetIds := expandCognitiveAccountNetworkAcls(d)

	// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
	virtualNetworkNames := make([]string, 0)
	for _, v := range subnetIds {
		id, err := networkParse.SubnetIDInsensitively(v)
		if err != nil {
			return err
		}
		if !utils.SliceContainsValue(virtualNetworkNames, id.VirtualNetworkName) {
			virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	publicNetworkAccess := cognitiveservices.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = cognitiveservices.PublicNetworkAccessDisabled
	}

	apiProps, err := expandCognitiveAccountAPIProperties(d)
	if err != nil {
		return err
	}

	props := cognitiveservices.Account{
		Sku: sku,
		Properties: &cognitiveservices.AccountProperties{
			APIProperties:                 apiProps,
			NetworkAcls:                   networkAcls,
			CustomSubDomainName:           utils.String(d.Get("custom_subdomain_name").(string)),
			AllowedFqdnList:               utils.ExpandStringSlice(d.Get("fqdns").([]interface{})),
			PublicNetworkAccess:           publicNetworkAccess,
			UserOwnedStorage:              expandCognitiveAccountStorage(d.Get("storage").([]interface{})),
			RestrictOutboundNetworkAccess: utils.Bool(d.Get("outbound_network_access_restrited").(bool)),
			DisableLocalAuth:              utils.Bool(!d.Get("local_auth_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	identityRaw := d.Get("identity").([]interface{})
	identity, err := expandCognitiveAccountIdentity(identityRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `identity`: %+v", err)
	}
	props.Identity = identity

	if _, err = client.Update(ctx, id.ResourceGroup, id.Name, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Accepted"},
		Target:     []string{"Succeeded"},
		Refresh:    cognitiveAccountStateRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}
	return resourceCognitiveAccountRead(d, meta)
}

func resourceCognitiveAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		// note for the resource we shouldn't gracefully fail since we have permission to CRUD it
		return fmt.Errorf("listing the Keys for %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("kind", resp.Kind)

	d.Set("location", location.NormalizeNilable(resp.Location))
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	identity, err := flattenCognitiveAccountIdentity(resp.Identity)
	if err != nil {
		return err
	}
	d.Set("identity", identity)

	if props := resp.Properties; props != nil {
		if apiProps := props.APIProperties; apiProps != nil {
			d.Set("qna_runtime_endpoint", apiProps.QnaRuntimeEndpoint)
			d.Set("metrics_advisor_aad_client_id", apiProps.AadClientID)
			d.Set("metrics_advisor_aad_tenant_id", apiProps.AadTenantID)
			d.Set("metrics_advisor_super_user_name", apiProps.SuperUser)
			d.Set("metrics_advisor_website_name", apiProps.WebsiteName)
		}
		d.Set("endpoint", props.Endpoint)
		d.Set("custom_subdomain_name", props.CustomSubDomainName)
		if err := d.Set("network_acls", flattenCognitiveAccountNetworkAcls(props.NetworkAcls)); err != nil {
			return fmt.Errorf("setting `network_acls` for Cognitive Account %q: %+v", id, err)
		}
		d.Set("fqdns", utils.FlattenStringSlice(props.AllowedFqdnList))
		d.Set("public_network_access_enabled", props.PublicNetworkAccess == cognitiveservices.PublicNetworkAccessEnabled)
		if err := d.Set("storage", flattenCognitiveAccountStorage(props.UserOwnedStorage)); err != nil {
			return fmt.Errorf("setting `storages` for Cognitive Account %q: %+v", id, err)
		}
		if props.RestrictOutboundNetworkAccess != nil {
			d.Set("outbound_network_access_restrited", props.RestrictOutboundNetworkAccess)
		}
		if props.DisableLocalAuth != nil {
			d.Set("local_auth_enabled", !*props.DisableLocalAuth)
		}
	}

	d.Set("primary_access_key", keys.Key1)
	d.Set("secondary_access_key", keys.Key2)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceCognitiveAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	accountsClient := meta.(*clients.Client).Cognitive.AccountsClient
	deletedAccountsClient := meta.(*clients.Client).Cognitive.DeletedAccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	// first we need to retrieve it, since we need the location to be able to purge it
	log.Printf("[DEBUG] Retrieving %s..", *id)
	account, err := accountsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	log.Printf("[DEBUG] Deleting %s..", *id)
	deleteFuture, err := accountsClient.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	if err := deleteFuture.WaitForCompletionRef(ctx, accountsClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	if meta.(*clients.Client).Features.CognitiveAccount.PurgeSoftDeleteOnDestroy {
		log.Printf("[DEBUG] Purging %s..", *id)
		purgeFuture, err := deletedAccountsClient.Purge(ctx, *account.Location, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("purging %s: %+v", *id, err)
		}
		if err := purgeFuture.WaitForCompletionRef(ctx, deletedAccountsClient.Client); err != nil {
			return fmt.Errorf("waiting for purge of %s: %+v", *id, err)
		}
	} else {
		log.Printf("[DEBUG] Skipping Purge of %s", *id)
	}

	return nil
}

func expandAccountSkuName(skuName string) (*cognitiveservices.Sku, error) {
	var tier cognitiveservices.SkuTier
	switch skuName[0:1] {
	case "F":
		tier = cognitiveservices.SkuTierFree
	case "S":
		tier = cognitiveservices.SkuTierStandard
	case "P":
		tier = cognitiveservices.SkuTierPremium
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, skuName[0:1])
	}

	return &cognitiveservices.Sku{
		Name: utils.String(skuName),
		Tier: tier,
	}, nil
}

func cognitiveAccountStateRefreshFunc(ctx context.Context, client *cognitiveservices.AccountsClient, id parse.AccountId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		return res, string(res.Properties.ProvisioningState), nil
	}
}

func expandCognitiveAccountNetworkAcls(d *pluginsdk.ResourceData) (*cognitiveservices.NetworkRuleSet, []string) {
	input := d.Get("network_acls").([]interface{})
	subnetIds := make([]string, 0)
	if len(input) == 0 || input[0] == nil {
		return nil, subnetIds
	}

	v := input[0].(map[string]interface{})

	defaultAction := v["default_action"].(string)

	ipRulesRaw := v["ip_rules"].(*pluginsdk.Set)
	ipRules := make([]cognitiveservices.IPRule, 0)

	for _, v := range ipRulesRaw.List() {
		rule := cognitiveservices.IPRule{
			Value: utils.String(v.(string)),
		}
		ipRules = append(ipRules, rule)
	}

	networkRules := make([]cognitiveservices.VirtualNetworkRule, 0)
	if d.HasChange("network_acls.0.virtual_network_subnet_ids") {
		networkRulesRaw := v["virtual_network_subnet_ids"]
		for _, v := range networkRulesRaw.(*pluginsdk.Set).List() {
			rawId := v.(string)
			subnetIds = append(subnetIds, rawId)
			rule := cognitiveservices.VirtualNetworkRule{
				ID: utils.String(rawId),
			}
			networkRules = append(networkRules, rule)
		}
	}
	if d.HasChange("network_acls.0.virtual_network_rules") {
		networkRulesRaw := v["virtual_network_rules"]
		for _, v := range networkRulesRaw.(*pluginsdk.Set).List() {
			value := v.(map[string]interface{})
			subnetId := value["subnet_id"].(string)
			subnetIds = append(subnetIds, subnetId)
			rule := cognitiveservices.VirtualNetworkRule{
				ID:                               utils.String(subnetId),
				IgnoreMissingVnetServiceEndpoint: utils.Bool(value["ignore_missing_vnet_service_endpoint"].(bool)),
			}
			networkRules = append(networkRules, rule)
		}
	}

	ruleSet := cognitiveservices.NetworkRuleSet{
		DefaultAction:       cognitiveservices.NetworkRuleAction(defaultAction),
		IPRules:             &ipRules,
		VirtualNetworkRules: &networkRules,
	}
	return &ruleSet, subnetIds
}

func expandCognitiveAccountStorage(input []interface{}) *[]cognitiveservices.UserOwnedStorage {
	if len(input) == 0 {
		return nil
	}
	results := make([]cognitiveservices.UserOwnedStorage, 0)
	for _, v := range input {
		value := v.(map[string]interface{})
		results = append(results, cognitiveservices.UserOwnedStorage{
			ResourceID:       utils.String(value["storage_account_id"].(string)),
			IdentityClientID: utils.String(value["identity_client_id"].(string)),
		})
	}
	return &results
}

func expandCognitiveAccountIdentity(vs []interface{}) (*cognitiveservices.Identity, error) {
	if len(vs) == 0 {
		return &cognitiveservices.Identity{
			Type: cognitiveservices.ResourceIdentityTypeNone,
		}, nil
	}

	v := vs[0].(map[string]interface{})
	managedServiceIdentity := cognitiveservices.Identity{
		Type: cognitiveservices.ResourceIdentityType(v["type"].(string)),
	}

	var identityIdSet []interface{}
	if identityIds, ok := v["identity_ids"]; ok {
		identityIdSet = identityIds.(*pluginsdk.Set).List()
	}

	// If type contains `UserAssigned`, `identity_ids` must be specified and have at least 1 element
	if managedServiceIdentity.Type == cognitiveservices.ResourceIdentityTypeUserAssigned || managedServiceIdentity.Type == cognitiveservices.ResourceIdentityTypeSystemAssignedUserAssigned {
		if len(identityIdSet) == 0 {
			return nil, fmt.Errorf("`identity_ids` must have at least 1 element when `type` includes `UserAssigned`")
		}

		userAssignedIdentities := make(map[string]*cognitiveservices.UserAssignedIdentity)
		for _, id := range identityIdSet {
			userAssignedIdentities[id.(string)] = &cognitiveservices.UserAssignedIdentity{}
		}

		managedServiceIdentity.UserAssignedIdentities = userAssignedIdentities
	} else if len(identityIdSet) > 0 {
		// If type does _not_ contain `UserAssigned` (i.e. is set to `SystemAssigned` or defaulted to `None`), `identity_ids` is not allowed
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`; but `type` is currently %q", managedServiceIdentity.Type)
	}

	return &managedServiceIdentity, nil
}

func expandCognitiveAccountAPIProperties(d *pluginsdk.ResourceData) (*cognitiveservices.APIProperties, error) {
	props := cognitiveservices.APIProperties{}
	kind := d.Get("kind")
	if kind == "QnAMaker" {
		if v, ok := d.GetOk("qna_runtime_endpoint"); ok && v != "" {
			props.QnaRuntimeEndpoint = utils.String(v.(string))
		} else {
			return nil, fmt.Errorf("the QnAMaker runtime endpoint `qna_runtime_endpoint` is required when kind is set to `QnAMaker`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_aad_client_id"); ok {
		if kind == "MetricsAdvisor" {
			props.AadClientID = utils.String(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_aad_client_id can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_aad_tenant_id"); ok {
		if kind == "MetricsAdvisor" {
			props.AadTenantID = utils.String(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_aad_tenant_id can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_super_user_name"); ok {
		if kind == "MetricsAdvisor" {
			props.SuperUser = utils.String(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_super_user_name can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_website_name"); ok {
		if kind == "MetricsAdvisor" {
			props.WebsiteName = utils.String(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_website_name can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	return &props, nil
}

func flattenCognitiveAccountNetworkAcls(input *cognitiveservices.NetworkRuleSet) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ipRules := make([]interface{}, 0)
	if input.IPRules != nil {
		for _, v := range *input.IPRules {
			if v.Value == nil {
				continue
			}

			ipRules = append(ipRules, *v.Value)
		}
	}

	virtualNetworkSubnetIds := make([]interface{}, 0)
	virtualNetworkRules := make([]interface{}, 0)
	if input.VirtualNetworkRules != nil {
		for _, v := range *input.VirtualNetworkRules {
			if v.ID == nil {
				continue
			}

			id := *v.ID
			subnetId, err := networkParse.SubnetIDInsensitively(*v.ID)
			if err == nil {
				id = subnetId.ID()
			}

			virtualNetworkSubnetIds = append(virtualNetworkSubnetIds, id)
			virtualNetworkRules = append(virtualNetworkRules, map[string]interface{}{
				"subnet_id":                            id,
				"ignore_missing_vnet_service_endpoint": *v.IgnoreMissingVnetServiceEndpoint,
			})
		}
	}
	return []interface{}{
		map[string]interface{}{
			"default_action":             string(input.DefaultAction),
			"ip_rules":                   pluginsdk.NewSet(pluginsdk.HashString, ipRules),
			"virtual_network_subnet_ids": pluginsdk.NewSet(pluginsdk.HashString, virtualNetworkSubnetIds),
			"virtual_network_rules":      virtualNetworkRules,
		},
	}
}

func flattenCognitiveAccountStorage(input *[]cognitiveservices.UserOwnedStorage) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	results := make([]interface{}, 0)
	for _, v := range *input {
		value := make(map[string]interface{})
		if v.ResourceID != nil {
			value["storage_account_id"] = *v.ResourceID
		}
		if v.IdentityClientID != nil {
			value["identity_client_id"] = *v.IdentityClientID
		}
		results = append(results, value)
	}
	return results
}

func flattenCognitiveAccountIdentity(identity *cognitiveservices.Identity) ([]interface{}, error) {
	if identity == nil || identity.Type == cognitiveservices.ResourceIdentityTypeNone {
		return make([]interface{}, 0), nil
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)

	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}

	if identity.TenantID != nil {
		result["tenant_id"] = *identity.TenantID
	}

	identityIds := make([]interface{}, 0)
	if identity.UserAssignedIdentities != nil {
		for key := range identity.UserAssignedIdentities {
			parsedId, err := msiparse.UserAssignedIdentityID(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
		result["identity_ids"] = pluginsdk.NewSet(pluginsdk.HashString, identityIds)
	}

	return []interface{}{result}, nil
}
