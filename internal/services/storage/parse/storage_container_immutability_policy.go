// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StorageContainerImmutabilityPolicyId struct {
	SubscriptionId         string
	ResourceGroup          string
	StorageAccountName     string
	BlobServiceName        string
	ContainerName          string
	ImmutabilityPolicyName string
}

func NewStorageContainerImmutabilityPolicyID(subscriptionId, resourceGroup, storageAccountName, blobServiceName, containerName, immutabilityPolicyName string) StorageContainerImmutabilityPolicyId {
	return StorageContainerImmutabilityPolicyId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		StorageAccountName:     storageAccountName,
		BlobServiceName:        blobServiceName,
		ContainerName:          containerName,
		ImmutabilityPolicyName: immutabilityPolicyName,
	}
}

func (id StorageContainerImmutabilityPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Immutability Policy Name %q", id.ImmutabilityPolicyName),
		fmt.Sprintf("Container Name %q", id.ContainerName),
		fmt.Sprintf("Blob Service Name %q", id.BlobServiceName),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Container Immutability Policy", segmentsStr)
}

func (id StorageContainerImmutabilityPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/%s/containers/%s/immutabilityPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.BlobServiceName, id.ContainerName, id.ImmutabilityPolicyName)
}

// StorageContainerImmutabilityPolicyID parses a StorageContainerImmutabilityPolicy ID into an StorageContainerImmutabilityPolicyId struct
func StorageContainerImmutabilityPolicyID(input string) (*StorageContainerImmutabilityPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an StorageContainerImmutabilityPolicy ID: %+v", input, err)
	}

	resourceId := StorageContainerImmutabilityPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.BlobServiceName, err = id.PopSegment("blobServices"); err != nil {
		return nil, err
	}
	if resourceId.ContainerName, err = id.PopSegment("containers"); err != nil {
		return nil, err
	}
	if resourceId.ImmutabilityPolicyName, err = id.PopSegment("immutabilityPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
