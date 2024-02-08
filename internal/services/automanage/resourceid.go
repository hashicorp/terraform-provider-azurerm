// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automanage

// leaving the Automanage prefix here to avoid stuttering the property name for now
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AutomanageConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AutomanageConfigurationHCIAssignment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AzureStackHci/clusters/clusterName1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfile1
