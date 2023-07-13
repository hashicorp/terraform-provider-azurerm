// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SecretId struct {
	SubscriptionId string
	ResourceGroup  string
	VaultName      string
	Name           string
	VersionName    string
}

func NewSecretID(subscriptionId, resourceGroup, vaultName, name, versionName string) SecretId {
	return SecretId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VaultName:      vaultName,
		Name:           name,
		VersionName:    versionName,
	}
}

func (id SecretId) String() string {
	segments := []string{
		fmt.Sprintf("Version Name %q", id.VersionName),
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Secret", segmentsStr)
}

func (id SecretId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/secrets/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.Name, id.VersionName)
}

// SecretID parses a Secret ID into an SecretId struct
func SecretID(input string) (*SecretId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Secret ID: %+v", input, err)
	}

	resourceId := SecretId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VaultName, err = id.PopSegment("vaults"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("secrets"); err != nil {
		return nil, err
	}
	if resourceId.VersionName, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
