package hpccache

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2021-03-01/storagecache"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func cacheGetExistingAccessPolicies(ctx context.Context, client *storagecache.CachesClient, id parse.CacheId) ([]storagecache.NfsAccessPolicy, error) {
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return nil, fmt.Errorf("Error checking for presence of existing HPC Cache %q: %s", id, err)
		}
	}

	prop := existing.CacheProperties
	if prop == nil {
		return nil, nil
	}
	securitySetting := prop.SecuritySettings
	if securitySetting == nil {
		return nil, nil
	}

	accessPolicies := securitySetting.AccessPolicies
	if accessPolicies == nil {
		return nil, nil
	}

	return *accessPolicies, nil
}

func cacheGetAccessPolicyByName(policies []storagecache.NfsAccessPolicy, name string) (storagecache.NfsAccessPolicy, bool) {
	for _, policy := range policies {
		if policy.Name == nil {
			continue
		}
		if *policy.Name == name {
			return policy, true
		}
	}
	return storagecache.NfsAccessPolicy{}, false
}

func cacheInsertOrUpdateAccessPolicy(policies []storagecache.NfsAccessPolicy, policy storagecache.NfsAccessPolicy) ([]storagecache.NfsAccessPolicy, error) {
	if policy.Name == nil {
		return nil, fmt.Errorf("the name of the HPC Cache access policy is nil")
	}
	var newPolicies []storagecache.NfsAccessPolicy

	isNew := true
	for _, existPolicy := range policies {
		if existPolicy.Name != nil && *existPolicy.Name == *policy.Name {
			newPolicies = append(newPolicies, policy)
			isNew = false
			continue
		}
		newPolicies = append(newPolicies, existPolicy)
	}

	if !isNew {
		return newPolicies, nil
	}

	return append(newPolicies, policy), nil
}

func cacheGetAccessPolicyRuleByScope(policyRules []storagecache.NfsAccessRule, scope storagecache.NfsAccessRuleScope) (storagecache.NfsAccessRule, bool) {
	for _, rule := range policyRules {
		if rule.Scope == scope {
			return rule, true
		}
	}

	return storagecache.NfsAccessRule{}, false
}
