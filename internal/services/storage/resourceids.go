// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StorageAccountDefaultBlob -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/blobServices/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StorageQueueResourceManager -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/queueServices/default/queues/queue1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StorageShareResourceManager -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/fileServices/fileService1/fileshares/share1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StorageTableResourceManager -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/tableService1/tables/table1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StorageAccountManagementPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/managementPolicies/policy1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StorageContainerImmutabilityPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/blobServices/default/containers/container1/immutabilityPolicies/default
