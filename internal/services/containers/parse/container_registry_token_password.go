// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ContainerRegistryTokenPasswordId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	TokenName      string
	PasswordName   string
}

func NewContainerRegistryTokenPasswordID(subscriptionId, resourceGroup, registryName, tokenName, passwordName string) ContainerRegistryTokenPasswordId {
	return ContainerRegistryTokenPasswordId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		TokenName:      tokenName,
		PasswordName:   passwordName,
	}
}

func (id ContainerRegistryTokenPasswordId) String() string {
	segments := []string{
		fmt.Sprintf("Password Name %q", id.PasswordName),
		fmt.Sprintf("Token Name %q", id.TokenName),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Container Registry Token Password", segmentsStr)
}

func (id ContainerRegistryTokenPasswordId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/tokens/%s/%ss/password"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName, id.PasswordName)
}

// ContainerRegistryTokenPasswordID parses a ContainerRegistryTokenPassword ID into an ContainerRegistryTokenPasswordId struct
func ContainerRegistryTokenPasswordID(input string) (*ContainerRegistryTokenPasswordId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ContainerRegistryTokenPassword ID: %+v", input, err)
	}

	resourceId := ContainerRegistryTokenPasswordId{
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
	if resourceId.TokenName, err = id.PopSegment("tokens"); err != nil {
		return nil, err
	}
	if resourceId.PasswordName, err = id.PopSegment("passwords"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
