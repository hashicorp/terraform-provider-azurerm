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

	if propertiesMatch(*resp, d.expected) {
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

func propertiesMatch(actual, expected tables.StorageServiceProperties) bool {
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
		if actual.Logging == nil {
			return false
		}
		if actual.Logging.Version != expected.Logging.Version ||
			actual.Logging.Delete != expected.Logging.Delete ||
			actual.Logging.Read != expected.Logging.Read ||
			actual.Logging.Write != expected.Logging.Write ||
			actual.Logging.RetentionPolicy.Enabled != expected.Logging.RetentionPolicy.Enabled ||
			actual.Logging.RetentionPolicy.Days != expected.Logging.RetentionPolicy.Days {
			return false
		}
	}

	if expected.HourMetrics != nil {
		if actual.HourMetrics == nil {
			return false
		}
		if actual.HourMetrics.Enabled != expected.HourMetrics.Enabled ||
			actual.HourMetrics.Version != expected.HourMetrics.Version ||
			actual.HourMetrics.RetentionPolicy.Enabled != expected.HourMetrics.RetentionPolicy.Enabled ||
			actual.HourMetrics.RetentionPolicy.Days != expected.HourMetrics.RetentionPolicy.Days {
			return false
		}
		if expected.HourMetrics.IncludeAPIs != nil {
			if actual.HourMetrics.IncludeAPIs == nil || *actual.HourMetrics.IncludeAPIs != *expected.HourMetrics.IncludeAPIs {
				return false
			}
		}
	}

	if expected.MinuteMetrics != nil {
		if actual.MinuteMetrics == nil {
			return false
		}
		if actual.MinuteMetrics.Enabled != expected.MinuteMetrics.Enabled ||
			actual.MinuteMetrics.Version != expected.MinuteMetrics.Version ||
			actual.MinuteMetrics.RetentionPolicy.Enabled != expected.MinuteMetrics.RetentionPolicy.Enabled ||
			actual.MinuteMetrics.RetentionPolicy.Days != expected.MinuteMetrics.RetentionPolicy.Days {
			return false
		}
		if expected.MinuteMetrics.IncludeAPIs != nil {
			if actual.MinuteMetrics.IncludeAPIs == nil || *actual.MinuteMetrics.IncludeAPIs != *expected.MinuteMetrics.IncludeAPIs {
				return false
			}
		}
	}

	return true
}
