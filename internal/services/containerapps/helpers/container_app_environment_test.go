// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
)

func TestFlattenWorkloadProfiles(t *testing.T) {
	cases := []struct {
		Name                    string
		Input                   *[]managedenvironments.WorkloadProfile
		UserDeclaredConsumption bool
		Expected                []WorkloadProfileModel
	}{
		{
			Name:     "nil input returns empty",
			Input:    nil,
			Expected: []WorkloadProfileModel{},
		},
		{
			Name:     "empty input returns empty",
			Input:    &[]managedenvironments.WorkloadProfile{},
			Expected: []WorkloadProfileModel{},
		},
		{
			Name: "single Consumption profile is preserved (Consumption-only env)",
			Input: &[]managedenvironments.WorkloadProfile{
				{
					Name:                string(WorkloadProfileSkuConsumption),
					WorkloadProfileType: string(WorkloadProfileSkuConsumption),
				},
			},
			UserDeclaredConsumption: true,
			Expected: []WorkloadProfileModel{
				{
					Name:                string(WorkloadProfileSkuConsumption),
					WorkloadProfileType: string(WorkloadProfileSkuConsumption),
				},
			},
		},
		{
			Name: "implicit Consumption is filtered when user did not declare it",
			Input: &[]managedenvironments.WorkloadProfile{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        pointer.To(int64(0)),
					MaximumCount:        pointer.To(int64(3)),
				},
				{
					Name:                string(WorkloadProfileSkuConsumption),
					WorkloadProfileType: string(WorkloadProfileSkuConsumption),
				},
			},
			UserDeclaredConsumption: false,
			Expected: []WorkloadProfileModel{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        0,
					MaximumCount:        3,
				},
			},
		},
		{
			Name: "Consumption is preserved when user explicitly declared it",
			Input: &[]managedenvironments.WorkloadProfile{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        pointer.To(int64(0)),
					MaximumCount:        pointer.To(int64(3)),
				},
				{
					Name:                string(WorkloadProfileSkuConsumption),
					WorkloadProfileType: string(WorkloadProfileSkuConsumption),
				},
			},
			UserDeclaredConsumption: true,
			Expected: []WorkloadProfileModel{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        0,
					MaximumCount:        3,
				},
				{
					Name:                string(WorkloadProfileSkuConsumption),
					WorkloadProfileType: string(WorkloadProfileSkuConsumption),
				},
			},
		},
		{
			Name: "multiple dedicated profiles + implicit Consumption: only Consumption filtered",
			Input: &[]managedenvironments.WorkloadProfile{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        pointer.To(int64(0)),
					MaximumCount:        pointer.To(int64(3)),
				},
				{
					Name:                "E4-01",
					WorkloadProfileType: string(WorkloadProfileSkuE4),
					MinimumCount:        pointer.To(int64(1)),
					MaximumCount:        pointer.To(int64(2)),
				},
				{
					Name:                string(WorkloadProfileSkuConsumption),
					WorkloadProfileType: string(WorkloadProfileSkuConsumption),
				},
			},
			UserDeclaredConsumption: false,
			Expected: []WorkloadProfileModel{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        0,
					MaximumCount:        3,
				},
				{
					Name:                "E4-01",
					WorkloadProfileType: string(WorkloadProfileSkuE4),
					MinimumCount:        1,
					MaximumCount:        2,
				},
			},
		},
		{
			Name: "Consumption-GPU profile is never filtered (different name from implicit Consumption)",
			Input: &[]managedenvironments.WorkloadProfile{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        pointer.To(int64(0)),
					MaximumCount:        pointer.To(int64(3)),
				},
				{
					Name:                string(WorkloadProfileSkuConsumptionGpuNc24A100),
					WorkloadProfileType: string(WorkloadProfileSkuConsumptionGpuNc24A100),
				},
			},
			UserDeclaredConsumption: false,
			Expected: []WorkloadProfileModel{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        0,
					MaximumCount:        3,
				},
				{
					Name:                string(WorkloadProfileSkuConsumptionGpuNc24A100),
					WorkloadProfileType: string(WorkloadProfileSkuConsumptionGpuNc24A100),
				},
			},
		},
		{
			Name: "Consumption preserved when flag is true even alongside dedicated profiles",
			Input: &[]managedenvironments.WorkloadProfile{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        pointer.To(int64(0)),
					MaximumCount:        pointer.To(int64(3)),
				},
				{
					Name:                string(WorkloadProfileSkuConsumption),
					WorkloadProfileType: string(WorkloadProfileSkuConsumption),
				},
			},
			UserDeclaredConsumption: true,
			Expected: []WorkloadProfileModel{
				{
					Name:                "D4-01",
					WorkloadProfileType: string(WorkloadProfileSkuD4),
					MinimumCount:        0,
					MaximumCount:        3,
				},
				{
					Name:                string(WorkloadProfileSkuConsumption),
					WorkloadProfileType: string(WorkloadProfileSkuConsumption),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := FlattenWorkloadProfiles(tc.Input, tc.UserDeclaredConsumption)
			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Fatalf("expected %#v, got %#v", tc.Expected, actual)
			}
		})
	}
}
