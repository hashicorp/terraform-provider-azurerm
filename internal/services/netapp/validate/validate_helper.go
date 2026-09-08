// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
)

type ProtocolType string

const (
	ProtocolTypeNfsV41 ProtocolType = "NFSv4.1"
	ProtocolTypeNfsV3  ProtocolType = "NFSv3"
	ProtocolTypeCifs   ProtocolType = "CIFS"
)

func PossibleValuesForProtocolType() []string {
	return []string{
		string(ProtocolTypeNfsV41),
		string(ProtocolTypeNfsV3),
		string(ProtocolTypeCifs),
	}
}

func findStringInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}
