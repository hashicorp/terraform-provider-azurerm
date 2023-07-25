// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StorageQueueResourceManagerId struct {
	SubscriptionId     string
	ResourceGroup      string
	StorageAccountName string
	QueueServiceName   string
	QueueName          string
}

func NewStorageQueueResourceManagerID(subscriptionId, resourceGroup, storageAccountName, queueServiceName, queueName string) StorageQueueResourceManagerId {
	return StorageQueueResourceManagerId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		QueueServiceName:   queueServiceName,
		QueueName:          queueName,
	}
}

func (id StorageQueueResourceManagerId) String() string {
	segments := []string{
		fmt.Sprintf("Queue Name %q", id.QueueName),
		fmt.Sprintf("Queue Service Name %q", id.QueueServiceName),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Queue Resource Manager", segmentsStr)
}

func (id StorageQueueResourceManagerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/queueServices/%s/queues/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.QueueServiceName, id.QueueName)
}

// StorageQueueResourceManagerID parses a StorageQueueResourceManager ID into an StorageQueueResourceManagerId struct
func StorageQueueResourceManagerID(input string) (*StorageQueueResourceManagerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an StorageQueueResourceManager ID: %+v", input, err)
	}

	resourceId := StorageQueueResourceManagerId{
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
	if resourceId.QueueServiceName, err = id.PopSegment("queueServices"); err != nil {
		return nil, err
	}
	if resourceId.QueueName, err = id.PopSegment("queues"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
