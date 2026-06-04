// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/shim"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

var _ pollers.PollerType = &DataPlaneTablesPropertiesPoller{}

type DataPlaneTablesPropertiesPoller struct {
	client   shim.StorageTableWrapper
	expected tables.StorageServiceProperties
}

func NewDataPlaneTablesPropertiesPoller(client shim.StorageTableWrapper, expected tables.StorageServiceProperties) *DataPlaneTablesPropertiesPoller {
	return &DataPlaneTablesPropertiesPoller{
		client:   client,
		expected: expected,
	}
}

func (d *DataPlaneTablesPropertiesPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := d.client.GetServiceProperties(ctx)
	if err != nil {
		return nil, pollers.PollingFailedError{
			Message:      fmt.Sprintf("retrieving Table Service Properties: %+v", err),
			HttpResponse: nil,
		}
	}
	if resp == nil {
		return &pollers.PollResult{
			HttpResponse: nil,
			PollInterval: 10 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	if tablePropertiesMatch(*resp, d.expected) {
		return &pollers.PollResult{
			HttpResponse: nil,
			PollInterval: 10 * time.Second,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	return &pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}

func tablePropertiesMatch(actual, expected tables.StorageServiceProperties) bool {
	if expected.Cors != nil {
		if actual.Cors == nil {
			return false
		}
		if len(expected.Cors.CorsRule) == 0 && len(actual.Cors.CorsRule) == 0 {
			// treat nil and empty slice as equivalent
		} else if !reflect.DeepEqual(actual.Cors.CorsRule, expected.Cors.CorsRule) {
			return false
		}
	}

	if expected.Logging != nil {
		if actual.Logging == nil || !reflect.DeepEqual(*actual.Logging, *expected.Logging) {
			return false
		}
	}

	if expected.HourMetrics != nil {
		if actual.HourMetrics == nil || !reflect.DeepEqual(*actual.HourMetrics, *expected.HourMetrics) {
			return false
		}
	}

	if expected.MinuteMetrics != nil {
		if actual.MinuteMetrics == nil || !reflect.DeepEqual(*actual.MinuteMetrics, *expected.MinuteMetrics) {
			return false
		}
	}

	return true
}
