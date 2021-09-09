package recoveryservices

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ReplicationNetworkMapping -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationFabrics/fabric1/replicationNetworks/network1/replicationNetworkMappings/mapping1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ReplicationFabric -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationFabrics/fabric1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ReplicationPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationPolicies/policy1
