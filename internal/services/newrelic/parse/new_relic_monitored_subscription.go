// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NewRelicMonitoredSubscriptionId struct {
	SubscriptionId            string
	ResourceGroup             string
	MonitorName               string
	MonitoredSubscriptionName string
}

func NewNewRelicMonitoredSubscriptionID(subscriptionId, resourceGroup, monitorName, monitoredSubscriptionName string) NewRelicMonitoredSubscriptionId {
	return NewRelicMonitoredSubscriptionId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		MonitorName:               monitorName,
		MonitoredSubscriptionName: monitoredSubscriptionName,
	}
}

func (id NewRelicMonitoredSubscriptionId) String() string {
	segments := []string{
		fmt.Sprintf("Monitored Subscription Name %q", id.MonitoredSubscriptionName),
		fmt.Sprintf("Monitor Name %q", id.MonitorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "New Relic Monitored Subscription", segmentsStr)
}

func (id NewRelicMonitoredSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/NewRelic.Observability/monitors/%s/monitoredSubscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MonitorName, id.MonitoredSubscriptionName)
}

// NewRelicMonitoredSubscriptionID parses a NewRelicMonitoredSubscription ID into an NewRelicMonitoredSubscriptionId struct
func NewRelicMonitoredSubscriptionID(input string) (*NewRelicMonitoredSubscriptionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an NewRelicMonitoredSubscription ID: %+v", input, err)
	}

	resourceId := NewRelicMonitoredSubscriptionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MonitorName, err = id.PopSegment("monitors"); err != nil {
		return nil, err
	}
	if resourceId.MonitoredSubscriptionName, err = id.PopSegment("monitoredSubscriptions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
