package hpccache

import (
	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2021-03-01/storagecache"
)

func cacheGetDefaultAccessPolicy(policies []storagecache.NfsAccessPolicy) *storagecache.NfsAccessPolicy {
	for _, policy := range policies {
		if policy.Name == nil {
			continue
		}
		if *policy.Name == "default" {
			return &policy
		}
	}
	return nil
}

func cacheGetCustomAccessPolicies(policies []storagecache.NfsAccessPolicy) []storagecache.NfsAccessPolicy {
	var out []storagecache.NfsAccessPolicy
	for _, policy := range policies {
		if policy.Name == nil {
			continue
		}
		if *policy.Name != "default" {
			out = append(out, policy)
		}
	}
	return out
}

func cacheGetAccessPolicyRuleByScope(policyRules []storagecache.NfsAccessRule, scope storagecache.NfsAccessRuleScope) (storagecache.NfsAccessRule, bool) {
	for _, rule := range policyRules {
		if rule.Scope == scope {
			return rule, true
		}
	}

	return storagecache.NfsAccessRule{}, false
}
