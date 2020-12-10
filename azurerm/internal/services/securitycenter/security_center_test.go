package securitycenter_test

import (
	"testing"
)

func TestAccAzureRMSecurityCenter_pricingAndWorkspace(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests
	// due to the workspace tests depending on the current pricing tier
	testCases := map[string]map[string]func(t *testing.T){
		"pricing": {
			"update": testAccAzureRMSecurityCenterSubscriptionPricing_update,
		},
		"workspace": {
			"basic":          testAccAzureRMSecurityCenterWorkspace_basic,
			"update":         testAccAzureRMSecurityCenterWorkspace_update,
			"requiresImport": testAccAzureRMSecurityCenterWorkspace_requiresImport,
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
