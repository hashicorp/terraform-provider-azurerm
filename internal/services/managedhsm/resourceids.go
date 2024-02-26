// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Key -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/managedHSMs/mhsm1/keys/key1/versions/version1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=KeyVersionless -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/managedHSMs/mhsm1/keys/key1
