// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package user_assigned_identity

func (r Resource) ModelObject() any {
	return &UserAssignedIdentityResourceSchema{}
}
