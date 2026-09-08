// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
)

func TestCurrentThroughputMode(t *testing.T) {
	testData := []struct {
		name     string
		input    cosmosdb.ThroughputSettingsGetResults
		expected ThroughputMode
	}{
		{
			name:     "empty response",
			input:    cosmosdb.ThroughputSettingsGetResults{},
			expected: ThroughputModeNone,
		},
		{
			name: "nil resource",
			input: cosmosdb.ThroughputSettingsGetResults{
				Properties: &cosmosdb.ThroughputSettingsGetProperties{},
			},
			expected: ThroughputModeNone,
		},
		{
			name: "manual throughput",
			input: cosmosdb.ThroughputSettingsGetResults{
				Properties: &cosmosdb.ThroughputSettingsGetProperties{
					Resource: &cosmosdb.ThroughputSettingsGetPropertiesResource{
						Throughput: pointer.To(int64(400)),
					},
				},
			},
			expected: ThroughputModeManual,
		},
		{
			// an autoscale offer reports the RU/s it is currently scaled to alongside the maximum,
			// so `AutoScaleSettings` is what distinguishes the two modes
			name: "autoscale reports both values",
			input: cosmosdb.ThroughputSettingsGetResults{
				Properties: &cosmosdb.ThroughputSettingsGetProperties{
					Resource: &cosmosdb.ThroughputSettingsGetPropertiesResource{
						Throughput: pointer.To(int64(100)),
						AutoScaleSettings: &cosmosdb.AutoscaleSettingsResource{
							MaxThroughput: 1000,
						},
					},
				},
			},
			expected: ThroughputModeAutoscale,
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			if actual := CurrentThroughputMode(v.input); actual != v.expected {
				t.Fatalf("expected %q but got %q", v.expected, actual)
			}
		})
	}
}

func TestThroughputMigrationRequired(t *testing.T) {
	testData := []struct {
		name     string
		current  ThroughputMode
		desired  ThroughputMode
		expected bool
	}{
		{
			name:     "manual to autoscale",
			current:  ThroughputModeManual,
			desired:  ThroughputModeAutoscale,
			expected: true,
		},
		{
			name:     "autoscale to manual",
			current:  ThroughputModeAutoscale,
			desired:  ThroughputModeManual,
			expected: true,
		},
		{
			name:     "manual unchanged",
			current:  ThroughputModeManual,
			desired:  ThroughputModeManual,
			expected: false,
		},
		{
			name:     "autoscale unchanged",
			current:  ThroughputModeAutoscale,
			desired:  ThroughputModeAutoscale,
			expected: false,
		},
		{
			// serverless and shared throughput resources have no offer of their own
			name:     "no current throughput",
			current:  ThroughputModeNone,
			desired:  ThroughputModeAutoscale,
			expected: false,
		},
		{
			name:     "not set in config",
			current:  ThroughputModeManual,
			desired:  ThroughputModeNone,
			expected: false,
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			if actual := ThroughputMigrationRequired(v.current, v.desired); actual != v.expected {
				t.Fatalf("expected %t but got %t", v.expected, actual)
			}
		})
	}
}
