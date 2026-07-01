package provider

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

// TODO: Remove this entire file and test once the legacy top-level block is removed in v5.0.
func TestProvider_LegacyEnhancedValidation(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("This test verifies legacy 4.x enhanced_validation schema logic and is intentionally skipped in 5.x")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cases := []struct {
		name        string
		setupEnv    func(*testing.T)
		config      map[string]any
		expectError bool
		expect      features.EnhancedValidationFeatures
	}{
		{
			name: "Top-level config only v4",
			config: map[string]any{
				"enhanced_validation": []any{
					map[string]any{
						"locations":          true,
						"resource_providers": true,
					},
				},
			},
			expect: features.EnhancedValidationFeatures{
				Locations:         true,
				ResourceProviders: true,
				PreflightEnabled:  false, // Not available at top-level
				LocationFallback:  nil,
			},
			expectError: false,
		},
		{
			name: "Features block only v4",
			config: map[string]any{
				"features": []any{
					map[string]any{
						"enhanced_validation": []any{
							map[string]any{
								"locations":                   false,
								"resource_providers":          false,
								"preflight_enabled":           false,
								"preflight_location_fallback": "",
							},
						},
					},
				},
			},
			expectError: false,
			expect: features.EnhancedValidationFeatures{
				Locations:         false,
				ResourceProviders: false,
				PreflightEnabled:  false,
				LocationFallback:  nil,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupEnv != nil {
				tc.setupEnv(t)
			}

			if tc.config == nil {
				tc.config = map[string]any{}
			}

			p := AzureProvider()
			d := schema.TestResourceDataRaw(t, p.Schema, tc.config)
			env, _ := environments.FromName("public")
			authConfig := &auth.Credentials{
				Environment: *env,
			}

			_, diags := buildClient(ctx, p, d, authConfig, "")

			if tc.expectError {
				if diags == nil || !diags.HasError() {
					t.Fatalf("Expected configuration to return a diagnostic error, but it succeeded")
				}

				foundConflictErr := false
				for _, err := range diags {
					if strings.Contains(err.Summary, "the `enhanced_validation` block is defined at both") {
						foundConflictErr = true
					}
				}
				if !foundConflictErr {
					t.Fatalf("Expected 'defined at both' error, but got: %v", diags)
				}

				return
			}

			if diags != nil && diags.HasError() {
				hasOtherError := false
				for _, err := range diags {
					if !strings.Contains(err.Summary, "unable to build authorizer") {
						hasOtherError = true
					}
				}
				if hasOtherError {
					t.Fatalf("Unexpected error building client: %v", diags)
				}
			}

			// We can't access client if authorizer failed to build because buildClient returns nil client
			// So for non-error cases we just ensure no unexpected errors occurred.
			// The actual feature flags are verified in `TestAccProvider_enhancedValidation`.
		})
	}
}
