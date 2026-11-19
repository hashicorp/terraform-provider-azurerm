// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func frontDoorCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
	if helpers.IsFrontDoorFullyRetired() {
		return fmt.Errorf("%s", helpers.FullyRetiredMessage)
	}

	// New resources are not supported, and since these fields are 'ForceNew' we also need to block changing them as
	// the re-create would fail with the create error from the service API...
	if helpers.IsFrontDoorDeprecatedForCreation() && d.HasChanges("name", "resource_group_name") {
		return fmt.Errorf("%s", helpers.CreateDeprecationMessage)
	}

	if err := frontDoorSettings(d); err != nil {
		return fmt.Errorf("validating Front Door %q (Resource Group %q): %+v", d.Get("name").(string), d.Get("resource_group_name").(string), err)
	}

	return nil
}

func frontDoorSettings(d *pluginsdk.ResourceDiff) error {
	routingRules := d.Get("routing_rule").([]any)
	configFrontendEndpoints := d.Get("frontend_endpoint").([]any)
	backendPools := d.Get("backend_pool").([]any)
	loadBalancingSettings := d.Get("backend_pool_load_balancing").([]any)
	healthProbeSettings := d.Get("backend_pool_health_probe").([]any)

	if len(configFrontendEndpoints) == 0 {
		return errors.New(`"frontend_endpoint": must have at least one "frontend_endpoint" defined, found 0`)
	}

	// Loop over all of the Routing Rules and validate that only one type of configuration is defined per Routing Rule
	for _, rr := range routingRules {
		routingRule := rr.(map[string]any)
		routingRuleName := routingRule["name"]
		redirectConfig := routingRule["redirect_configuration"].([]any)
		forwardConfig := routingRule["forwarding_configuration"].([]any)

		// Check 0. validate that at least one routing configuration exists per routing rule
		if len(redirectConfig) == 0 && len(forwardConfig) == 0 {
			return fmt.Errorf(`routing_rule %s block is invalid. you must have either a "redirect_configuration" or a "forwarding_configuration" defined for the routing_rule %s`, routingRuleName, routingRuleName)
		}

		// Check 1. validate that only one configuration type is defined per routing rule
		if len(redirectConfig) == 1 && len(forwardConfig) == 1 {
			return fmt.Errorf(`routing_rule %s block is invalid. "redirect_configuration" conflicts with "forwarding_configuration". You can only have one configuration type per each routing rule`, routingRuleName)
		}

		// Check 2. routing rule is a forwarding_configuration type make sure the backend_pool_name exists in the configuration file
		if len(forwardConfig) > 0 {
			fc := forwardConfig[0].(map[string]any)

			if err := verifyBackendPoolExists(fc["backend_pool_name"].(string), backendPools); err != nil {
				return fmt.Errorf(`routing_rule %s is invalid. %+v`, routingRuleName, err)
			}

			// cacheConfiguration validation
			cacheEnabled := fc["cache_enabled"].(bool)
			cacheQueryParameters := fc["cache_query_parameters"].([]any)

			// set cacheQueryParameters to nil if empty
			if len(cacheQueryParameters) < 1 {
				cacheQueryParameters = nil
			}

			cacheDuration := fc["cache_duration"].(string)
			cacheQueryParameterStripDirective := fc["cache_query_parameter_strip_directive"].(string)

			// cacheQueryParameters cannot be set when cacheEnabled is false
			if !cacheEnabled && cacheQueryParameters != nil {
				return fmt.Errorf(`"cache_query_parameters" (%s) cannot be configured when "cache_enabled" is set to "false"`, cacheQueryParameters)
			}

			// cacheQueryParametersStripDirective cannot be different than StripAll (default) when cacheEnabled is false
			if !cacheEnabled && strings.TrimSpace(cacheQueryParameterStripDirective) != "StripAll" {
				return fmt.Errorf(`"cache_query_parameter_strip_directive" (%s) cannot be configured when "cache_enabled" is set to "false"`, cacheQueryParameterStripDirective)
			}

			// cacheDuration cannot be set when cacheEnabled is false
			if !cacheEnabled && strings.TrimSpace(cacheDuration) != "" {
				return fmt.Errorf(`"cache_duration" (%s) cannot be configured when "cache_enabled" is set to "false"`, cacheDuration)
			}

			// cacheQueryParameters cannot be empty when cacheEnabled is true and cacheQueryParameterStripDirective is other than StripAllExcept or StripOnly
			if cacheEnabled && (strings.TrimSpace(cacheQueryParameterStripDirective) == "StripAllExcept" || strings.TrimSpace(cacheQueryParameterStripDirective) == "StripOnly") && cacheQueryParameters == nil {
				return fmt.Errorf(`"cache_query_parameters" cannot be empty when "cache_query_parameter_strip_directive" (%s) is set to "StripAllExcept" or "StripOnly"`, cacheQueryParameterStripDirective)
			}

			// cacheQueryParameters cannot be set when cacheQueryParameterStripDirective is set to StripNone or StripAll
			if cacheEnabled && (strings.TrimSpace(cacheQueryParameterStripDirective) == "StripNone" || strings.TrimSpace(cacheQueryParameterStripDirective) == "StripAll") && cacheQueryParameters != nil {
				return fmt.Errorf(`"cache_query_parameters" cannot be set when "cache_query_parameter_strip_directive" (%s) is set to "StripNone" or "StripAll"`, cacheQueryParameterStripDirective)
			}
		}

		// Check 3. validate that each routing rule frontend_endpoints are actually defined in the resource schema
		if routingRuleFrontends := routingRule["frontend_endpoints"].([]any); len(routingRuleFrontends) > 0 {
			if err := verifyRoutingRuleFrontendEndpoints(routingRuleFrontends, configFrontendEndpoints); err != nil {
				return fmt.Errorf(`"routing_rule":%q %+v`, routingRuleName, err)
			}
		} else {
			return fmt.Errorf(`"routing_rule": %q must have at least one "frontend_endpoints" defined`, routingRuleName)
		}
	}

	// Verify backend pool load balancing settings and health probe settings are defined in the resource schema
	if err := verifyLoadBalancingAndHealthProbeSettings(backendPools, loadBalancingSettings, healthProbeSettings); err != nil {
		return fmt.Errorf(`%+v`, err)
	}

	return nil
}

