// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SystemTopic -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.EventGrid/systemTopics/systemTopic1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Topic -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.EventGrid/topics/topic1
