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

type StorageShareResourceManagerId struct {
	SubscriptionId     string
	ResourceGroup      string
	StorageAccountName string
	FileServiceName    string
	FileshareName      string
}

func NewStorageShareResourceManagerID(subscriptionId, resourceGroup, storageAccountName, fileServiceName, fileshareName string) StorageShareResourceManagerId {
	return StorageShareResourceManagerId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		FileServiceName:    fileServiceName,
		FileshareName:      fileshareName,
	}
}

func (id StorageShareResourceManagerId) String() string {
	segments := []string{
		fmt.Sprintf("Fileshare Name %q", id.FileshareName),
		fmt.Sprintf("File Service Name %q", id.FileServiceName),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Share Resource Manager", segmentsStr)
}

func (id StorageShareResourceManagerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/%s/fileshares/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.FileServiceName, id.FileshareName)
}

// StorageShareResourceManagerID parses a StorageShareResourceManager ID into an StorageShareResourceManagerId struct
func StorageShareResourceManagerID(input string) (*StorageShareResourceManagerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an StorageShareResourceManager ID: %+v", input, err)
	}

	resourceId := StorageShareResourceManagerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.FileServiceName, err = id.PopSegment("fileServices"); err != nil {
		return nil, err
	}
	if resourceId.FileshareName, err = id.PopSegment("fileshares"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
