// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	searchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSearchServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := searchClient.From(meta).ServicesClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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
			cmkEnforcement := false     // cmkEnforcement defaults to false...
			endpoint := ""
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
				publicNetworkAccess = strings.EqualFold(pointer.FromEnum(props.PublicNetworkAccess), string(services.PublicNetworkAccessEnabled))
			}

			if props.HostingMode != nil {
				hostingMode = *props.HostingMode
			}

			if props.EncryptionWithCmk != nil {
				cmkEnforcement = strings.EqualFold(pointer.FromEnum(props.EncryptionWithCmk.Enforcement), string(services.SearchEncryptionWithCmkEnabled))
				d.Set("customer_managed_key_encryption_compliance_status", pointer.FromEnum(props.EncryptionWithCmk.EncryptionComplianceStatus))
			}

			if props.Endpoint != nil {
				endpoint = pointer.From(props.Endpoint)
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
						authFailureMode = pointer.FromEnum(props.AuthOptions.AadOrApiKey.AadAuthFailureMode)
					}
				}
			}

			if props.SemanticSearch != nil && pointer.From(props.SemanticSearch) != services.SearchSemanticSearchDisabled {
				semanticSearchSku = pointer.FromEnum(props.SemanticSearch)
			}

			d.Set("authentication_failure_mode", authFailureMode)
			d.Set("local_authentication_enabled", localAuthEnabled)
			d.Set("partition_count", partitionCount)
			d.Set("replica_count", replicaCount)
			d.Set("public_network_access_enabled", publicNetworkAccess)
			d.Set("hosting_mode", hostingMode)
			d.Set("endpoint", endpoint)
			d.Set("customer_managed_key_enforcement_enabled", cmkEnforcement)
			d.Set("allowed_ips", flattenSearchServiceIPRules(props.NetworkRuleSet))
			d.Set("semantic_search_sku", semanticSearchSku)

			if props.NetworkRuleSet != nil {
				d.Set("network_rule_bypass_option", pointer.FromEnum(props.NetworkRuleSet.Bypass))
			}
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err = d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	adminKeysClient := searchClient.From(meta).AdminKeysClient
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

	queryKeysClient := searchClient.From(meta).QueryKeysClient
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

func flattenSearchQueryKeys(input *[]querykeys.QueryKey) []any {
	results := make([]any, 0)

	if input != nil {
		for _, v := range *input {
			results = append(results, map[string]any{
				"name": pointer.From(v.Name),
				"key":  pointer.From(v.Key),
			})
		}
	}

	return results
}

func flattenSearchServiceIPRules(input *services.NetworkRuleSet) []any {
	result := make([]any, 0)
	if input != nil && input.IPRules != nil {
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

func validateSearchServiceApiAccessControlRbac(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
	auth := diff.Get("local_authentication_enabled").(bool)
	failureMode := diff.Get("authentication_failure_mode").(string)

	if !auth && failureMode != "" {
		return errors.New("'authentication_failure_mode' cannot be defined if 'local_authentication_enabled' has been set to 'false'")
	}

	return nil
}
