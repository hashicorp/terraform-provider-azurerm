// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedDatabase -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1/databases/database1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedInstance -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedInstanceAzureActiveDirectoryAdministrator -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1/administrators/activeDirectory
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedInstanceEncryptionProtector -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/managedInstances/instance1/encryptionProtector/current
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedInstanceFailoverGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/locations/Location/instanceFailoverGroups/failoverGroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedInstanceVulnerabilityAssessment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1/vulnerabilityAssessments/assessment1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedInstancesSecurityAlertPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/managedInstances/instance1/securityAlertPolicies/Default
