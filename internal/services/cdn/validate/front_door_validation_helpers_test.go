// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-09-01/rules"
)

func TestCdnFrontDoorActionsBlock(t *testing.T) {
	tests := []struct {
		name        string
		actions     []rules.DeliveryRuleAction
		errContains string
	}{
		{
			name: "valid actions",
			actions: []rules.DeliveryRuleAction{
				urlRewriteAction(),
				requestHeaderAction(),
				responseHeaderAction(),
				routeConfigurationOverrideAction(),
			},
		},
		{
			name: "rewrite and redirect conflict",
			actions: []rules.DeliveryRuleAction{
				urlRewriteAction(),
				urlRedirectAction(),
			},
			errContains: "both present",
		},
		{
			name: "duplicate rewrite action",
			actions: []rules.DeliveryRuleAction{
				urlRewriteAction(),
				urlRewriteAction(),
			},
			errContains: "url_rewrite_action",
		},
		{
			name: "duplicate redirect action",
			actions: []rules.DeliveryRuleAction{
				urlRedirectAction(),
				urlRedirectAction(),
			},
			errContains: "url_redirect_action",
		},
		{
			name: "duplicate route override action",
			actions: []rules.DeliveryRuleAction{
				routeConfigurationOverrideAction(),
				routeConfigurationOverrideAction(),
			},
			errContains: "route_configuration_override_action",
		},
		{
			name: "too many actions",
			actions: []rules.DeliveryRuleAction{
				requestHeaderAction(),
				requestHeaderAction(),
				requestHeaderAction(),
				responseHeaderAction(),
				responseHeaderAction(),
				routeConfigurationOverrideAction(),
			},
			errContains: "up to 5 match actions",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CdnFrontDoorActionsBlock(test.actions)
			assertFrontDoorValidationHelperError(t, err, test.errContains)
		})
	}
}

func urlRewriteAction() rules.DeliveryRuleAction {
	return rules.URLRewriteAction{
		Name: rules.DeliveryRuleActionNameURLRewrite,
	}
}

func urlRedirectAction() rules.DeliveryRuleAction {
	return rules.URLRedirectAction{
		Name: rules.DeliveryRuleActionNameURLRedirect,
	}
}

func requestHeaderAction() rules.DeliveryRuleAction {
	return rules.DeliveryRuleRequestHeaderAction{
		Name: rules.DeliveryRuleActionNameModifyRequestHeader,
	}
}

func responseHeaderAction() rules.DeliveryRuleAction {
	return rules.DeliveryRuleResponseHeaderAction{
		Name: rules.DeliveryRuleActionNameModifyResponseHeader,
	}
}

func routeConfigurationOverrideAction() rules.DeliveryRuleAction {
	return rules.DeliveryRuleRouteConfigurationOverrideAction{
		Name: rules.DeliveryRuleActionNameRouteConfigurationOverride,
	}
}

func assertFrontDoorValidationHelperError(t *testing.T, err error, errContains string) {
	t.Helper()

	if errContains == "" {
		if err != nil {
			t.Fatalf("expected no error but got %q", err)
		}
		return
	}

	if err == nil {
		t.Fatalf("expected error containing %q but got nil", errContains)
	}

	if !strings.Contains(err.Error(), errContains) {
		t.Fatalf("expected error containing %q but got %q", errContains, err)
	}
}
