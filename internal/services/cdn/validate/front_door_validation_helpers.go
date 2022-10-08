package validate

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
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

	if m, _ := validate.RegExHelper(i, k, `^([1-9]|([1-9][0-9])|([1-3][0-6][0-5])).((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$|^((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between in the d.HH:MM:SS or HH:MM:SS format and must be equal to or lower than %q, got %q`, k, "365.23:59:59", v)}
	}

	if v == "00:00:00" {
		return nil, []error{fmt.Errorf(`%q must be longer than zero seconds, got %q`, k, v)}
	}

	return nil, nil
}

func CdnFrontDoorUrlPathConditionMatchValue(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.HasPrefix(v, "/") {
		return nil, []error{fmt.Errorf(`%q must not begin with the URLs leading slash(e.g. /), got %q`, k, v)}
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
		if strings.ContainsAny(v, "?&") {
			return nil, []error{fmt.Errorf("%q is invalid: %q must not include the %q or the %q characters in the %q field. They will be automatically added by Frontdoor, got %q", "url_redirect_action", k, "?", "&", "query_string", v)}
		}

		if m, _ := validate.RegExHelper(i, k, `^(\b[\da-zA-Z\-\._~]*)(={1})((\b[\da-zA-Z\-\._~]*)|(\{{1}\b(socket_ip|client_ip|client_port|hostname|geo_country|http_method|http_version|query_string|request_scheme|request_uri|ssl_protocol|server_port|url_path){1}\}){1})$`); !m {
			return nil, []error{fmt.Errorf("%q is invalid: %q must be in the <key>=<value> or <key>={action_server_variable} format, got %q", "url_redirect_action", k, v)}
		}
	} else {
		return nil, []error{fmt.Errorf("%q is invalid: %q must not be empty, got %q", "url_redirect_action", k, v)}
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
