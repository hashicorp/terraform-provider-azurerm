// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-15-preview/scheduledqueryrules"
)

func TestExpandScheduledQueryRulesAlertV2DimensionModel(t *testing.T) {
	cases := []struct {
		name     string
		input    []ScheduledQueryRulesAlertV2DimensionModel
		expected *[]scheduledqueryrules.Dimension
	}{
		{
			name:     "empty input returns nil",
			input:    []ScheduledQueryRulesAlertV2DimensionModel{},
			expected: nil,
		},
		{
			name:     "nil input returns nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "single dimension",
			input: []ScheduledQueryRulesAlertV2DimensionModel{
				{
					Name:     "location",
					Operator: scheduledqueryrules.DimensionOperatorInclude,
					Values:   []string{"eastus2"},
				},
			},
			expected: &[]scheduledqueryrules.Dimension{
				{
					Name:     "location",
					Operator: scheduledqueryrules.DimensionOperatorInclude,
					Values:   []string{"eastus2"},
				},
			},
		},
		{
			name: "multiple dimensions",
			input: []ScheduledQueryRulesAlertV2DimensionModel{
				{
					Name:     "location",
					Operator: scheduledqueryrules.DimensionOperatorInclude,
					Values:   []string{"eastus2"},
				},
				{
					Name:     "quotaName",
					Operator: scheduledqueryrules.DimensionOperatorInclude,
					Values:   []string{"cores"},
				},
			},
			expected: &[]scheduledqueryrules.Dimension{
				{
					Name:     "location",
					Operator: scheduledqueryrules.DimensionOperatorInclude,
					Values:   []string{"eastus2"},
				},
				{
					Name:     "quotaName",
					Operator: scheduledqueryrules.DimensionOperatorInclude,
					Values:   []string{"cores"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := expandScheduledQueryRulesAlertV2DimensionModel(tc.input)

			if tc.expected == nil {
				if result != nil {
					t.Fatalf("expected nil, got %v", result)
				}
				return
			}

			if result == nil {
				t.Fatalf("expected non-nil result, got nil")
			}

			if len(*result) != len(*tc.expected) {
				t.Fatalf("expected %d dimensions, got %d", len(*tc.expected), len(*result))
			}

			for i, expected := range *tc.expected {
				actual := (*result)[i]
				if actual.Name != expected.Name {
					t.Errorf("dimension[%d].Name: expected %s, got %s", i, expected.Name, actual.Name)
				}
				if actual.Operator != expected.Operator {
					t.Errorf("dimension[%d].Operator: expected %s, got %s", i, expected.Operator, actual.Operator)
				}
				if len(actual.Values) != len(expected.Values) {
					t.Errorf("dimension[%d].Values: expected %d values, got %d", i, len(expected.Values), len(actual.Values))
				}
			}
		})
	}
}
