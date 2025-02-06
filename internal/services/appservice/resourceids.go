// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AppServiceEnvironment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/hostingEnvironments/hostingEnvironment1 -rewrite=true

// @tombuildsstuff: this Resource is going to need a State Migration `serverfarms` -> `serverFarms`
// //go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServicePlan -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/serverfarms/farm1 -rewrite=true
