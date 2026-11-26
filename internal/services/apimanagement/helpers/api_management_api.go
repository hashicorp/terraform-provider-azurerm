// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"strings"
)

func GetApiName(apiId string) string {
	name := apiId
	if strings.Contains(apiId, ";") {
		name = strings.Split(apiId, ";")[0]
	}

	return name
}
