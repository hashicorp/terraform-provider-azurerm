// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AppService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AppServiceEnvironment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/hostingEnvironments/hostingEnvironment1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AppServiceSlot -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AppServiceSlotCustomHostnameBinding -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1/hostNameBindings/binding1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Certificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/certificate1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CertificateOrder -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.CertificateRegistration/certificateOrders/order1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CertificateOrderOld -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificateOrders/order1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FunctionApp -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FunctionAppSlot -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HostnameBinding -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/binding1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HybridConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/hybridConnectionNamespaces/hybridConnectionNamespace1/relays/relay1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedCertificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/customhost.contoso.com
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PublicCertificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/publicCertificates/publicCertificate1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SlotVirtualNetworkSwiftConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1/config/virtualNetwork
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StaticSite -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Web/staticSites/my-static-site1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StaticSiteCustomDomain -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Web/staticSites/my-static-site1/customDomains/name.contoso.com
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkSwiftConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/config/virtualNetwork

// @tombuildsstuff: this is going to require a State Migration to account for `serverfarms` -> `serverFarms` prior to migrating to `hashicorp/go-azure-sdk`
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AppServicePlan -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/serverFarms/farm1
