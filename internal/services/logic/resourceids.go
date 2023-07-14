// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

// @tombuildsstuff: @mbfrahry is going to send a PR to remove/clean these up

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogicAppStandard -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Action -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Logic/workflows/workflow1/actions/action1
