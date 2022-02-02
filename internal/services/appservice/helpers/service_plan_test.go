package helpers_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func TestPlanIsConsumption(t *testing.T) {
	input := []struct {
		name          string
		isConsumption bool
	}{
		{
			name:          "",
			isConsumption: false,
		},
		{
			name:          "Y1",
			isConsumption: true,
		},
		{
			name:          "EP1",
			isConsumption: false,
		},
		{
			name:          "S1",
			isConsumption: false,
		},
	}

	for _, v := range input {
		if actual := helpers.PlanIsConsumption(v.name); actual != v.isConsumption {
			t.Fatalf("expected %s to be %t, got %t", v.name, v.isConsumption, actual)
		}
	}
}

func TestPlanIsElastic(t *testing.T) {
	input := []struct {
		name      string
		isElastic bool
	}{
		{
			name:      "",
			isElastic: false,
		},
		{
			name:      "Y1",
			isElastic: false,
		},
		{
			name:      "EP1",
			isElastic: true,
		},
		{
			name:      "S1",
			isElastic: false,
		},
	}

	for _, v := range input {
		if actual := helpers.PlanIsElastic(v.name); actual != v.isElastic {
			t.Fatalf("expected %s to be %t, got %t", v.name, v.isElastic, actual)
		}
	}
}

func TestPlanIsIsolated(t *testing.T) {
	input := []struct {
		name       string
		isIsolated bool
	}{
		{
			name:       "",
			isIsolated: false,
		},
		{
			name:       "Y1",
			isIsolated: false,
		},
		{
			name:       "EP1",
			isIsolated: false,
		},
		{
			name:       "S1",
			isIsolated: false,
		},
		{
			name:       "I1",
			isIsolated: true,
		},
		{
			name:       "I1v2",
			isIsolated: true,
		},
	}

	for _, v := range input {
		if actual := helpers.PlanIsIsolated(v.name); actual != v.isIsolated {
			t.Fatalf("expected %s to be %t, got %t", v.name, v.isIsolated, actual)
		}
	}
}

func TestPlanIsAppPlan(t *testing.T) {
	input := []struct {
		name      string
		isAppPlan bool
	}{
		{
			name:      "",
			isAppPlan: false,
		},
		{
			name:      "Y1",
			isAppPlan: false,
		},
		{
			name:      "EP1",
			isAppPlan: false,
		},
		{
			name:      "B1",
			isAppPlan: true,
		},
		{
			name:      "S1",
			isAppPlan: true,
		},
		{
			name:      "P1v3",
			isAppPlan: true,
		},
		{
			name:      "I1",
			isAppPlan: false,
		},
		{
			name:      "I1v2",
			isAppPlan: false,
		},
	}

	for _, v := range input {
		if actual := helpers.PlanIsAppPlan(v.name); actual != v.isAppPlan {
			t.Fatalf("expected %s to be %t, got %t", v.name, v.isAppPlan, actual)
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
			expected: "app",
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
