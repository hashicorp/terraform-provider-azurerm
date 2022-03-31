package cdn

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	legacyfrontdoor "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/legacysdk/2020-11-01"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

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

func ConvertBoolToEnabledState(isEnabled bool) track1.EnabledState {
	out := track1.EnabledStateDisabled

	if isEnabled {
		out = track1.EnabledStateEnabled
	}

	return out
}

func ConvertEnabledStateToBool(enabledState *track1.EnabledState) bool {
	if enabledState == nil {
		return false
	}

	return (*enabledState == track1.EnabledStateEnabled)
}

func expandResourceReference(input string) *track1.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &track1.ResourceReference{
		ID: utils.String(input),
	}
}

// Takes a Slice of strings and transforms it into a CSV formatted string.
func ExpandStringSliceToCsvFormat(input []interface{}) *string {
	if len(input) == 0 {
		return nil
	}

	v := utils.ExpandStringSlice(input)
	csv := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(*v)), ","), "[]")

	return &csv
}

// Takes a CSV formatted string and transforms it into a Slice of strings.
func FlattenCsvToStringSlice(input *string) []interface{} {
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

func flattenResourceReference(input *track1.ResourceReference) string {
	result := ""
	if input == nil {
		return result
	}

	if input.ID != nil {
		result = *input.ID
	}

	return result
}

func flattenTransformSlice(input *[]legacyfrontdoor.TransformType) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

func flattenFrontendEndpointLinkSlice(input *[]legacyfrontdoor.FrontendEndpointLink) []interface{} {
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

func ConvertBoolToRouteHttpsRedirect(isEnabled bool) track1.HTTPSRedirect {
	out := track1.HTTPSRedirectDisabled

	if isEnabled {
		out = track1.HTTPSRedirectEnabled
	}

	return out
}

func ConvertRouteHttpsRedirectToBool(httpsRedirect *track1.HTTPSRedirect) bool {
	if httpsRedirect == nil {
		return false
	}

	return (*httpsRedirect == track1.HTTPSRedirect(track1.HTTPSRedirectEnabled))
}

func ConvertBoolToRouteLinkToDefaultDomain(isLinked bool) track1.LinkToDefaultDomain {
	out := track1.LinkToDefaultDomainDisabled

	if isLinked {
		out = track1.LinkToDefaultDomainEnabled
	}

	return out
}

func ConvertRouteLinkToDefaultDomainToBool(linkToDefaultDomain *track1.LinkToDefaultDomain) bool {
	if linkToDefaultDomain == nil {
		return false
	}

	return (*linkToDefaultDomain == track1.LinkToDefaultDomainEnabled)
}

func IsValidDomain(i interface{}, k string) (warnings []string, errors []error) {
	if warn, err := validation.IsIPv6Address(i, k); len(err) == 0 {
		return warn, err
	}

	if warn, err := validation.IsIPv4Address(i, k); len(err) == 0 {
		return warn, err
	}

	// TODO: Figure out a better way to validate Doman Name if not and IP Address
	if warn, err := validation.StringIsNotEmpty(i, k); len(err) == 0 {
		return warn, err
	}

	return warnings, errors
}

func ValidateCdnFrontdoorUrlRedirectActionDestinationPath(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("%q is invalid: expected type of %q to be string", "url_redirect_action", k)}
	}

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

func ValidateCdnFrontdoorRuleSetName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `(^[a-zA-Z])([\da-zA-Z]{1,88})([a-zA-Z]$)`); !m {
		return nil, []error{fmt.Errorf(`%q must be between 1 and 90 characters in length and begin with a letter, end with a letter and may contain only letters and numbers, got %q`, v, k)}
	}

	return nil, nil
}

func ValidateCdnFrontdoorCacheDuration(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^([0-3]|([1-9][0-9])|([1-3][0-6][0-5])).((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between in the d.HH:MM:SS format and must be equal to or lower than %q, got %q`, k, "365.23:59:59", v)}
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

func validateActionsBlock(actions []track1.BasicDeliveryRuleAction) error {
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
			string(track1.OperatorAny),
			string(track1.OperatorEqual),
			string(track1.OperatorContains),
			string(track1.OperatorBeginsWith),
			string(track1.OperatorEndsWith),
			string(track1.OperatorLessThan),
			string(track1.OperatorLessThanOrEqual),
			string(track1.OperatorGreaterThan),
			string(track1.OperatorGreaterThanOrEqual),
			string(track1.OperatorRegEx),
		}, false),
	}
}

func SchemaCdnFrontdoorOperatorEqualOnly() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(track1.OperatorEqual),
		ValidateFunc: validation.StringInSlice([]string{
			string(track1.OperatorEqual),
		}, false),
	}
}

func SchemaCdnFrontdoorOperatorRemoteAddress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(track1.OperatorIPMatch),
		ValidateFunc: validation.StringInSlice([]string{
			string(track1.OperatorAny),
			string(track1.OperatorIPMatch),
			string(track1.OperatorGeoMatch),
		}, false),
	}
}

func SchemaCdnFrontdoorOperatorSocketAddress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(track1.OperatorIPMatch),
		ValidateFunc: validation.StringInSlice([]string{
			string(track1.OperatorAny),
			string(track1.OperatorIPMatch),
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
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 2,

		// In some cases it is valid for this to be an empty string
		Elem: &pluginsdk.Schema{
			Type:    pluginsdk.TypeString,
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
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 3,

		Elem: &pluginsdk.Schema{
			Type:    pluginsdk.TypeString,
			Default: string(track1.SslProtocolTLSv12),
			ValidateFunc: validation.StringInSlice([]string{
				string(track1.SslProtocolTLSv1),
				string(track1.SslProtocolTLSv11),
				string(track1.SslProtocolTLSv12),
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
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,

		Elem: &pluginsdk.Schema{
			Type:    pluginsdk.TypeString,
			Default: "HTTP",
			ValidateFunc: validation.StringInSlice([]string{
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
			Type:    pluginsdk.TypeString,
			Default: "Mobile",
			ValidateFunc: validation.StringInSlice([]string{
				"Mobile",
				"Desktop",
			}, false),
		},
	}
}

func SchemaCdnFrontdoorHttpVersionMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 4,

		Elem: &pluginsdk.Schema{
			Type:    pluginsdk.TypeString,
			Default: string(track1.TransformLowercase),
			ValidateFunc: validation.StringInSlice([]string{
				string(track1.TransformLowercase),
				string(track1.TransformRemoveNulls),
				string(track1.TransformTrim),
				string(track1.TransformUppercase),
				string(track1.TransformURLDecode),
				string(track1.TransformURLEncode),
			}, false),
		},
	}
}
