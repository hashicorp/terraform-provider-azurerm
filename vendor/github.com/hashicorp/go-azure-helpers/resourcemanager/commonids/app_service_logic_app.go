// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import "fmt"

// NOTE: A Logic App is just an App Service instance - but to allow for potential differentiation in the future
// we're wrapping this to provide a unique Type, Parse and Validation functions.

type LogicAppId = AppServiceId

// ParseLogicAppId parses 'input' into a LogicAppId
func ParseLogicAppId(input string) (*LogicAppId, error) {
	parsed, err := ParseAppServiceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Logic App ID: %+v", input, err)
	}

	return &LogicAppId{
		SubscriptionId:    parsed.SubscriptionId,
		ResourceGroupName: parsed.ResourceGroupName,
		SiteName:          parsed.SiteName,
	}, nil
}

// ParseLogicAppIdInsensitively parses 'input' case-insensitively into a LogicAppId
// note: this method should only be used for API response data and not user input
func ParseLogicAppIdInsensitively(input string) (*LogicAppId, error) {
	parsed, err := ParseAppServiceIDInsensitively(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Logic App ID: %+v", input, err)
	}

	return &LogicAppId{
		SubscriptionId:    parsed.SubscriptionId,
		ResourceGroupName: parsed.ResourceGroupName,
		SiteName:          parsed.SiteName,
	}, nil
}

// ValidateLogicAppId checks that 'input' can be parsed as a Logic App ID
func ValidateLogicAppId(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLogicAppId(v); err != nil {
		errors = append(errors, err)
	}

	return
}
