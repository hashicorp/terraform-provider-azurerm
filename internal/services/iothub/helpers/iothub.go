// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func FileUploadConnectionStringDiffSuppress(k, old, new string, d *pluginsdk.ResourceData) bool {
	// The access keys are always masked by Azure and the ordering of the parameters in the connection string
	// differs across services, so we will compare the fields individually instead.
	secretKeyRegex := regexp.MustCompile("(SharedAccessKey|AccountKey)=[^;]+")

	if secretKeyRegex.MatchString(new) {
		maskedNew := secretKeyRegex.ReplaceAllString(new, "$1=****")

		oldSplit := strings.Split(old, ";")
		newSplit := strings.Split(maskedNew, ";")

		sort.Strings(oldSplit)
		sort.Strings(newSplit)

		if len(oldSplit) != len(newSplit) {
			return false
		}

		for i := range oldSplit {
			if !strings.EqualFold(oldSplit[i], newSplit[i]) {
				return false
			}
		}
		return true
	}
	return false
}
