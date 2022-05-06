package cdn

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-11-01/frontdoor"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO: the validation methods want moving into ./validation
// TODO: the convert methods need to be made private (since they're only related to this package)

func ValidPrivateLinkTargetTypes() []string {
	return []string{"blob", "blob_secondary", "sites", "web"}
}

func ConvertCdnFrontdoorTags(tagMap *map[string]string) map[string]*string {
	t := make(map[string]*string)

	if tagMap != nil {
		for k, v := range *tagMap {
			tagKey := k
			tagValue := v
			t[tagKey] = &tagValue
		}
	}

	return t
}

func ConvertCdnFrontdoorTagsToTagsFlatten(tagMap map[string]*string) *map[string]string {
	t := make(map[string]string)

	for k, v := range tagMap {
		tagKey := k
		tagValue := v
		t[tagKey] = *tagValue
	}

	return &t
}

func convertBoolToEnabledState(isEnabled bool) cdn.EnabledState {
	out := cdn.EnabledStateDisabled

	if isEnabled {
		out = cdn.EnabledStateEnabled
	}

	return out
}

func convertEnabledStateToBool(enabledState *cdn.EnabledState) bool {
	if enabledState == nil {
		return false
	}

	return (*enabledState == cdn.EnabledStateEnabled)
}

func expandResourceReference(input string) *cdn.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &cdn.ResourceReference{
		ID: utils.String(input),
	}
}

// Takes a Slice of strings and transforms it into a CSV formatted string.
func expandStringSliceToCsvFormat(input []interface{}) *string {
	if len(input) == 0 {
		return nil
	}

	v := utils.ExpandStringSlice(input)
	// TODO: this can be reduced to:
	// csv := fmt.Sprintf("[%s]", strings.Join(*v, ","))
	csv := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(*v)), ","), "[]")

	return &csv
}

// Takes a CSV formatted string and transforms it into a Slice of strings.
func flattenCsvToStringSlice(input *string) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	v := strings.Split(*input, ",")

	for _, s := range v {
		results = append(results, s)
	}

	return results
}

func flattenResourceReference(input *cdn.ResourceReference) string {
	result := ""
	if input == nil {
		return result
	}

	if input.ID != nil {
		result = *input.ID
	}

	return result
}

func flattenTransformSlice(input *[]frontdoor.TransformType) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

func flattenFrontendEndpointLinkSlice(input *[]frontdoor.FrontendEndpointLink) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			if item.ID == nil {
				continue
			}

			result = append(result, *item.ID)
		}
	}
	return result
}

func ConvertBoolToRouteHttpsRedirect(isEnabled bool) cdn.HTTPSRedirect {
	out := cdn.HTTPSRedirectDisabled

	if isEnabled {
		out = cdn.HTTPSRedirectEnabled
	}

	return out
}

func ConvertRouteHttpsRedirectToBool(httpsRedirect *cdn.HTTPSRedirect) bool {
	if httpsRedirect == nil {
		return false
	}

	return (*httpsRedirect == cdn.HTTPSRedirectEnabled)
}

func ConvertBoolToRouteLinkToDefaultDomain(isLinked bool) cdn.LinkToDefaultDomain {
	out := cdn.LinkToDefaultDomainDisabled

	if isLinked {
		out = cdn.LinkToDefaultDomainEnabled
	}

	return out
}

func ConvertRouteLinkToDefaultDomainToBool(linkToDefaultDomain *cdn.LinkToDefaultDomain) bool {
	if linkToDefaultDomain == nil {
		return false
	}

	return (*linkToDefaultDomain == cdn.LinkToDefaultDomainEnabled)
}

func IsValidDomain(i interface{}, k string) (warnings []string, errors []error) {
	// TODO: this can be replaced by:
	// Validation.Any(validation.IsIPv6Address, validation.IsIPv4Address, validation.StringIsNotEmpty)

	isIPv6 := true
	isIPv4 := true

	if _, err := validation.IsIPv6Address(i, k); len(err) > 0 {
		isIPv6 = false
	}

	if _, err := validation.IsIPv4Address(i, k); len(err) > 0 {
		isIPv4 = false
	}

	// If it's not an IPv4 or IPv6 check to see if it is a valid Domain name...
	if !isIPv4 && !isIPv6 {
		// TODO: Figure out a better way to validate Domain Name if not an IP Address
		if warn, err := validation.StringIsNotEmpty(i, k); len(err) > 0 {
			errors = append(errors, fmt.Errorf("expected %q to be an IPv4 IP Address, IPv6 IP Address or a valid domain name, got %q", k, i))
			return warn, errors
		}
	}

	return warnings, errors
}

func ValidateCdnFrontdoorName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z])([-\da-zA-Z]{0,88})([\da-zA-Z]$)`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 2 and 90 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens`, k))
	}
	return nil, nil
}

