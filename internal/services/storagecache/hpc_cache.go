// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/caches"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/storagetargets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func CacheGetAccessPolicyByName(policies []caches.NfsAccessPolicy, name string) *caches.NfsAccessPolicy {
	for _, policy := range policies {
		if policy.Name == name {
			return &policy
		}
	}
	return nil
}

func CacheDeleteAccessPolicyByName(policies []caches.NfsAccessPolicy, name string) []caches.NfsAccessPolicy {
	var newPolicies []caches.NfsAccessPolicy
	for _, policy := range policies {
		if policy.Name != name {
			newPolicies = append(newPolicies, policy)
		}
	}
	return newPolicies
}

func CacheInsertOrUpdateAccessPolicy(policies []caches.NfsAccessPolicy, policy caches.NfsAccessPolicy) ([]caches.NfsAccessPolicy, error) {
	var newPolicies []caches.NfsAccessPolicy

	isNew := true
	for _, existPolicy := range policies {
		if existPolicy.Name == policy.Name {
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

func resourceHPCCacheWaitForCreating(ctx context.Context, client *caches.CachesClient, id caches.CacheId, d *pluginsdk.ResourceData) (caches.GetOperationResponse, error) {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{string(storagetargets.ProvisioningStateTypeCreating)},
		Target:     []string{string(storagetargets.ProvisioningStateTypeSucceeded)},
		Refresh:    resourceHPCCacheRefresh(ctx, client, id),
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	resp, err := state.WaitForStateContext(ctx)
	if err != nil {
		return resp.(caches.GetOperationResponse), fmt.Errorf("waiting for the HPC Cache to be missing (%q): %+v", id.String(), err)
	}

	return resp.(caches.GetOperationResponse), nil
}

func resourceHPCCacheRefresh(ctx context.Context, client *caches.CachesClient, id caches.CacheId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("making Read request on HPC Cache (%q): %+v", id.String(), err)
		}

		if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.ProvisioningState == nil {
			return nil, "Error", fmt.Errorf("unexpected nil pointer")
		}

		return resp, string(*resp.Model.Properties.ProvisioningState), nil
	}
}
