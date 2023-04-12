package search

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	validateSearch "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSearchService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSearchServiceCreate,
		Read:   resourceSearchServiceRead,
		Update: resourceSearchServiceUpdate,
		Delete: resourceSearchServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := services.ParseSearchServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(services.SkuNameFree),
					string(services.SkuNameBasic),
					string(services.SkuNameStandard),
					string(services.SkuNameStandardTwo),
					string(services.SkuNameStandardThree),
					string(services.SkuNameStorageOptimizedLOne),
					string(services.SkuNameStorageOptimizedLTwo),
				}, false),
			},

			"replica_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 12),
			},

			"partition_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateSearch.PartitionCount,
			},

			"api_access_control": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "api_keys",
							ValidateFunc: validation.StringInSlice([]string{
								"api_keys",
								"role_based_access_control",
								"role_based_access_control_and_api_keys",
							}, false),
						},

						"authentication_failure_mode": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(services.AadAuthFailureModeHTTPFourZeroOneWithBearerChallenge),
								string(services.AadAuthFailureModeHTTPFourZeroThree),
							}, false),
						},
					},
				},
			},

			"hosting_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  services.HostingModeDefault,
				ValidateFunc: validation.StringInSlice([]string{
					string(services.HostingModeDefault),
					string(services.HostingModeHighDensity),
				}, false),
			},

			"cmk_enforcement_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"cmk_enforcement_compliance": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"query_keys": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"allowed_ips": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.Any(
						validate.IPv4Address,
						validate.CIDR,
					),
				},
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceSearchServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := services.NewSearchServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, services.GetOperationOptions{})
	if err != nil && !response.WasNotFound(existing.HttpResponse) {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_search_service", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	publicNetworkAccess := services.PublicNetworkAccessEnabled
	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		publicNetworkAccess = services.PublicNetworkAccessDisabled
	}

	skuName := services.SkuName(d.Get("sku").(string))
	ipRulesRaw := d.Get("allowed_ips").(*pluginsdk.Set).List()
	hostingMode := services.HostingMode(d.Get("hosting_mode").(string))
	cmkEnforcementEnabled := d.Get("cmk_enforcement_enabled").(bool)
	apiAccessControl := d.Get("api_access_control").(*pluginsdk.Set).List()
	localAuthDisabled := false

	cmkEnforcement := services.SearchEncryptionWithCmkDisabled
	if cmkEnforcementEnabled {
		cmkEnforcement = services.SearchEncryptionWithCmkEnabled
	}

	// NOTE: hosting mode is only valid if the SKU is 'standard3'
	if skuName != services.SkuNameStandardThree && hostingMode == services.HostingModeHighDensity {
		return fmt.Errorf("'hosting_mode' can only be defined if the 'sku' field is set to the 'standard3' SKU, got %q", skuName)
	}

	// NOTE: 'partition_count' values greater than 1 are not valid for 'free' or 'basic' SKUs...
	partitionCount := int64(d.Get("partition_count").(int))

	if (skuName == services.SkuNameFree || skuName == services.SkuNameBasic) && partitionCount > 1 {
		return fmt.Errorf("'partition_count' values greater than 1 cannot be set for the %q SKU, got %d)", skuName, partitionCount)
	}

	// NOTE: 'standard3' services with 'hostingMode' set to 'highDensity' the
	// 'partition_count' must be between 1 and 3.
	if skuName == services.SkuNameStandardThree && partitionCount > 3 && hostingMode == services.HostingModeHighDensity {
		return fmt.Errorf("'standard3' SKUs in 'highDensity' mode can have a maximum of 3 partitions, got %d", partitionCount)
	}

	// The number of replicas can be between 1 and 12 for 'standard', 'storage_optimized_l1' and storage_optimized_l2' SKUs
	// or between 1 and 3 for 'basic' SKU. Defaults to 1.
	replicaCount, err := validateSearchServiceReplicaCount(int64(d.Get("replica_count").(int)), skuName)
	if err != nil {
		return err
	}

	searchService := services.SearchService{
		Location: location,
		Sku: pointer.To(services.Sku{
			Name: pointer.To(skuName),
		}),
		Properties: &services.SearchServiceProperties{
			PublicNetworkAccess: pointer.To(publicNetworkAccess),
			NetworkRuleSet: pointer.To(services.NetworkRuleSet{
				IPRules: expandSearchServiceIPRules(ipRulesRaw),
			}),
			EncryptionWithCmk: pointer.To(services.EncryptionWithCmk{
				Enforcement: pointer.To(cmkEnforcement),
			}),
			HostingMode: pointer.To(hostingMode),
			// AuthOptions:      authenticationOptions,
			// DisableLocalAuth: pointer.To(localAuthDisabled),
			PartitionCount: pointer.To(partitionCount),
			ReplicaCount:   pointer.To(replicaCount),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	// fix for issue #10151, if the identity type is TypeNone do not include it
	// in the create call, only in the update call when 'identity' is removed from the
	// configuration file...
	if expandedIdentity.Type != identity.TypeNone {
		searchService.Identity = expandedIdentity
	}

	if len(apiAccessControl) > 0 {
		authenticationOptions, err := expandSearchServiceAuthOptions(apiAccessControl)
		if err != nil {
			return err
		}

		// This means that it is RBAC only
		if authenticationOptions == nil {
			localAuthDisabled = true
		}

		searchService.Properties.AuthOptions = authenticationOptions
		searchService.Properties.DisableLocalAuth = pointer.To(localAuthDisabled)
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, searchService, services.CreateOrUpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSearchServiceRead(d, meta)
}

func resourceSearchServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := services.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, services.GetOperationOptions{})
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if d.HasChange("public_network_access_enabled") {
			publicNetworkAccess := services.PublicNetworkAccessEnabled
			if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
				publicNetworkAccess = services.PublicNetworkAccessDisabled
			}

			model.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
		}

		if d.HasChange("identity") {
			expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			// NOTE: Passing type 'None' on update will remove all identities from the service.
			model.Identity = expandedIdentity
		}

		if d.HasChange("hosting_mode") {
			hostingMode := services.HostingMode(d.Get("hosting_mode").(string))
			if pointer.From(model.Sku.Name) != services.SkuNameStandardThree && hostingMode == services.HostingModeHighDensity {
				return fmt.Errorf("'hosting_mode' can only be set to 'highDensity' if the 'sku' is 'standard3', got %q", pointer.From(model.Sku.Name))
			}

			model.Properties.HostingMode = pointer.To(hostingMode)
		}

		if d.HasChange("cmk_enforcement_enabled") {
			cmkEnforcement := services.SearchEncryptionWithCmkDisabled
			if enabled := d.Get("cmk_enforcement_enabled").(bool); enabled {
				cmkEnforcement = services.SearchEncryptionWithCmkEnabled
			}

			model.Properties.EncryptionWithCmk.Enforcement = pointer.To(cmkEnforcement)
		}

		apiAccessControl := d.Get("api_access_control").(*pluginsdk.Set).List()

		if d.HasChange("api_access_control") {
			authenticationOptions, err := expandSearchServiceAuthOptions(apiAccessControl)
			if err != nil {
				return err
			}

			model.Properties.DisableLocalAuth = pointer.To(false)

			// it looks like from my debugging if this is removed it needs to be a pointer to a []interface{map[]}
			if authenticationOptions == nil {
				// This means that it is RBAC only
				model.Properties.DisableLocalAuth = pointer.To(true)
			}

			model.Properties.AuthOptions = authenticationOptions
		}

		if d.HasChange("replica_count") {
			replicaCount, err := validateSearchServiceReplicaCount(int64(d.Get("replica_count").(int)), pointer.From(model.Sku.Name))
			if err != nil {
				return err
			}

			model.Properties.ReplicaCount = pointer.To(replicaCount)
		}

		if d.HasChange("partition_count") {
			partitionCount := int64(d.Get("partition_count").(int))
			// NOTE: 'partition_count' values greater than 1 are not valid for 'free' or 'basic' SKUs...
			if (pointer.From(model.Sku.Name) == services.SkuNameFree || pointer.From(model.Sku.Name) == services.SkuNameBasic) && partitionCount > 1 {
				return fmt.Errorf("'partition_count' values greater than 1 cannot be set for the %q SKU, got %d)", pointer.From(model.Sku.Name), partitionCount)
			}

			// NOTE: If SKU is 'standard3' and the 'hosting_mode' is set to 'highDensity' the maximum number of partitions allowed is 3
			// where if 'hosting_mode' is set to 'default' the maximum number of partitions is 12...
			if pointer.From(model.Sku.Name) == services.SkuNameStandardThree && partitionCount > 3 && pointer.From(model.Properties.HostingMode) == services.HostingModeHighDensity {
				return fmt.Errorf("'standard3' SKUs in 'highDensity' mode can have a maximum of 3 partitions, got %d", partitionCount)
			}

			model.Properties.PartitionCount = pointer.To(partitionCount)
		}

		if d.HasChange("allowed_ips") {
			ipRulesRaw := d.Get("allowed_ips").(*pluginsdk.Set).List()
			model.Properties.NetworkRuleSet.IPRules = expandSearchServiceIPRules(ipRulesRaw)
		}

		if d.HasChange("tags") {
			model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
		}

		err = client.CreateOrUpdateThenPoll(ctx, pointer.From(id), pointer.From(model), services.CreateOrUpdateOperationOptions{})
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		// If you remove this value from your config you need to set this back to the 'default' value in Azure(e.g. 'api_keys')
		// in the expand func, but you also need to remove this code block from your state file else you will get a diff...
		if len(apiAccessControl) == 0 {
			d.Set("api_access_control", apiAccessControl)
		}

		return resourceSearchServiceRead(d, meta)
	}

	return nil
}

func resourceSearchServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := services.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, services.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SearchServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		skuName := ""
		if sku := model.Sku; sku != nil && sku.Name != nil {
			skuName = string(*sku.Name)
		}
		d.Set("sku", skuName)

		if props := model.Properties; props != nil {
			partitionCount := 1         // Default
			replicaCount := 1           // Default
			publicNetworkAccess := true // publicNetworkAccess defaults to true...
			cmkEnforcement := false     // cmkEnforcment defaults to false...
			hostingMode := services.HostingModeDefault

			if count := props.PartitionCount; count != nil {
				partitionCount = int(pointer.From(count))
			}

			if count := props.ReplicaCount; count != nil {
				replicaCount = int(pointer.From(count))
			}

			// NOTE: There is a bug in the API where it returns the PublicNetworkAccess value
			// as 'Disabled' instead of 'disabled'
			if props.PublicNetworkAccess != nil {
				publicNetworkAccess = strings.EqualFold(string(pointer.From(props.PublicNetworkAccess)), string(services.PublicNetworkAccessEnabled))
			}

			if props.HostingMode != nil {
				hostingMode = *props.HostingMode
			}

			var cmkCompliance string
			if props.EncryptionWithCmk != nil {
				cmkEnforcement = strings.EqualFold(string(pointer.From(props.EncryptionWithCmk.Enforcement)), string(services.SearchEncryptionWithCmkEnabled))
				cmkCompliance = string(pointer.From(props.EncryptionWithCmk.EncryptionComplianceStatus))
			}

			// If the 'authenticationOptions' are not in the config file
			// we need to remove them from the state file...
			authenticationOptions := make([]interface{}, 0)

			// **************************************************************************************************************************************
			// I cannot trust props.DisableLocalAuth here, since the service will automatically set a value if it was omitted from the create call...
			// **************************************************************************************************************************************

			// I am using 'DisableLocalAuth' here because when you set your 'api_access_control'
			// to 'RBAC only', the 'props.AuthOptions' will be 'nil'...

			if props.DisableLocalAuth != nil {
				// since I cannot trust the values coming back from Azure, due to their default values, I will have to pull the values from state...
				o, n := d.GetChange("api_access_control")
				newValue := n.(*pluginsdk.Set).List()
				oldValue := o.(*pluginsdk.Set).List()

				if len(newValue) > 0 {
					authenticationOptions = flattenSearchServiceDataPlaneAuthOptions(props.AuthOptions, props.DisableLocalAuth)
				}

				log.Println("************************************************************************")
				log.Println("resourceSearchServiceRead:")
				log.Println("************************************************************************")
				log.Printf("  'api_access_control' New Value         : %+v\n", newValue)
				log.Printf("  'api_access_control' Old Value         : %+v\n", oldValue)
				log.Printf("  Flattened 'authenticationOptions' value: %+v\n", authenticationOptions)
				log.Printf("  Azure RP value                         : %+v\n", props.AuthOptions)
				log.Printf("  Azure RP 'ApiKeyOnly' value            : %+v\n", pointer.From(props.AuthOptions.ApiKeyOnly))
				log.Println("************************************************************************")
			}

			// if props.DisableLocalAuth != nil {
			// 	o, _ := d.GetChange("api_access_control")
			// 	oldValue := o.([]interface{})

			// 	log.Printf("\n\n\n\n\n\n\n\n\n\n******************************\nOld Value: %+v\n******************************\n\n\n\n\n\n\n\n", oldValue)

			// 	// Only call flatten if the value exists in state...
			// 	if len(oldValue) > 0 {
			// 		authenticationOptions = flattenSearchServiceDataPlaneAadOrApiKeyAuthOption(props.AuthOptions, props.DisableLocalAuth)
			// 	}
			// }

			d.Set("api_access_control", authenticationOptions)
			d.Set("partition_count", partitionCount)
			d.Set("replica_count", replicaCount)
			d.Set("public_network_access_enabled", publicNetworkAccess)
			d.Set("hosting_mode", hostingMode)
			d.Set("cmk_enforcement_enabled", cmkEnforcement)
			d.Set("cmk_enforcement_compliance", cmkCompliance)
			d.Set("allowed_ips", flattenSearchServiceIPRules(props.NetworkRuleSet))
		}

		if err = d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	adminKeysClient := meta.(*clients.Client).Search.AdminKeysClient
	adminKeysId, err := adminkeys.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	adminKeysResp, err := adminKeysClient.Get(ctx, *adminKeysId, adminkeys.GetOperationOptions{})
	if err == nil {
		if model := adminKeysResp.Model; model != nil {
			d.Set("primary_key", model.PrimaryKey)
			d.Set("secondary_key", model.SecondaryKey)
		}
	}

	queryKeysClient := meta.(*clients.Client).Search.QueryKeysClient
	queryKeysId, err := querykeys.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}
	queryKeysResp, err := queryKeysClient.ListBySearchService(ctx, *queryKeysId, querykeys.ListBySearchServiceOperationOptions{})
	if err == nil {
		if model := queryKeysResp.Model; model != nil {
			d.Set("query_keys", flattenSearchQueryKeys(*model))
		}
	}

	return nil
}

func resourceSearchServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := services.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id, services.DeleteOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func flattenSearchQueryKeys(input []querykeys.QueryKey) []interface{} {
	results := make([]interface{}, 0)

	for _, v := range input {
		result := make(map[string]interface{})

		if v.Name != nil {
			result["name"] = *v.Name
		}
		result["key"] = *v.Key

		results = append(results, result)
	}

	return results
}

func expandSearchServiceIPRules(input []interface{}) *[]services.IPRule {
	output := make([]services.IPRule, 0)
	if input == nil {
		return &output
	}

	for _, rule := range input {
		if rule != nil {
			output = append(output, services.IPRule{
				Value: utils.String(rule.(string)),
			})
		}
	}

	return &output
}

func expandSearchServiceAuthOptions(input []interface{}) (*services.DataPlaneAuthOptions, error) {
	var foo interface{}
	apiKeyOnlyDefault := make(map[string]interface{}, 0)
	foo = apiKeyOnlyDefault

	// the default(e.g. 'ApiKeyOnly'), only requires an empty 'DataPlaneAuthOptions.ApiKeyOnly'
	// interface which must be an empty map...
	defaultAuthOptions := pointer.To(services.DataPlaneAuthOptions{
		ApiKeyOnly: pointer.To(foo),
	})

	log.Println("************************************************************************")
	log.Println("expandSearchServiceDataPlaneAadOrApiKeyAuthOption:")
	log.Println("************************************************************************")
	log.Printf("  'apiKeyOnlyDefault'  : %+v\n", apiKeyOnlyDefault)
	log.Printf("  'apiKeyOnly'         : %+v\n", foo)
	log.Printf("  'defaultAuthOptions' : %+v\n", defaultAuthOptions)
	log.Println("************************************************************************")

	if len(input) == 0 {
		return defaultAuthOptions, nil
	}

	apiAccessControl := input[0].(map[string]interface{})
	accessControlType := apiAccessControl["type"].(string)
	authFailureMode := apiAccessControl["authentication_failure_mode"].(string)

	switch accessControlType {
	case "role_based_access_control":
		return nil, nil
	case "api_keys":
		return defaultAuthOptions, nil
	case "role_based_access_control_and_api_keys":
		return pointer.To(services.DataPlaneAuthOptions{
			AadOrApiKey: pointer.To(services.DataPlaneAadOrApiKeyAuthOption{
				AadAuthFailureMode: (*services.AadAuthFailureMode)(pointer.To(authFailureMode)),
			}),
		}), nil
	}

	return defaultAuthOptions, nil
}

