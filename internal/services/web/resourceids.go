// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AppService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CertificateOrderOld -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificateOrders/order1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SlotVirtualNetworkSwiftConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1/config/virtualNetwork
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkSwiftConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/config/virtualNetwork

// @tombuildsstuff: this is going to require a State Migration to account for `serverfarms` -> `serverFarms` prior to migrating to `hashicorp/go-azure-sdk`
