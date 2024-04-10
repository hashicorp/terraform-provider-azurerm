// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2023-11-01/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2023-11-01/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2023-11-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: validation.IntInSlice([]int{
					1,
					2,
					3,
					4,
					6,
					12,
				}),
			},

			"local_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"authentication_failure_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(services.AadAuthFailureModeHTTPFourZeroOneWithBearerChallenge),
					string(services.AadAuthFailureModeHTTPFourZeroThree),
				}, false),
			},

			"hosting_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(services.HostingModeDefault),
				ValidateFunc: validation.StringInSlice([]string{
					string(services.HostingModeDefault),
					string(services.HostingModeHighDensity),
				}, false),
			},

			"customer_managed_key_enforcement_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
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

			"semantic_search_sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(services.SearchSemanticSearchFree),
					string(services.SearchSemanticSearchStandard),
				}, false),
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
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := services.NewSearchServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, services.GetOperationOptions{})
	if err != nil && !response.WasNotFound(existing.HttpResponse) {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_search_service", id.ID())
	}

	publicNetworkAccess := services.PublicNetworkAccessEnabled
	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		publicNetworkAccess = services.PublicNetworkAccessDisabled
	}

	var apiKeyOnly interface{} = make(map[string]interface{}, 0)
	skuName := services.SkuName(d.Get("sku").(string))
	ipRulesRaw := d.Get("allowed_ips").(*pluginsdk.Set).List()
	hostingMode := services.HostingMode(d.Get("hosting_mode").(string))
	cmkEnforcementEnabled := d.Get("customer_managed_key_enforcement_enabled").(bool)
	localAuthenticationEnabled := d.Get("local_authentication_enabled").(bool)
	authenticationFailureMode := d.Get("authentication_failure_mode").(string)

	semanticSearchSku := services.SearchSemanticSearchDisabled
	if v := d.Get("semantic_search_sku").(string); v != "" {
		semanticSearchSku = services.SearchSemanticSearch(v)
	}

	cmkEnforcement := services.SearchEncryptionWithCmkDisabled
	if cmkEnforcementEnabled {
		cmkEnforcement = services.SearchEncryptionWithCmkEnabled
	}

	// NOTE: hosting mode is only valid if the SKU is 'standard3'
	if skuName != services.SkuNameStandardThree && hostingMode == services.HostingModeHighDensity {
		return fmt.Errorf("'hosting_mode' can only be defined if the 'sku' field is set to the %q SKU, got %q", string(services.SkuNameStandardThree), skuName)
	}

	// NOTE: 'partition_count' values greater than 1 are not valid for 'free' or 'basic' SKUs...
	partitionCount := int64(d.Get("partition_count").(int))

	if (skuName == services.SkuNameFree || skuName == services.SkuNameBasic) && partitionCount > 1 {
		return fmt.Errorf("'partition_count' values greater than 1 cannot be set for the %q SKU, got %d)", string(skuName), partitionCount)
	}

	// NOTE: 'standard3' services with 'hostingMode' set to 'highDensity' the
	// 'partition_count' must be between 1 and 3.
	if skuName == services.SkuNameStandardThree && partitionCount > 3 && hostingMode == services.HostingModeHighDensity {
		return fmt.Errorf("%q SKUs in %q mode can have a maximum of 3 partitions, got %d", string(services.SkuNameStandardThree), string(services.HostingModeHighDensity), partitionCount)
	}

	// NOTE: Semantic Search SKU cannot be set if the SKU is 'free'
	if skuName == services.SkuNameFree && semanticSearchSku != services.SearchSemanticSearchDisabled {
		return fmt.Errorf("`semantic_search_sku` can only be specified when `sku` is not set to %q", string(services.SkuNameFree))
	}

	// The number of replicas can be between 1 and 12 for 'standard', 'storage_optimized_l1' and storage_optimized_l2' SKUs
	// or between 1 and 3 for 'basic' SKU. Defaults to 1.
	replicaCount, err := validateSearchServiceReplicaCount(int64(d.Get("replica_count").(int)), skuName)
	if err != nil {
		return err
	}

	if !localAuthenticationEnabled && authenticationFailureMode != "" {
		return fmt.Errorf("'authentication_failure_mode' cannot be defined if 'local_authentication_enabled' has been set to 'true'")
	}

	// API Only Mode (Default) (e.g. localAuthenticationEnabled = true)...
	authenticationOptions := pointer.To(services.DataPlaneAuthOptions{
		ApiKeyOnly: pointer.To(apiKeyOnly),
	})

	if localAuthenticationEnabled && authenticationFailureMode != "" {
		// API & RBAC Mode..
		authenticationOptions = pointer.To(services.DataPlaneAuthOptions{
			AadOrApiKey: pointer.To(services.DataPlaneAadOrApiKeyAuthOption{
				AadAuthFailureMode: pointer.To(services.AadAuthFailureMode(authenticationFailureMode)),
			}),
		})
	}

	if !localAuthenticationEnabled {
		// RBAC Only Mode...
		authenticationOptions = nil
	}

	payload := services.SearchService{
		Location: location.Normalize(d.Get("location").(string)),
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
			HostingMode:      pointer.To(hostingMode),
			AuthOptions:      authenticationOptions,
			DisableLocalAuth: pointer.To(!localAuthenticationEnabled),
			PartitionCount:   pointer.To(partitionCount),
			ReplicaCount:     pointer.To(replicaCount),
			SemanticSearch:   pointer.To(semanticSearchSku),
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
		payload.Identity = expandedIdentity
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, payload, services.CreateOrUpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSearchServiceRead(d, meta)
}

func resourceSearchServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := services.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, services.GetOperationOptions{})
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving existing %s: %+v", id, err)
	}

	model := *resp.Model
	if props := model.Properties; props == nil {
		return fmt.Errorf("retrieving existing %s: `properties` was nil", id)
	}

	// The service API has changed where it will not allow the updated model to be
	// passed to the update PATCH call. You must now create a new update payload
	// object by removing all of the READ-ONLY fields from the model...
	// (e.g., privateEndpointConnections, provisioningState, sharedPrivateLinkResources,
	// status and statusDetails)
	model.Properties.PrivateEndpointConnections = nil
	model.Properties.ProvisioningState = nil
	model.Properties.SharedPrivateLinkResources = nil
	model.Properties.Status = nil
	model.Properties.StatusDetails = nil

	if d.HasChange("customer_managed_key_enforcement_enabled") {
		cmkEnforcement := services.SearchEncryptionWithCmkDisabled
		if enabled := d.Get("customer_managed_key_enforcement_enabled").(bool); enabled {
			cmkEnforcement = services.SearchEncryptionWithCmkEnabled
		}

		model.Properties.EncryptionWithCmk = &services.EncryptionWithCmk{
			Enforcement: pointer.To(cmkEnforcement),
		}
	}

	if d.HasChange("hosting_mode") {
		hostingMode := services.HostingMode(d.Get("hosting_mode").(string))
		if model.Sku == nil {
			return fmt.Errorf("updating `hosting_mode` for %s: unable to validate the hosting_mode since `model.Sku` was nil", *id)
		}

		if pointer.From(model.Sku.Name) != services.SkuNameStandardThree && hostingMode == services.HostingModeHighDensity {
			return fmt.Errorf("'hosting_mode' can only be set to %q if the 'sku' is %q, got %q", services.HostingModeHighDensity, services.SkuNameStandardThree, pointer.From(model.Sku.Name))
		}

		model.Properties.HostingMode = pointer.To(hostingMode)
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		model.Identity = expandedIdentity
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := services.PublicNetworkAccessEnabled
		if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
			publicNetworkAccess = services.PublicNetworkAccessDisabled
		}

		model.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChanges("authentication_failure_mode", "local_authentication_enabled") {
		authenticationFailureMode := d.Get("authentication_failure_mode").(string)
		localAuthenticationEnabled := d.Get("local_authentication_enabled").(bool)
		if !localAuthenticationEnabled && authenticationFailureMode != "" {
			return fmt.Errorf("'authentication_failure_mode' cannot be defined if 'local_authentication_enabled' has been set to 'false'")
		}

		var apiKeyOnly interface{} = make(map[string]interface{}, 0)

		// API Only Mode (Default)...
		authenticationOptions := pointer.To(services.DataPlaneAuthOptions{
			ApiKeyOnly: pointer.To(apiKeyOnly),
		})

		if localAuthenticationEnabled && authenticationFailureMode != "" {
			// API & RBAC Mode..
			authenticationOptions = pointer.To(services.DataPlaneAuthOptions{
				AadOrApiKey: pointer.To(services.DataPlaneAadOrApiKeyAuthOption{
					AadAuthFailureMode: (*services.AadAuthFailureMode)(pointer.To(authenticationFailureMode)),
				}),
			})
		}

		if !localAuthenticationEnabled {
			// RBAC Only Mode...
			authenticationOptions = nil
		}

		model.Properties.DisableLocalAuth = pointer.To(!localAuthenticationEnabled)
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
			return fmt.Errorf("%q SKUs in %q mode can have a maximum of 3 partitions, got %d", string(services.SkuNameStandardThree), string(services.HostingModeHighDensity), partitionCount)
		}

		model.Properties.PartitionCount = pointer.To(partitionCount)
	}

	if d.HasChange("allowed_ips") {
		ipRulesRaw := d.Get("allowed_ips").(*pluginsdk.Set).List()

		model.Properties.NetworkRuleSet = &services.NetworkRuleSet{
			IPRules: expandSearchServiceIPRules(ipRulesRaw),
		}
	}

	if d.HasChange("semantic_search_sku") {
		semanticSearchSku := services.SearchSemanticSearchDisabled
		if v := d.Get("semantic_search_sku").(string); v != "" {
			semanticSearchSku = services.SearchSemanticSearch(v)
		}

		// NOTE: Semantic Search SKU cannot be set if the SKU is 'free'
		if pointer.From(model.Sku.Name) == services.SkuNameFree && semanticSearchSku != services.SearchSemanticSearchDisabled {
			return fmt.Errorf("`semantic_search_sku` can only be specified when `sku` is not set to %q", string(services.SkuNameFree))
		}

		model.Properties.SemanticSearch = pointer.To(semanticSearchSku)
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, model, services.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceSearchServiceRead(d, meta)
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
			localAuthEnabled := true
			authFailureMode := ""
			semanticSearchSku := ""

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

			if props.EncryptionWithCmk != nil {
				cmkEnforcement = strings.EqualFold(string(pointer.From(props.EncryptionWithCmk.Enforcement)), string(services.SearchEncryptionWithCmkEnabled))
			}

			// I am using 'DisableLocalAuth' here because when you are in
			// RBAC Only Mode, the 'props.AuthOptions' will be 'nil'...
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !pointer.From(props.DisableLocalAuth)

				// if the AuthOptions are nil that means you are in RBAC Only Mode...
				if props.AuthOptions != nil {
					// If AuthOptions are not nil that means that you are in either
					// API Keys Only Mode or RBAC & API Keys Mode...
					if props.AuthOptions.AadOrApiKey != nil && props.AuthOptions.AadOrApiKey.AadAuthFailureMode != nil {
						// You are in RBAC & API Keys Mode...
						authFailureMode = string(pointer.From(props.AuthOptions.AadOrApiKey.AadAuthFailureMode))
					}
				}
			}

			if props.SemanticSearch != nil && pointer.From(props.SemanticSearch) != services.SearchSemanticSearchDisabled {
				semanticSearchSku = string(pointer.From(props.SemanticSearch))
			}

			d.Set("authentication_failure_mode", authFailureMode)
			d.Set("local_authentication_enabled", localAuthEnabled)
			d.Set("partition_count", partitionCount)
			d.Set("replica_count", replicaCount)
			d.Set("public_network_access_enabled", publicNetworkAccess)
			d.Set("hosting_mode", hostingMode)
			d.Set("customer_managed_key_enforcement_enabled", cmkEnforcement)
			d.Set("allowed_ips", flattenSearchServiceIPRules(props.NetworkRuleSet))
			d.Set("semantic_search_sku", semanticSearchSku)
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
	if err != nil {
		return fmt.Errorf("retrieving Admin Keys for %s: %+v", *id, err)
	}
	if model := adminKeysResp.Model; model != nil {
		d.Set("primary_key", model.PrimaryKey)
		d.Set("secondary_key", model.SecondaryKey)
	}

	queryKeysClient := meta.(*clients.Client).Search.QueryKeysClient
	queryKeysId, err := querykeys.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}
	queryKeysResp, err := queryKeysClient.ListBySearchService(ctx, *queryKeysId, querykeys.ListBySearchServiceOperationOptions{})
	if err != nil {
		return fmt.Errorf("retrieving Query Keys for %s: %+v", *id, err)
	}
	if err := d.Set("query_keys", flattenSearchQueryKeys(queryKeysResp.Model)); err != nil {
		return fmt.Errorf("setting `query_keys`: %+v", err)
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

	if _, err := client.Delete(ctx, *id, services.DeleteOperationOptions{}); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func flattenSearchQueryKeys(input *[]querykeys.QueryKey) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			results = append(results, map[string]interface{}{
				"name": utils.NormalizeNilableString(v.Name),
				"key":  utils.NormalizeNilableString(v.Key),
			})
		}
	}

	return results
}

func expandSearchServiceIPRules(input []interface{}) *[]services.IPRule {
	output := make([]services.IPRule, 0)

	for _, rule := range input {
		if rule != nil {
			output = append(output, services.IPRule{
				Value: utils.String(rule.(string)),
			})
		}
	}

	return &output
}

func flattenSearchServiceIPRules(input *services.NetworkRuleSet) []interface{} {
	result := make([]interface{}, 0)
	if input != nil || input.IPRules != nil {
		for _, rule := range *input.IPRules {
			result = append(result, rule.Value)
		}
	}
	return result
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
