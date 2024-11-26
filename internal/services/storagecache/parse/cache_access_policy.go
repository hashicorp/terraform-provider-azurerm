// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type CacheAccessPolicyId struct {
	SubscriptionId string
	ResourceGroup  string
	CacheName      string
	Name           string
}

func NewCacheAccessPolicyID(subscriptionId, resourceGroup, cacheName, name string) CacheAccessPolicyId {
	return CacheAccessPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		CacheName:      cacheName,
		Name:           name,
	}
}

func (id CacheAccessPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Cache Name %q", id.CacheName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cache Access Policy", segmentsStr)
}

func (id CacheAccessPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageCache/caches/%s/cacheAccessPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CacheName, id.Name)
}

// CacheAccessPolicyID parses a CacheAccessPolicy ID into an CacheAccessPolicyId struct
func CacheAccessPolicyID(input string) (*CacheAccessPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an CacheAccessPolicy ID: %+v", input, err)
	}

	resourceId := CacheAccessPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CacheName, err = id.PopSegment("caches"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("cacheAccessPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
