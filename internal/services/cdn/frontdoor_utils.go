package cdn

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func ConvertFrontdoorTags(tagMap *map[string]string) map[string]*string {
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

func ConvertBoolToEnabledState(isEnabled bool) track1.EnabledState {
	out := track1.EnabledState(track1.EnabledStateDisabled)

	if isEnabled {
		out = track1.EnabledState(track1.EnabledStateEnabled)
	}

	return out
}

func ConvertEnabledStateToBool(enabledState *track1.EnabledState) bool {
	if enabledState == nil {
		return false
	}

	return (*enabledState == track1.EnabledState(track1.EnabledStateEnabled))
}

func expandResourceReference(input string) *track1.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &track1.ResourceReference{
		ID: utils.String(input),
	}
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

// func ConvertBoolToOriginsEnabledState(isEnabled bool) *afdorigins.EnabledState {
// 	out := afdorigins.EnabledState(afdorigins.EnabledStateDisabled)

// 	if isEnabled {
// 		out = afdorigins.EnabledState(afdorigins.EnabledStateEnabled)
// 	}

// 	return &out
// }

// func ConvertOriginsEnabledStateToBool(enabledState *afdorigins.EnabledState) bool {
// 	if enabledState == nil {
// 		return false
// 	}

// 	return (*enabledState == afdorigins.EnabledState(afdorigins.EnabledStateEnabled))
// }

func ConvertBoolToRouteHttpsRedirect(isEnabled bool) track1.HTTPSRedirect {
	out := track1.HTTPSRedirect(track1.HTTPSRedirectDisabled)

	if isEnabled {
		out = track1.HTTPSRedirect(track1.HTTPSRedirectEnabled)
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
	out := track1.LinkToDefaultDomain(track1.LinkToDefaultDomainDisabled)

	if isLinked {
		out = track1.LinkToDefaultDomain(track1.LinkToDefaultDomainEnabled)
	}

	return out
}

func ConvertRouteLinkToDefaultDomainToBool(linkToDefaultDomain *track1.LinkToDefaultDomain) bool {
	if linkToDefaultDomain == nil {
		return false
	}

	return (*linkToDefaultDomain == track1.LinkToDefaultDomain(track1.LinkToDefaultDomainEnabled))
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

func ValidateFrontdoorRuleSetName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, regexErrs := validate.RegExHelper(i, k, `(^[a-zA-Z])([\da-zA-Z]{1,88})([a-zA-Z]$)`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 1 and 90 characters in length and begin with a letter, end with a letter and may contain only letters and numbers, got %q`, v, k))
	}

	return nil, nil
}

func ValidateFrontdoorCacheDuration(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, regexErrs := validate.RegExHelper(i, k, `^([0-3]|([1-9][0-9])|([1-3][0-6][0-5])).((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between in the d.HH:MM:SS format and must be equal to or lower than %q, got %q`, v, "365.23:59:59", k))
	}

	return nil, nil
}

func ValidateFrontdoorUrlPathConditionMatchValue(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.HasPrefix(v, "/") {
		return nil, append(errors, fmt.Errorf(`%q must not start with the URLs paths leading slash(e.g. /), got %q`, k, v))
	}

	return nil, nil
}

func validContentTypes() []string {
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

func SchemaFrontdoorOperator() *pluginsdk.Schema {
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

func SchemaFrontdoorOperatorEqualOnly() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(track1.OperatorEqual),
		ValidateFunc: validation.StringInSlice([]string{
			string(track1.OperatorEqual),
		}, false),
	}
}

func SchemaFrontdoorNegateCondition() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
}

func SchemaFrontdoorMatchValues() *pluginsdk.Schema {
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

func SchemaFrontdoorUrlPathConditionMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 25,

		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: ValidateFrontdoorUrlPathConditionMatchValue,
		},
	}
}

func SchemaFrontdoorMatchValuesRequired() *pluginsdk.Schema {
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

func SchemaFrontdoorRequestMethodMatchValues() *pluginsdk.Schema {
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

func SchemaFrontdoorProtocolMatchValues() *pluginsdk.Schema {
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

func SchemaFrontdoorIsDeviceMatchValues() *pluginsdk.Schema {
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

func SchemaFrontdoorHttpVersionMatchValues() *pluginsdk.Schema {
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

func SchemaFrontdoorRuleTransforms() *pluginsdk.Schema {
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