func flattenSearchServiceIPRules(input *services.NetworkRuleSet) []interface{} {
	if input == nil || *input.IPRules == nil || len(*input.IPRules) == 0 {
		return nil
	}
	result := make([]interface{}, 0)
	for _, rule := range *input.IPRules {
		result = append(result, rule.Value)
	}
	return result
}

func flattenSearchServiceDataPlaneAuthOptions(input *services.DataPlaneAuthOptions, localAuthenticationDisabled *bool) []interface{} {
	// TODO: Validate what I should be checking here...
	// For RBAC Only DataPlaneAuthOptions will be nil...
	if localAuthenticationDisabled == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	// if 'localAuthenticationDisabled' is 'true' that means the this in RBAC Only mode...
	if pointer.From(localAuthenticationDisabled) {
		// Auth is RBAC Only...
		// "localAuthenticationDisabled": true
		result["type"] = "role_based_access_control"
		result["authentication_failure_mode"] = ""
	} else {
		// Auth can be API Only or RBAC and API...
		if input.AadOrApiKey != nil && input.AadOrApiKey.AadAuthFailureMode != nil {
			result["type"] = "role_based_access_control_and_api_keys"
			result["authentication_failure_mode"] = pointer.From(input.AadOrApiKey.AadAuthFailureMode)
		} else {
			result["type"] = "api_keys"
			result["authentication_failure_mode"] = ""
		}
	}

	return []interface{}{result}
}

func validateSearchServiceReplicaCount(replicaCount int64, skuName services.SkuName) (int64, error) {
	switch skuName {
	case services.SkuNameFree:
		if replicaCount > 1 {
			return 0, fmt.Errorf("'replica_count' cannot be greater than 1 for the %q SKU, got %d", skuName, replicaCount)
		}
	case services.SkuNameBasic:
		if replicaCount > 3 {
			return 0, fmt.Errorf("'replica_count' must be between 1 and 3 for the %q SKU, got %d)", skuName, replicaCount)
		}
	}

	return replicaCount, nil
}
