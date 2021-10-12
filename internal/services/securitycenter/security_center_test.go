package securitycenter_test

import (
	"testing"
)

func TestAccSecurityCenter_pricingAndWorkspace(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests
	// due to the workspace tests depending on the current pricing tier
	testCases := map[string]map[string]func(t *testing.T){
		"pricing": {
			"update": testAccSecurityCenterSubscriptionPricing_update,
		},
		"workspace": {
			"basic":          testAccSecurityCenterWorkspace_basic,
			"update":         testAccSecurityCenterWorkspace_update,
			"requiresImport": testAccSecurityCenterWorkspace_requiresImport,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}