func ValidateCdnFrontdoorUrlRedirectActionDestinationPath(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("%q is invalid: expected type of %q to be string", "url_redirect_action", k)}
	}

	// TODO: we don't mention the below in the error, but we probably should?
	// Path cannot be empty and must start with /. Leave empty to use the incoming path as destination path.
	if v != "" {
		if !strings.HasPrefix(v, "/") {
			return nil, []error{fmt.Errorf("%q is invalid: %q must start with a %q, got %q", "url_redirect_action", k, "/", v)}
		}
	}

	return nil, nil
}

func ValidateCdnFrontdoorUrlRedirectActionQueryString(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("%q is invalid: expected type of %q to be string", "url_redirect_action", k)}
	}

	// Query string must be in <key>=<value> format. ? and & will be added automatically so do not include them.
	if v != "" {
		if strings.ContainsAny(v, "?&") {
			return nil, []error{fmt.Errorf("%q is invalid: %q must not include the %q or the %q characters in the %q field. They will be automatically added by Frontdoor, got %q", "url_redirect_action", k, "?", "&", "query_string", v)}
		}
		// ^(\b[a-zA-Z0-9\-\._~]*)(={1})(\b[a-zA-Z0-9\-\._~]*)$
		// ^(\b[a-zA-Z0-9\-\._~]*)(={1})((\b[a-zA-Z0-9\-\._~]*)|(\{{1}\b[a-zA-Z0-9\-\._~]*)\}{1})$
		if m, _ := validate.RegExHelper(i, k, `^(\b[a-zA-Z0-9\-\._~]*)(={1})((\b[a-zA-Z0-9\-\._~]*)|(\{{1}\b(socket_ip|client_ip|client_port|hostname|geo_country|http_method|http_version|query_string|request_scheme|request_uri|ssl_protocol|server_port|url_path){1}\}){1})$`); !m {
			return nil, []error{fmt.Errorf("%q is invalid: %q must be in the <key>=<value> or <key>={action_server_variable} format, got %q", "url_redirect_action", k, v)}
		}
	} else {
		return nil, []error{fmt.Errorf("%q is invalid: %q must not be empty, got %q", "url_redirect_action", k, v)}
	}

	return nil, nil
}

func ValidateCdnFrontdoorCacheDuration(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.HasPrefix(v, "0.") {
		return nil, []error{fmt.Errorf(`%q must not start with %q if the duration is less than 1 day. If the %q is less than 1 day it should be in the HH:MM:SS format, got %q`, k, "0.", k, v)}
	}

	if m, _ := validate.RegExHelper(i, k, `^([1-3]|([1-9][0-9])|([1-3][0-6][0-5])).((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$|^((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between in the d.HH:MM:SS or HH:MM:SS format and must be equal to or lower than %q, got %q`, k, "365.23:59:59", v)}
	}

	return nil, nil
}

func ValidateCdnFrontdoorUrlPathConditionMatchValue(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.HasPrefix(v, "/") {
		return nil, []error{fmt.Errorf(`%q must not start with the URLs paths leading slash(e.g. /), got %q`, k, v)}
	}

	return nil, nil
}

func ValidateCdnFrontdoorRuleName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^[a-zA-Z][a-zA-Z0-9]{0,259}$`); !m {
		return nil, []error{fmt.Errorf(`%q must start with a letter and contain only numbers and letters with a maximum length of 260 characters, got %q`, k, v)}
	}

	return nil, nil
}

func validCdnFrontdoorContentTypes() []string {
	return []string{
		"application/eot",
		"application/font",
		"application/font-sfnt",
		"application/javascript",
		"application/json",
		"application/opentype",
		"application/otf",
		"application/pkcs7-mime",
		"application/truetype",
		"application/ttf",
		"application/vnd.ms-fontobject",
		"application/xhtml+xml",
		"application/xml",
		"application/xml+rss",
		"application/x-font-opentype",
		"application/x-font-truetype",
		"application/x-font-ttf",
		"application/x-httpd-cgi",
		"application/x-mpegurl",
		"application/x-opentype",
		"application/x-otf",
		"application/x-perl",
		"application/x-ttf",
		"application/x-javascript",
		"font/eot",
		"font/ttf",
		"font/otf",
		"font/opentype",
		"image/svg+xml",
		"text/css",
		"text/csv",
		"text/html",
		"text/javascript",
		"text/js",
		"text/plain",
		"text/richtext",
		"text/tab-separated-values",
		"text/xml",
		"text/x-script",
		"text/x-component",
		"text/x-java-source",
	}
}

func validateActionsBlock(actions []cdn.BasicDeliveryRuleAction) error {
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

func ValidatedLegacyFrontdoorWAFName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[a-zA-Z])([\da-zA-Z]{0,127})$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 1 and 128 characters in length, must begin with a letter and may only contain letters and numbers`, k))
	}

	return nil, nil
}

