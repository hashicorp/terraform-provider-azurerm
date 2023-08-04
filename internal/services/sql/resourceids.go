// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AzureActiveDirectoryAdministrator -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/administrators/activeDirectory -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Database -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/databases/database1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ElasticPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/elasticPools/elasticPool1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FailoverGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/failoverGroups/failoverGroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/firewallRules/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedInstance -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedDatabase -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1/databases/database1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedInstanceAzureActiveDirectoryAdministrator -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1/administrators/activeDirectory
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Server -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/virtualNetworkRules/virtualNetworkRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=InstanceFailoverGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/locations/Location/instanceFailoverGroups/failoverGroup1
