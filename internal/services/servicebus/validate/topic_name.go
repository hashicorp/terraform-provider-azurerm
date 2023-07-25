// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func TopicName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9]([-._~/a-zA-Z0-9]{0,258}[a-zA-Z0-9])?$"),
		"The topic name can contain only letters, numbers, periods, hyphens, tildas, forward slashes and underscores. The namespace must start with a letter or number, and it must end with a letter or number and be less then 260 characters long.",
	)
}
