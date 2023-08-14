// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

// @tombuildsstuff: Resources using this Resource ID are going to need a state migration to account for `streamingjobs` -> `streamingJobs` prior to migrating to `hashicorp/go-azure-sdk`
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StreamingJobSchedule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/streamingJob1/schedule/default -rewrite=true
