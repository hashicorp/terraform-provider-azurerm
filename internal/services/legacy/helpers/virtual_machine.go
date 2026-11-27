// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
)

func UserDataStateFunc(v any) string {
	switch s := v.(type) {
	case string:
		s = base64.EncodeIfNot(s)
		hash := sha1.Sum([]byte(s))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}
