// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	searchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSearchServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := searchClient.From(meta).ServicesClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
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

	if d.HasChange("sku") {
		model.Sku = &services.Sku{
			Name: pointer.ToEnum[services.SkuName](d.Get("sku").(string)),
		}
	}

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

		if pointer.From(model.Sku.Name) != services.SkuNameStandardThree && strings.EqualFold(string(hostingMode), string(services.HostingModeHighDensity)) {
			return fmt.Errorf("'hosting_mode' can only be set to %q if the 'sku' is %q, got %q", services.HostingModeHighDensity, services.SkuNameStandardThree, pointer.From(model.Sku.Name))
		}

		model.Properties.HostingMode = pointer.To(hostingMode)
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
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

		var apiKeyOnly any = make(map[string]any, 0)

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
		// NOTE: 'partition_count' values greater than 1 are not valid for 'free' SKUs...
		if (pointer.From(model.Sku.Name) == services.SkuNameFree) && partitionCount > 1 {
			return fmt.Errorf("'partition_count' values greater than 1 cannot be set for the %q SKU, got %d)", pointer.From(model.Sku.Name), partitionCount)
		}
		// NOTE: 'partition_count' values greater than 3 are not valid for 'basic' SKUs...
		if (pointer.From(model.Sku.Name) == services.SkuNameBasic) && partitionCount > 3 {
			return fmt.Errorf("'partition_count' values greater than 3 cannot be set for the %q SKU, got %d)", pointer.From(model.Sku.Name), partitionCount)
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

		if model.Properties.NetworkRuleSet == nil {
			model.Properties.NetworkRuleSet = &services.NetworkRuleSet{}
		}
		model.Properties.NetworkRuleSet.IPRules = expandSearchServiceIPRules(ipRulesRaw)
	}

	if d.HasChange("network_rule_bypass_option") {
		networkBypassOptions := services.SearchBypass(d.Get("network_rule_bypass_option").(string))
		if model.Properties.NetworkRuleSet == nil {
			model.Properties.NetworkRuleSet = &services.NetworkRuleSet{}
		}
		model.Properties.NetworkRuleSet.Bypass = pointer.To(networkBypassOptions)
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
		model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, model, services.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceSearchServiceRead(d, meta)
}

func validateSearchServiceSKUUpdate(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
	// only validate if the resource already exists
	if diff.Id() == "" {
		return nil
	}

	old, new := diff.GetChange("sku")
	if old == new {
		return nil
	}

	oldSku := old.(string)
	newSku := new.(string)

	// Define SKU hierarchy for validation - excludes Free tier
	skuHierarchy := map[string]int{
		string(services.SkuNameBasic):         1, // basic
		string(services.SkuNameStandard):      2, // standard (S1)
		string(services.SkuNameStandardTwo):   3, // standard2 (S2)
		string(services.SkuNameStandardThree): 4, // standard3 (S3)
		// Free and Storage optimized SKUs are not included as they're not part of the Basic->Standard upgrade path
	}

	_, oldExists := skuHierarchy[oldSku]
	_, newExists := skuHierarchy[newSku]

	// If it's not a valid upgrade (upgrades between basic and standard skus), force recreation instead of blocking the change
	if !oldExists || !newExists {
		return diff.ForceNew("sku")
	}

	return nil
}
