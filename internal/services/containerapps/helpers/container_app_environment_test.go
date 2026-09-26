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
		{Input: "Consumption-GPU-NC8as-T4", Expected: true},
		{Input: "Consumption-GPU-NC24-A100", Expected: true},
		{Input: "consumption", Expected: false},
		{Input: "Consumption-GPU-Unknown", Expected: false},
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

func TestValidateWorkloadProfileCounts(t *testing.T) {
	cases := []struct {
		Name        string
		Input       []WorkloadProfileModel
		ExpectError bool
	}{
		{
			Name: "consumption without counts",
			Input: []WorkloadProfileModel{
				{Name: "Consumption", WorkloadProfileType: "Consumption"},
			},
		},
		{
			Name: "gpu consumption without counts",
			Input: []WorkloadProfileModel{
				{Name: "Consump-GPU-T4", WorkloadProfileType: "Consumption-GPU-NC8as-T4"},
			},
		},
		{
			Name: "dedicated with counts",
			Input: []WorkloadProfileModel{
				{Name: "E4-01", WorkloadProfileType: "E4", MinimumCount: 1, MaximumCount: 3},
			},
		},
		{
			Name: "consumption with minimum_count",
			Input: []WorkloadProfileModel{
				{Name: "Consumption", WorkloadProfileType: "Consumption", MinimumCount: 1},
			},
			ExpectError: true,
		},
		{
			Name: "gpu consumption with maximum_count",
			Input: []WorkloadProfileModel{
				{Name: "E4-01", WorkloadProfileType: "E4", MinimumCount: 1, MaximumCount: 3},
				{Name: "Consump-GPU-A100", WorkloadProfileType: "Consumption-GPU-NC24-A100", MaximumCount: 2},
			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			err := ValidateWorkloadProfileCounts(tc.Input)
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but didn't get one")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("expected no error but got: %+v", err)
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

	for i, v := range input {
		if profiles[i].Name != v.Name {
			t.Errorf("expected Name %q, got %q", v.Name, profiles[i].Name)
		}
		if profiles[i].WorkloadProfileType != v.WorkloadProfileType {
			t.Errorf("expected WorkloadProfileType %q for %q, got %q", v.WorkloadProfileType, v.Name, profiles[i].WorkloadProfileType)
		}
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