func ValidateLegacyCustomBlockResponseBody(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q contains invalid characters, %q must contain only alphanumeric and equals sign characters`, k, k))
	}

	return nil, nil
}

func SchemaCdnFrontdoorOperator() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			string(cdn.OperatorAny),
			string(cdn.OperatorEqual),
			string(cdn.OperatorContains),
			string(cdn.OperatorBeginsWith),
			string(cdn.OperatorEndsWith),
			string(cdn.OperatorLessThan),
			string(cdn.OperatorLessThanOrEqual),
			string(cdn.OperatorGreaterThan),
			string(cdn.OperatorGreaterThanOrEqual),
			string(cdn.OperatorRegEx),
		}, false),
	}
}

func SchemaCdnFrontdoorOperatorEqualOnly() *pluginsdk.Schema {
	// TODO: if there's only one possible value, and it's defaulted - we don't need to expose this field for now?
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(cdn.OperatorEqual),
		ValidateFunc: validation.StringInSlice([]string{
			string(cdn.OperatorEqual),
		}, false),
	}
}

func SchemaCdnFrontdoorOperatorRemoteAddress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(cdn.OperatorIPMatch),
		ValidateFunc: validation.StringInSlice([]string{
			string(cdn.OperatorAny),
			string(cdn.OperatorIPMatch),
			string(cdn.OperatorGeoMatch),
		}, false),
	}
}

func SchemaCdnFrontdoorOperatorSocketAddress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(cdn.OperatorIPMatch),
		ValidateFunc: validation.StringInSlice([]string{
			string(cdn.OperatorAny),
			string(cdn.OperatorIPMatch),
		}, false),
	}
}

func SchemaCdnFrontdoorNegateCondition() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
}

func SchemaCdnFrontdoorMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 25,

		// In some cases it is valid for this to be an empty string
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func SchemaCdnFrontdoorServerPortMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: should this be a set?
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 2,

		// In some cases it is valid for this to be an empty string
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			// TODO: a list element can't have a default value
			Default: "80",
			ValidateFunc: validation.StringInSlice([]string{
				"80",
				"443",
			}, false),
		},
	}
}

func SchemaCdnFrontdoorSslProtocolMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: should this be a Set?
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 3,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			// TODO: a list element can't have a default value
			Default: string(cdn.SslProtocolTLSv12),
			ValidateFunc: validation.StringInSlice([]string{
				string(cdn.SslProtocolTLSv1),
				string(cdn.SslProtocolTLSv11),
				string(cdn.SslProtocolTLSv12),
			}, false),
		},
	}
}

func SchemaCdnFrontdoorUrlPathConditionMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 25,

		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: ValidateCdnFrontdoorUrlPathConditionMatchValue,
		},
	}
}

func SchemaCdnFrontdoorMatchValuesRequired() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 25,

		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func SchemaCdnFrontdoorRequestMethodMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: should this be a Set?
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 7,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				"GET",
				"POST",
				"PUT",
				"DELETE",
				"HEAD",
				"OPTIONS",
				"TRACE",
			}, false),
		},
	}
}

func SchemaCdnFrontdoorProtocolMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: if this is MaxItems: 1, should this be a string?
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,

		Elem: &pluginsdk.Schema{
			Type:    pluginsdk.TypeString,
			Default: "HTTP",
			ValidateFunc: validation.StringInSlice([]string{
				// TODO: are there constants for these?
				// TODO: other APIs use `Http` and `Https`, is that casing consistent in the API?
				"HTTP",
				"HTTPS",
			}, false),
		},
	}
}

func SchemaCdnFrontdoorIsDeviceMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			// TODO: an element can't have a default
			Default: "Mobile",
			ValidateFunc: validation.StringInSlice([]string{
				// TODO: are there constants for these?
				"Mobile",
				"Desktop",
			}, false),
		},
	}
}

func SchemaCdnFrontdoorHttpVersionMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: should this be a set?
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 4,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				"2.0",
				"1.1",
				"1.0",
				"0.9",
			}, false),
		},
	}
}

func SchemaCdnFrontdoorRuleTransforms() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: should this be a set?
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 4,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			// TODO: an element can't have a default value like this - besides which it's a weird default?
			Default: string(cdn.TransformLowercase),
			ValidateFunc: validation.StringInSlice([]string{
				string(cdn.TransformLowercase),
				string(cdn.TransformRemoveNulls),
				string(cdn.TransformTrim),
				string(cdn.TransformUppercase),
				string(cdn.TransformURLDecode),
				string(cdn.TransformURLEncode),
			}, false),
		},
	}
}
