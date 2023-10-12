// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=WebApp -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=WebAppSlot -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FunctionApp -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FunctionAppSlot -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AppServiceEnvironment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/hostingEnvironments/hostingEnvironment1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FunctionAppFunction -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/functions/function1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AppHybridConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/hybridConnectionNamespaces/hybridConnectionNamespace1/relays/relay1

// @tombuildsstuff: this Resource is going to need a State Migration `serverfarms` -> `serverFarms`
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServicePlan -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/serverfarms/farm1 -rewrite=true
