// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package springcloud

// @tombuildsstuff: the following Resource IDs use the incorrect casing and will need State Migrations to fix the casing prior to moving to `hashicorp/go-azure-sdk` (e.g. `Spring` -> `spring`)
// in addition, Resources ending in `/default` shouldn't be exposed as separate resources - and instead should be embedded within the parent resource, as such these Resources will need to be deprecated and instead inlined within the parent Resource
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SpringCloudDevToolPortal -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/Spring/service1/DevToolPortals/default
