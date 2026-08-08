// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"testing"
)

func TestIsConsumptionProfileType(t *testing.T) {
	cases := []struct {
		Input    string
		Expected bool
	}{
		{Input: "Consumption", Expected: true},
		{Input: "consumption", Expected: true},
		{Input: "Consumption-GPU-NC8as-T4", Expected: true},
		{Input: "Consumption-GPU-NC24-A100", Expected: true},
		{Input: "consumption-gpu-nc8as-t4", Expected: true},
		{Input: "D4", Expected: false},
		{Input: "E4", Expected: false},
		{Input: "NC24-A100", Expected: false},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			result := isConsumptionProfileType(tc.Input)
			if result != tc.Expected {
				t.Fatalf("expected %v for %q, got %v", tc.Expected, tc.Input, result)
			}
		})
	}
}

func TestExpandWorkloadProfiles_ConsumptionGPU(t *testing.T) {
	input := []WorkloadProfileModel{
		{
			Name:                "Consumption",
			WorkloadProfileType: "Consumption",
		},
		{
			Name:                "Consump-GPU-T4",
			WorkloadProfileType: "Consumption-GPU-NC8as-T4",
		},
		{
			Name:                "E4-01",
			WorkloadProfileType: "E4",
			MaximumCount:        2,
			MinimumCount:        0,
		},
	}

	result := ExpandWorkloadProfiles(input)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	profiles := *result
	if len(profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(profiles))
	}

	// Consumption profile should not have MinimumCount/MaximumCount
	if profiles[0].MinimumCount != nil {
		t.Errorf("Consumption profile should not have MinimumCount, got %v", *profiles[0].MinimumCount)
	}
	if profiles[0].MaximumCount != nil {
		t.Errorf("Consumption profile should not have MaximumCount, got %v", *profiles[0].MaximumCount)
	}

	// GPU Consumption profile should not have MinimumCount/MaximumCount
	if profiles[1].MinimumCount != nil {
		t.Errorf("GPU Consumption profile should not have MinimumCount, got %v", *profiles[1].MinimumCount)
	}
	if profiles[1].MaximumCount != nil {
		t.Errorf("GPU Consumption profile should not have MaximumCount, got %v", *profiles[1].MaximumCount)
	}

	// Dedicated profile should have MinimumCount/MaximumCount
	if profiles[2].MinimumCount == nil {
		t.Error("E4 profile should have MinimumCount")
	} else if *profiles[2].MinimumCount != 0 {
		t.Errorf("expected MinimumCount 0, got %d", *profiles[2].MinimumCount)
	}
	if profiles[2].MaximumCount == nil {
		t.Error("E4 profile should have MaximumCount")
	} else if *profiles[2].MaximumCount != 2 {
		t.Errorf("expected MaximumCount 2, got %d", *profiles[2].MaximumCount)
	}
}
