// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import "encoding/json"

var _ json.Marshaler = UserAssignedIdentityDetails{}

type UserAssignedIdentityDetails struct {
	ClientId    *string `json:"clientId,omitempty"`
	PrincipalId *string `json:"principalId,omitempty"`
}

func (u UserAssignedIdentityDetails) MarshalJSON() ([]byte, error) {
	// none of these properties can be set, so we'll just flatten an empty struct
	return json.Marshal(map[string]interface{}{})
}
