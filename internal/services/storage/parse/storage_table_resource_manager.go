// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StorageTableResourceManagerId struct {
	SubscriptionId     string
	ResourceGroup      string
	StorageAccountName string
	TableServiceName   string
	TableName          string
}

func NewStorageTableResourceManagerID(subscriptionId, resourceGroup, storageAccountName, tableServiceName, tableName string) StorageTableResourceManagerId {
	return StorageTableResourceManagerId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		TableServiceName:   tableServiceName,
		TableName:          tableName,
	}
}

func (id StorageTableResourceManagerId) String() string {
	segments := []string{
		fmt.Sprintf("Table Name %q", id.TableName),
		fmt.Sprintf("Table Service Name %q", id.TableServiceName),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Table Resource Manager", segmentsStr)
}

func (id StorageTableResourceManagerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/tableServices/%s/tables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.TableServiceName, id.TableName)
}

// StorageTableResourceManagerID parses a StorageTableResourceManager ID into an StorageTableResourceManagerId struct
func StorageTableResourceManagerID(input string) (*StorageTableResourceManagerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an StorageTableResourceManager ID: %+v", input, err)
	}

	resourceId := StorageTableResourceManagerId{
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
	if resourceId.TableServiceName, err = id.PopSegment("tableServices"); err != nil {
		return nil, err
	}
	if resourceId.TableName, err = id.PopSegment("tables"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
