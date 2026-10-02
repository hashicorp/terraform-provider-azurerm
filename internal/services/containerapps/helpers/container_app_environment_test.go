// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func workloadProfileTestData(name, profileType string, minimumCount, maximumCount int) map[string]any {
	return map[string]any{
		"name":                  name,
		"workload_profile_type": profileType,
		"minimum_count":         minimumCount,
		"maximum_count":         maximumCount,
	}
}

func implicitConsumptionProfileTestData() map[string]any {
	return workloadProfileTestData(string(WorkloadProfileSkuConsumption), string(WorkloadProfileSkuConsumption), 0, 0)
}

func workloadProfileTestSet(profiles ...map[string]any) *pluginsdk.Set {
	items := make([]any, 0, len(profiles))
	for _, v := range profiles {
		items = append(items, v)
	}

	return pluginsdk.NewSet(pluginsdk.HashResource(WorkloadProfileSchema().Elem.(*pluginsdk.Resource)), items)
}

func TestOnlyImplicitConsumptionProfileDiffers(t *testing.T) {
	d4 := workloadProfileTestData("D4-01", "D4", 0, 3)
	e4 := workloadProfileTestData("E4-01", "E4", 1, 2)

	testCases := []struct {
		name     string
		returned *pluginsdk.Set
		defined  *pluginsdk.Set
		expected bool
	}{
		{
			name:     "no changes besides the implicit consumption profile",
			returned: workloadProfileTestSet(d4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(d4),
			expected: true,
		},
		{
			name:     "no changes besides the implicit consumption profile with multiple profiles",
			returned: workloadProfileTestSet(d4, e4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(d4, e4),
			expected: true,
		},
		{
			name:     "maximum_count changed",
			returned: workloadProfileTestSet(d4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(workloadProfileTestData("D4-01", "D4", 0, 5)),
			expected: false,
		},
		{
			name:     "minimum_count changed",
			returned: workloadProfileTestSet(d4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(workloadProfileTestData("D4-01", "D4", 1, 3)),
			expected: false,
		},
		{
			name:     "workload_profile_type changed",
			returned: workloadProfileTestSet(d4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(workloadProfileTestData("D4-01", "D8", 0, 3)),
			expected: false,
		},
		{
			name:     "profile renamed",
			returned: workloadProfileTestSet(d4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(workloadProfileTestData("D4-02", "D4", 0, 3)),
			expected: false,
		},
		{
			name:     "profile replaced",
			returned: workloadProfileTestSet(d4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(e4),
			expected: false,
		},
		{
			name:     "one of multiple profiles changed",
			returned: workloadProfileTestSet(d4, e4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(d4, workloadProfileTestData("E4-01", "E4", 1, 4)),
			expected: false,
		},
		{
			name:     "consumption profile defined",
			returned: workloadProfileTestSet(d4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(d4, implicitConsumptionProfileTestData()),
			expected: false,
		},
		{
			name:     "profile added",
			returned: workloadProfileTestSet(d4, implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(d4, e4),
			expected: false,
		},
		{
			name:     "additional returned profile is not a consumption profile",
			returned: workloadProfileTestSet(d4, e4),
			defined:  workloadProfileTestSet(d4),
			expected: false,
		},
		{
			name:     "all profiles removed",
			returned: workloadProfileTestSet(implicitConsumptionProfileTestData()),
			defined:  workloadProfileTestSet(),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := OnlyImplicitConsumptionProfileDiffers(testCase.returned, testCase.defined); actual != testCase.expected {
				t.Fatalf("expected %t but got %t", testCase.expected, actual)
			}
		})
	}
}

// TestWorkloadProfileSchemaDiff exercises the `DiffSuppressFunc` through the Plugin SDK's diff, which applies it to
// every nested attribute of the `workload_profile` block.
func TestWorkloadProfileSchemaDiff(t *testing.T) {
	d4 := workloadProfileTestData("D4-01", "D4", 0, 3)

	testCases := []struct {
		name          string
		returned      []any
		defined       []any
		expectChanges bool
	}{
		{
			name:          "implicit consumption profile is suppressed",
			returned:      []any{d4, implicitConsumptionProfileTestData()},
			defined:       []any{d4},
			expectChanges: false,
		},
		{
			name:          "maximum_count change is not suppressed",
			returned:      []any{d4, implicitConsumptionProfileTestData()},
			defined:       []any{workloadProfileTestData("D4-01", "D4", 0, 5)},
			expectChanges: true,
		},
		{
			name:          "minimum_count change is not suppressed",
			returned:      []any{d4, implicitConsumptionProfileTestData()},
			defined:       []any{workloadProfileTestData("D4-01", "D4", 1, 3)},
			expectChanges: true,
		},
		{
			name:          "workload_profile_type change is not suppressed",
			returned:      []any{d4, implicitConsumptionProfileTestData()},
			defined:       []any{workloadProfileTestData("D4-01", "E4", 0, 3)},
			expectChanges: true,
		},
		{
			name:          "no changes without the implicit consumption profile",
			returned:      []any{d4},
			defined:       []any{d4},
			expectChanges: false,
		},
	}

	resource := &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"workload_profile": WorkloadProfileSchema(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			createDiff, err := resource.Diff(ctx, nil, terraform.NewResourceConfigRaw(map[string]any{"workload_profile": testCase.returned}), nil)
			if err != nil {
				t.Fatalf("building initial diff: %v", err)
			}

			attributes, err := createDiff.Apply(nil, resource.CoreConfigSchema())
			if err != nil {
				t.Fatalf("building initial state: %v", err)
			}

			state := &terraform.InstanceState{
				ID:         "test",
				Attributes: attributes,
			}

			diff, err := resource.Diff(ctx, state, terraform.NewResourceConfigRaw(map[string]any{"workload_profile": testCase.defined}), nil)
			if err != nil {
				t.Fatalf("building update diff: %v", err)
			}

			if actual := diff != nil && !diff.Empty(); actual != testCase.expectChanges {
				t.Fatalf("expected changes %t but got %t: %#v", testCase.expectChanges, actual, diff)
			}
		})
	}
}
