package redisenterprise

// leaving the Redisenterprise prefix here to avoid stuttering the property name for now
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RedisEnterpriseCluster -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/cluster1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RedisEnterpriseDatabase -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/cluster1/databases/database1
