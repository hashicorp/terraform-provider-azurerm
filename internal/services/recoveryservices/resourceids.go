// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ReplicationNetworkMapping -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationFabrics/fabric1/replicationNetworks/network1/replicationNetworkMappings/mapping1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ReplicationFabric -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationFabrics/fabric1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ReplicationPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationPolicies/policy1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ReplicationProtectedItem -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationFabrics/fabric1/replicationProtectionContainers/container1/replicationProtectedItems/item1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ReplicationProtectionContainerMappings -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationFabrics/fabric1/replicationProtectionContainers/container1/replicationProtectionContainerMappings/mapping1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ReplicationProtectionContainer -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationFabrics/fabric1/replicationProtectionContainers/container1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ProtectionContainer -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/backupFabrics/fabric1/protectionContainers/container1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BackupPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/backupPolicies/policy1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ProtectedItem -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/backupFabrics/Azure/protectionContainers/container1/protectedItems/protectedItem1
