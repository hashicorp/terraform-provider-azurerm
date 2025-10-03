// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"testing"
)

func TestValidateContainerWorkloadProfiles(t *testing.T) {
	cases := []struct {
		Input WorkloadProfileModel
		Valid bool
	}{
		{
			Input: WorkloadProfileModel{
				Name:                "Consumption",
				WorkloadProfileType: "Consumption",
				MaximumCount:        10,
				MinimumCount:        1,
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "Consumption-GPU-NC24-A100",
				WorkloadProfileType: "Consumption-GPU-NC24-A100",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "Consumption-GPU-NC8as-T4",
				WorkloadProfileType: "Consumption-GPU-NC8as-T4",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "D4",
				WorkloadProfileType: "D4",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "D8",
				WorkloadProfileType: "D8",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "D16",
				WorkloadProfileType: "D16",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "D32",
				WorkloadProfileType: "D32",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "E4",
				WorkloadProfileType: "E4",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "E8",
				WorkloadProfileType: "E8",
				MaximumCount:        8,
				MinimumCount:        4,
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "E16",
				WorkloadProfileType: "E16",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "E32",
				WorkloadProfileType: "E32",
				MaximumCount:        10,
				MinimumCount:        3,
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "NC24-A100",
				WorkloadProfileType: "NC24-A100",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "NC48-A100",
				WorkloadProfileType: "NC48-A100",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "NC96-A100",
				WorkloadProfileType: "NC96-A100",
			},
			Valid: true,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "Invalid-Profile",
				WorkloadProfileType: "Invalid-Profile",
			},
			Valid: false,
		},
		{
			Input: WorkloadProfileModel{
				Name:                "D48",
				WorkloadProfileType: "D48",
				MaximumCount:        10,
				MinimumCount:        1,
			},
			Valid: false,
		},	
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Workload Profile %s", tc.Input.WorkloadProfileType)
		err := ValidateContainerWorkloadProfiles(tc.Input)
		valid := err == nil
		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s", tc.Valid, valid, tc.Input.WorkloadProfileType)
		}
	}
}