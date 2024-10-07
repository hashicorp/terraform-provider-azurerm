// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

// @tombuildsstuff: these resources are going to need state migrations prior to switching to `hashicorp/go-azure-sdk`
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabasePrincipal -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/Clusters/cluster1/Databases/database1/Role/Viewer/FQN/aaduser=11111111-1111-1111-1111-111111111111;22222222-2222-2222-2222-222222222222
