// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AttachedDatabaseConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/attachedDatabaseConfigurations/config1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Cluster -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ClusterPrincipalAssignment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/principalAssignments/assignment1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Database -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabasePrincipalAssignment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/principalAssignments/assignment1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DataConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/dataConnections/connection1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Script -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/scripts/script1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedPrivateEndpoints -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/managedPrivateEndpoints/endpoint1 -rewrite=true

// @tombuildsstuff: these resources are going to need state migrations prior to switching to `hashicorp/go-azure-sdk`
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabasePrincipal -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/Clusters/cluster1/Databases/database1/Role/Viewer/FQN/aaduser=11111111-1111-1111-1111-111111111111;22222222-2222-2222-2222-222222222222
