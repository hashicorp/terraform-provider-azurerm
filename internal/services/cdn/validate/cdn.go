// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func EndpointDeliveryRuleName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$"),
		"the Delivery Rule Name must start with a letter any may only contain letters and numbers",
	)
}

func RuleActionCacheExpirationDuration() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^(\d+\.)?([0-1][0-9]|[2][0-3]):[0-5][0-9]:[0-5][0-9]$`),
		"the Cache duration must be in this format [d.]hh:mm:ss",
	)
}

func RuleActionUrlRedirectPath() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^(/.*)?$"),
		"the Url Redirect Path must start with a slash",
	)
}

func RuleActionUrlRedirectQueryString() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, s string) ([]string, []error) {
		querystring := i.(string)

		re := regexp.MustCompile("^[?&]")
		if re.MatchString(querystring) {
			return nil, []error{fmt.Errorf("the Url Query String must not start with a question mark or ampersand")}
		}

		kvre := regexp.MustCompile("^[^?&]+=[^?&]+$")
		kvs := strings.Split(querystring, "&")
		for _, kv := range kvs {
			if len(kv) > 0 && !kvre.MatchString(kv) {
				return nil, []error{fmt.Errorf("the Url Query String must be in <key>=<value> format and separated by an ampersand")}
			}
		}

		return nil, nil
	}
}

func RuleActionUrlRedirectFragment() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^([^#].*)?$"),
		"the Url Fragment must not start with a hash.",
	)
}

func RuleActionUrlRewriteSourcePattern() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^/[^\n]{0,259}$"),
		"the Url Rewrite Source Pattern must start with a slash and can not have more than 260 characters",
	)
}

func RuleActionUrlRewriteDestination() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^/[^\n]{0,259}$"),
		"the Url Rewrite Destination must start with a slash and can not have more than 260 characters",
	)
}
