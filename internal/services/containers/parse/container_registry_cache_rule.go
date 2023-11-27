// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ContainerRegistryCacheRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	CacheRuleName  string
}

func NewContainerRegistryCacheRuleID(subscriptionId, resourceGroup, registryName, cacheRuleName string) ContainerRegistryCacheRuleId {
	return ContainerRegistryCacheRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		CacheRuleName:  cacheRuleName,
	}
}

func (id ContainerRegistryCacheRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Cache Rule Name %q", id.CacheRuleName),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Container Registry Cache Rule", segmentsStr)
}

func (id ContainerRegistryCacheRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/cacheRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.CacheRuleName)
}

// ContainerRegistryCacheRuleID parses a ContainerRegistryCacheRule ID into an ContainerRegistryCacheRuleId struct
func ContainerRegistryCacheRuleID(input string) (*ContainerRegistryCacheRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ContainerRegistryCacheRule ID: %+v", input, err)
	}

	resourceId := ContainerRegistryCacheRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.RegistryName, err = id.PopSegment("registries"); err != nil {
		return nil, err
	}
	if resourceId.CacheRuleName, err = id.PopSegment("cacheRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
