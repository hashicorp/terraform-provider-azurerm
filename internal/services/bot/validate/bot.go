// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func BotMSTeamsCallingWebHook() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errs []error) {
		value := i.(string)
		if !strings.HasPrefix(value, "https://") || !strings.HasSuffix(value, "/") {
			errs = append(errs, errors.New("invalid `calling_web_hook`, must start with `https://` and end with `/`"))
		}

		return warnings, errs
	}
}
