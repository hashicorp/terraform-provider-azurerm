// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StorageAccountDefaultBlobId struct {
	SubscriptionId     string
	ResourceGroup      string
	StorageAccountName string
	BlobServiceName    string
}

func NewStorageAccountDefaultBlobID(subscriptionId, resourceGroup, storageAccountName, blobServiceName string) StorageAccountDefaultBlobId {
	return StorageAccountDefaultBlobId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		BlobServiceName:    blobServiceName,
	}
}

func (id StorageAccountDefaultBlobId) String() string {
	segments := []string{
		fmt.Sprintf("Blob Service Name %q", id.BlobServiceName),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Account Default Blob", segmentsStr)
}

func (id StorageAccountDefaultBlobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.BlobServiceName)
}

// StorageAccountDefaultBlobID parses a StorageAccountDefaultBlob ID into an StorageAccountDefaultBlobId struct
func StorageAccountDefaultBlobID(input string) (*StorageAccountDefaultBlobId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an StorageAccountDefaultBlob ID: %+v", input, err)
	}

	resourceId := StorageAccountDefaultBlobId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
