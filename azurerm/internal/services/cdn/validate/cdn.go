package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func EndpointDeliveryRuleName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$"),
		"The Delivery Rule Name must start with a letter any may only contain letters and numbers.",
	)
}

func RuleActionCacheExpirationDuration() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^(\d+\.)?([0-1][0-9]|[2][0-3]):[0-5][0-9]:[0-5][0-9]$`),
		"The Cache duration must be in this format [d.]hh:mm:ss.",
	)
}

func RuleActionUrlRedirectPath() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^(/.*)?$"),
		"The Url Redirect Path must start with a slash.",
	)
}

func RuleActionUrlRedirectQueryString() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, s string) ([]string, []error) {
		querystring := i.(string)

		if len(querystring) > 100 {
			return nil, []error{fmt.Errorf("The Url Query String's max length is 100.")}
		}

		re := regexp.MustCompile("^[?&]")
		if re.MatchString(querystring) {
			return nil, []error{fmt.Errorf("The Url Query String must not start with a question mark or ampersand.")}
		}

		kvre := regexp.MustCompile("^[^?&]+=[^?&]+$")
		kvs := strings.Split(querystring, "&")
		for _, kv := range kvs {
			if len(kv) > 0 && !kvre.MatchString(kv) {
				return nil, []error{fmt.Errorf("The Url Query String must be in <key>=<value> format and separated by an ampersand.")}
			}
		}

		return nil, nil
	}
}

func RuleActionUrlRedirectFragment() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^([^#].*)?$"),
		"The Url Fragment must not start with a hash.",
	)
}

func RuleActionUrlRewriteSourcePattern() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^/[^\n]{0,259}$"),
		"The Url Rewrite Source Pattern must start with a slash and can not have more than 260 characters.",
	)
}

func RuleActionUrlRewriteDestination() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^/[^\n]{0,259}$"),
		"The Url Rewrite Destination must start with a slash and can not have more than 260 characters.",
	)
}
