// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Cluster -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NodePool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ContainerRegistryTaskSchedule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/tasks/task1/schedule/schedule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ContainerRegistryTokenPassword -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/tokens/token1/passwords/password
