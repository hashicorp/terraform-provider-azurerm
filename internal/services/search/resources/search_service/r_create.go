// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSearchServiceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Search.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := services.NewSearchServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, services.GetOperationOptions{})
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_search_service", id.ID())
		}
	}

	publicNetworkAccess := services.PublicNetworkAccessEnabled
	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		publicNetworkAccess = services.PublicNetworkAccessDisabled
	}

	var apiKeyOnly any = make(map[string]any, 0)
	skuName := services.SkuName(d.Get("sku").(string))
	ipRulesRaw := d.Get("allowed_ips").(*pluginsdk.Set).List()
	hostingMode := services.HostingMode(d.Get("hosting_mode").(string))
	cmkEnforcementEnabled := d.Get("customer_managed_key_enforcement_enabled").(bool)
	localAuthenticationEnabled := d.Get("local_authentication_enabled").(bool)
	authenticationFailureMode := d.Get("authentication_failure_mode").(string)
	networkRuleBypassOptions := services.SearchBypass(d.Get("network_rule_bypass_option").(string))

	semanticSearchSku := services.SearchSemanticSearchDisabled
	if v := d.Get("semantic_search_sku").(string); v != "" {
		semanticSearchSku = services.SearchSemanticSearch(v)
	}

	cmkEnforcement := services.SearchEncryptionWithCmkDisabled
	if cmkEnforcementEnabled {
		cmkEnforcement = services.SearchEncryptionWithCmkEnabled
	}

	// NOTE: hosting mode is only valid if the SKU is 'standard3'
	if skuName != services.SkuNameStandardThree && strings.EqualFold(string(hostingMode), string(services.HostingModeHighDensity)) {
		return fmt.Errorf("'hosting_mode' can only be defined if the 'sku' field is set to the %q SKU, got %q", string(services.SkuNameStandardThree), skuName)
	}

	// NOTE: 'partition_count' values greater than 1 are not valid for 'free' SKU...
	partitionCount := int64(d.Get("partition_count").(int))

	if (skuName == services.SkuNameFree) && partitionCount > 1 {
		return fmt.Errorf("'partition_count' values greater than 1 cannot be set for the %q SKU, got %d)", string(skuName), partitionCount)
	}

	// NOTE: 'partition_count' values greater than 3 are not valid for 'basic' SKU...

	if (skuName == services.SkuNameBasic) && partitionCount > 3 {
		return fmt.Errorf("'partition_count' values greater than 3 cannot be set for the %q SKU, got %d)", string(skuName), partitionCount)
	}

	// NOTE: 'standard3' services with 'hostingMode' set to 'highDensity' the
	// 'partition_count' must be between 1 and 3.
	if skuName == services.SkuNameStandardThree && partitionCount > 3 && strings.EqualFold(string(hostingMode), string(services.HostingModeHighDensity)) {
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

	// API Only Mode (Default) (e.g. localAuthenticationEnabled = true)...
	authenticationOptions := pointer.To(services.DataPlaneAuthOptions{
		ApiKeyOnly: pointer.To(apiKeyOnly),
	})

	if localAuthenticationEnabled && authenticationFailureMode != "" {
		// API & RBAC Mode..
		authenticationOptions = pointer.To(services.DataPlaneAuthOptions{
			AadOrApiKey: pointer.To(services.DataPlaneAadOrApiKeyAuthOption{
				AadAuthFailureMode: pointer.ToEnum[services.AadAuthFailureMode](authenticationFailureMode),
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
				Bypass:  pointer.To(networkRuleBypassOptions),
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
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	// fix for issue #10151, if the identity type is TypeNone do not include it
	// in the create call, only in the update call when 'identity' is removed from the
	// configuration file...
	if expandedIdentity.Type != identity.TypeNone {
		payload.Identity = expandedIdentity
	}

	if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, services.CreateOrUpdateOperationOptions{}, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceSearchServiceRead(d, meta)
}

func expandSearchServiceIPRules(input []any) *[]services.IPRule {
	output := make([]services.IPRule, 0)

	for _, rule := range input {
		if rule != nil {
			output = append(output, services.IPRule{
				Value: pointer.To(rule.(string)),
			})
		}
	}

	return &output
}
