// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers_test

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func TestPlanIsConsumption(t *testing.T) {
	input := []struct {
		name          *string
		isConsumption bool
	}{
		{
			name:          pointer.To(""),
			isConsumption: false,
		},
		{
			name:          pointer.To("Y1"),
			isConsumption: true,
		},
		{
			name:          pointer.To("EP1"),
			isConsumption: false,
		},
		{
			name:          pointer.To("S1"),
			isConsumption: false,
		},
	}

	for _, v := range input {
		if actual := helpers.PlanIsConsumption(v.name); actual != v.isConsumption {
			t.Fatalf("expected %s to be %t, got %t", *v.name, v.isConsumption, actual)
		}
	}
}

func TestPlanIsElastic(t *testing.T) {
	input := []struct {
		name      *string
		isElastic bool
	}{
		{
			name:      pointer.To(""),
			isElastic: false,
		},
		{
			name:      pointer.To("Y1"),
			isElastic: false,
		},
		{
			name:      pointer.To("EP1"),
			isElastic: true,
		},
		{
			name:      pointer.To("S1"),
			isElastic: false,
		},
	}

	for _, v := range input {
		if actual := helpers.PlanIsElastic(v.name); actual != v.isElastic {
			t.Fatalf("expected %s to be %t, got %t", *v.name, v.isElastic, actual)
		}
	}
}

func TestPlanIsIsolated(t *testing.T) {
	input := []struct {
		name       *string
		isIsolated bool
	}{
		{
			name:       pointer.To(""),
			isIsolated: false,
		},
		{
			name:       pointer.To("Y1"),
			isIsolated: false,
		},
		{
			name:       pointer.To("EP1"),
			isIsolated: false,
		},
		{
			name:       pointer.To("S1"),
			isIsolated: false,
		},
		{
			name:       pointer.To("I1"),
			isIsolated: true,
		},
		{
			name:       pointer.To("I1v2"),
			isIsolated: true,
		},
	}

	for _, v := range input {
		if actual := helpers.PlanIsIsolated(v.name); actual != v.isIsolated {
			t.Fatalf("expected %s to be %t, got %t", *v.name, v.isIsolated, actual)
		}
	}
}

func TestPlanIsAppPlan(t *testing.T) {
	input := []struct {
		name      *string
		isAppPlan bool
	}{
		{
			name:      pointer.To(""),
			isAppPlan: false,
		},
		{
			name:      pointer.To("Y1"),
			isAppPlan: false,
		},
		{
			name:      pointer.To("EP1"),
			isAppPlan: false,
		},
		{
			name:      pointer.To("B1"),
			isAppPlan: true,
		},
		{
			name:      pointer.To("S1"),
			isAppPlan: true,
		},
		{
			name:      pointer.To("P1v3"),
			isAppPlan: true,
		},
		{
			name:      pointer.To("I1"),
			isAppPlan: false,
		},
		{
			name:      pointer.To("I1v2"),
			isAppPlan: false,
		},
	}

	for _, v := range input {
		if actual := helpers.PlanIsAppPlan(v.name); actual != v.isAppPlan {
			t.Fatalf("expected %s to be %t, got %t", *v.name, v.isAppPlan, actual)
		}
	}
}

func TestPlanTypeFromSku(t *testing.T) {
	input := []struct {
		name     string
		expected string
	}{
		{
			name:     "",
			expected: "unknown",
		},
		{
			name:     "Y1",
			expected: "consumption",
		},
		{
			name:     "EP1",
			expected: "elastic",
		},
		{
			name:     "B1",
			expected: "app",
		},
		{
			name:     "S1",
			expected: "app",
		},
		{
			name:     "P1v3",
			expected: "premium",
		},
		{
			name:     "I1",
			expected: "isolated",
		},
		{
			name:     "I1v2",
			expected: "isolated",
		},
	}

	for _, v := range input {
		if actual := helpers.PlanTypeFromSku(v.name); actual != v.expected {
			t.Fatalf("expected %s to be %s, got %s", v.name, v.expected, actual)
		}
	}
}
