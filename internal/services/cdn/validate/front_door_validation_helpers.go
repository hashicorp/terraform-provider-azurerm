// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func CdnFrontDoorRouteName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[\da-zA-Z][-\da-zA-Z]{0,88}[\da-zA-Z]$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 2 and 90 characters begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens, got %q`, k, i))
	}

	return nil, nil
}

func CdnFrontDoorCacheDuration(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.HasPrefix(v, "0.") {
		return nil, []error{fmt.Errorf(`%q must not begin with %q if the duration is less than 1 day. If the %q is less than 1 day it should be in the HH:MM:SS format, got %q`, k, "0.", k, v)}
	}

	// fix for issue #22668
	durationParts := strings.Split(v, ".")

	if len(durationParts) > 1 {
		days, err := strconv.Atoi(durationParts[0])
		if err != nil {
			return nil, []error{fmt.Errorf(`%q 'days' segment is invalid, the 'days' segment must be a valid number and have a value that is between 1 and 365, got %q`, k, v)}
		}

		if days > 365 {
			return nil, []error{fmt.Errorf(`%q must be in the d.HH:MM:SS or HH:MM:SS format and must be equal to or lower than %q, got %q`, k, "365.23:59:59", v)}
		}
	}

	// the old regular expersion was broken because it wouldn't allow the value in the tens
	// position to be greater than 6 and the ones position greater than 5
	if m, _ := validate.RegExHelper(i, k, `^([1-9]|([1-9][0-9])|([1-3][0-9][0-9])).((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$|^((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$`); !m {
		return nil, []error{fmt.Errorf(`%q must be in the d.HH:MM:SS or HH:MM:SS format and must be equal to or lower than %q, got %q`, k, "365.23:59:59", v)}
	}

	return nil, nil
}

func CdnFrontDoorUrlPathConditionMatchValue(i interface{}, k string) (_ []string, errors []error) {
	if _, ok := i.(string); !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	return nil, nil
}

func CdnFrontDoorCustomDomainName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^[a-zA-Z0-9][a-zA-Z0-9-]{0,258}[a-zA-Z0-9]$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between 2 and 260 characters in length, must begin with a letter or number, end with a letter or number and contain only letters, numbers and hyphens, got %q`, k, v)}
	}

	return nil, nil
}

func CdnFrontDoorSecretName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^[a-zA-Z0-9][a-zA-Z0-9-]{0,258}[a-zA-Z0-9]$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between 2 and 260 characters in length, must begin with a letter or number, end with a letter or number and contain only letters, numbers and hyphens, got %q`, k, v)}
	}

	return nil, nil
}

func CdnFrontDoorActionsBlock(actions []cdn.BasicDeliveryRuleAction) error {
	routeConfigurationOverride := false
	responseHeader := false
	requestHeader := false
	urlRewrite := false
	urlRedirect := false

	for _, rule := range actions {
		if !routeConfigurationOverride {
			_, routeConfigurationOverride = rule.AsDeliveryRuleRouteConfigurationOverrideAction()
		}

		if !responseHeader {
			_, responseHeader = rule.AsDeliveryRuleResponseHeaderAction()
		}

		if !requestHeader {
			_, requestHeader = rule.AsDeliveryRuleRequestHeaderAction()
		}

		if !urlRewrite {
			_, urlRewrite = rule.AsURLRewriteAction()
		}

		if !urlRedirect {
			_, urlRedirect = rule.AsURLRedirectAction()
		}
	}

	if urlRedirect && urlRewrite {
		return fmt.Errorf("the %q and the %q are both present in the %q match block", "url_redirect_action", "url_rewrite_action", "actions")
	}

	return nil
}

func CdnFrontDoorRuleName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^[a-zA-Z][\da-zA-Z]{0,259}$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between 1 and 260 characters in length, begin with a letter and may contain only letters and numbers, got %q`, k, v)}
	}

	return nil, nil
}

func CdnFrontDoorUrlRedirectActionQueryString(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("%q is invalid: expected type of %q to be string", "url_redirect_action", k)}
	}

	// Query string must be in <key>=<value> format. ? and & will be added automatically so do not include them.
	if v != "" {
		if strings.HasPrefix(v, "?") {
			return nil, []error{fmt.Errorf("'url_redirect_action' is invalid: %q must not start with the '?' character in the 'query_string' field. It will be automatically added by Frontdoor, got %q", k, v)}
		}

		// NOTE: This matches the service code validation logic for this field
		if len(v) > 2048 {
			return nil, []error{fmt.Errorf("'url_redirect_action' is invalid: %q cannot be longer than 2048 characters in length, got %d", k, len(v))}
		}
	}

	return nil, nil
}

func CdnFrontDoorUrlRedirectActionDestinationPath(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("%q is invalid: expected type of %q to be string", "url_redirect_action", k)}
	}

	if v != "" {
		if !strings.HasPrefix(v, "/") {
			return nil, []error{fmt.Errorf("'url_redirect_action' is invalid: %q must begin with a '/', got %q. If you are trying to preserve the incoming path leave the 'destination_path' value empty", k, v)}
		}
	}

	return nil, nil
}
