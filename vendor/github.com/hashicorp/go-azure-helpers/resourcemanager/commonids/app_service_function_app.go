// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import "fmt"

// NOTE: A Function App is just an App Service instance - but to allow for potential differentiation in the future
// we're wrapping this to provide a unique Type, Parse and Validation functions.

type FunctionAppId = AppServiceId

// ParseFunctionAppID parses 'input' into a FunctionAppId
func ParseFunctionAppID(input string) (*FunctionAppId, error) {
	parsed, err := ParseAppServiceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Function App ID: %+v", input, err)
	}

	return &FunctionAppId{
		SubscriptionId:    parsed.SubscriptionId,
		ResourceGroupName: parsed.ResourceGroupName,
		SiteName:          parsed.SiteName,
	}, nil
}

// ParseFunctionAppIDInsensitively parses 'input' case-insensitively into a FunctionAppId
// note: this method should only be used for API response data and not user input
func ParseFunctionAppIDInsensitively(input string) (*FunctionAppId, error) {
	parsed, err := ParseAppServiceIDInsensitively(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Function App ID: %+v", input, err)
	}

	return &FunctionAppId{
		SubscriptionId:    parsed.SubscriptionId,
		ResourceGroupName: parsed.ResourceGroupName,
		SiteName:          parsed.SiteName,
	}, nil
}

// ValidateFunctionAppID checks that 'input' can be parsed as a Function App ID
func ValidateFunctionAppID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFunctionAppID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
