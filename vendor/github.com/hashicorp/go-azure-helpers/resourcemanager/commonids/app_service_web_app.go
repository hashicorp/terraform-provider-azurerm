// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import "fmt"

// NOTE: A Web App is just an App Service instance - but to allow for potential differentiation in the future
// we're wrapping this to provide a unique Type, Parse and Validation functions.

type WebAppId = AppServiceId

// ParseWebAppID parses 'input' into a WebAppId
func ParseWebAppID(input string) (*WebAppId, error) {
	parsed, err := ParseAppServiceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Web App ID: %+v", input, err)
	}

	return &WebAppId{
		SubscriptionId:    parsed.SubscriptionId,
		ResourceGroupName: parsed.ResourceGroupName,
		SiteName:          parsed.SiteName,
	}, nil
}

// ParseWebAppIDInsensitively parses 'input' case-insensitively into a WebAppId
// note: this method should only be used for API response data and not user input
func ParseWebAppIDInsensitively(input string) (*WebAppId, error) {
	parsed, err := ParseAppServiceIDInsensitively(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Web App ID: %+v", input, err)
	}

	return &WebAppId{
		SubscriptionId:    parsed.SubscriptionId,
		ResourceGroupName: parsed.ResourceGroupName,
		SiteName:          parsed.SiteName,
	}, nil
}

// ValidateWebAppID checks that 'input' can be parsed as a Web App Id
func ValidateWebAppID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWebAppID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