func verifyBackendPoolExists(backendPoolName string, backendPools []any) error {
	if backendPoolName == "" {
		return fmt.Errorf(`"backend_pool_name" cannot be empty`)
	}

	for _, bps := range backendPools {
		backendPool := bps.(map[string]any)
		if backendPool["name"].(string) == backendPoolName {
			return nil
		}
	}

	return fmt.Errorf(`unable to locate "backend_pool_name":%q in configuration file`, backendPoolName)
}

func verifyRoutingRuleFrontendEndpoints(routingRuleFrontends []any, configFrontendEndpoints []any) error {
	for _, routingRuleFrontend := range routingRuleFrontends {
		// Get the name of the frontend defined in the routing rule
		routingRulefrontendName := routingRuleFrontend.(string)
		found := false

		// Loop over all of the defined frontend endpoints in the config
		// seeing if we find the routing rule frontend in the list
		for _, configFrontendEndpoint := range configFrontendEndpoints {
			configFrontend := configFrontendEndpoint.(map[string]any)
			configFrontendName := configFrontend["name"]
			if routingRulefrontendName == configFrontendName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf(`"frontend_endpoints":%q was not found in the configuration file. verify you have the "frontend_endpoint":%q defined in the configuration file`, routingRulefrontendName, routingRulefrontendName)
		}
	}

	return nil
}

func verifyLoadBalancingAndHealthProbeSettings(backendPools []any, loadBalancingSettings []any, healthProbeSettings []any) error {
	for _, bps := range backendPools {
		backendPool := bps.(map[string]any)
		backendPoolName := backendPool["name"]
		backendPoolLoadBalancingName := backendPool["load_balancing_name"]
		backendPoolHealthProbeName := backendPool["health_probe_name"]
		found := false

		// Verify backend pool load balancing settings name exists
		if len(loadBalancingSettings) > 0 {
			for _, lbs := range loadBalancingSettings {
				loadBalancing := lbs.(map[string]any)
				loadBalancingName := loadBalancing["name"]

				if loadBalancingName == backendPoolLoadBalancingName {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf(`"backend_pool":%q "load_balancing_name":%q was not found in the configuration file. verify you have the "backend_pool_load_balancing":%q defined in the configuration file`, backendPoolName, backendPoolLoadBalancingName, backendPoolLoadBalancingName)
			}
		}

		found = false

		// Verify health probe settings name exists
		if len(healthProbeSettings) > 0 {
			for _, hps := range healthProbeSettings {
				healthProbe := hps.(map[string]any)
				healthProbeName := healthProbe["name"]

				if healthProbeName == backendPoolHealthProbeName {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf(`"backend_pool":%q "health_probe_name":%q was not found in the configuration file. verify you have the "backend_pool_health_probe":%q defined in the configuration file`, backendPoolName, backendPoolHealthProbeName, backendPoolHealthProbeName)
			}
		}
	}

	return nil
}
