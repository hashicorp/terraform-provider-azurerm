// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StreamingJobScheduleId struct {
	SubscriptionId   string
	ResourceGroup    string
	StreamingJobName string
	ScheduleName     string
}

func NewStreamingJobScheduleID(subscriptionId, resourceGroup, streamingJobName, scheduleName string) StreamingJobScheduleId {
	return StreamingJobScheduleId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		StreamingJobName: streamingJobName,
		ScheduleName:     scheduleName,
	}
}

func (id StreamingJobScheduleId) String() string {
	segments := []string{
		fmt.Sprintf("Schedule Name %q", id.ScheduleName),
		fmt.Sprintf("Streaming Job Name %q", id.StreamingJobName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Streaming Job Schedule", segmentsStr)
}

func (id StreamingJobScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StreamAnalytics/streamingJobs/%s/schedule/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StreamingJobName, id.ScheduleName)
}

// StreamingJobScheduleID parses a StreamingJobSchedule ID into an StreamingJobScheduleId struct
func StreamingJobScheduleID(input string) (*StreamingJobScheduleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an StreamingJobSchedule ID: %+v", input, err)
	}

	resourceId := StreamingJobScheduleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StreamingJobName, err = id.PopSegment("streamingJobs"); err != nil {
		return nil, err
	}
	if resourceId.ScheduleName, err = id.PopSegment("schedule"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// StreamingJobScheduleIDInsensitively parses an StreamingJobSchedule ID into an StreamingJobScheduleId struct, insensitively
// This should only be used to parse an ID for rewriting, the StreamingJobScheduleID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func StreamingJobScheduleIDInsensitively(input string) (*StreamingJobScheduleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StreamingJobScheduleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'streamingJobs' segment
	streamingJobsKey := "streamingJobs"
	for key := range id.Path {
		if strings.EqualFold(key, streamingJobsKey) {
			streamingJobsKey = key
			break
		}
	}
	if resourceId.StreamingJobName, err = id.PopSegment(streamingJobsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'schedule' segment
	scheduleKey := "schedule"
	for key := range id.Path {
		if strings.EqualFold(key, scheduleKey) {
			scheduleKey = key
			break
		}
	}
	if resourceId.ScheduleName, err = id.PopSegment(scheduleKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
