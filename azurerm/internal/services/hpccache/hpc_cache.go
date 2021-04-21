package hpccache

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2021-03-01/storagecache"
)

func CacheGetAccessPolicyRuleByScope(policyRules []storagecache.NfsAccessRule, scope storagecache.NfsAccessRuleScope) (storagecache.NfsAccessRule, bool) {
	for _, rule := range policyRules {
		if rule.Scope == scope {
			return rule, true
		}
	}

	return storagecache.NfsAccessRule{}, false
}

func CacheGetAccessPolicyByName(policies []storagecache.NfsAccessPolicy, name string) *storagecache.NfsAccessPolicy {
	for _, policy := range policies {
		if policy.Name != nil && *policy.Name == name {
			return &policy
		}
	}
	return nil
}

func CacheDeleteAccessPolicyByName(policies []storagecache.NfsAccessPolicy, name string) []storagecache.NfsAccessPolicy {
	var newPolicies []storagecache.NfsAccessPolicy
	for _, policy := range policies {
		if policy.Name != nil && *policy.Name != name {
			newPolicies = append(newPolicies, policy)
		}
	}
	return newPolicies
}

func CacheInsertOrUpdateAccessPolicy(policies []storagecache.NfsAccessPolicy, policy storagecache.NfsAccessPolicy) ([]storagecache.NfsAccessPolicy, error) {
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
