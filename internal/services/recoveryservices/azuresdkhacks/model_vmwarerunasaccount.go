// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/runasaccounts"

type VMwareRunAsAccount struct {
	Id         *string                 `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *RunAsAccountProperties `json:"properties,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}
type RunAsAccountProperties struct {
	CreatedTimestamp *string                       `json:"createdTimestamp,omitempty"`
	CredentialType   *runasaccounts.CredentialType `json:"credentialType,omitempty"`
	DisplayName      *string                       `json:"displayName,omitempty"`
	UpdatedTimestamp *string                       `json:"updatedTimestamp,omitempty"`
	// ApplianceName was not defined in the original code. This is the only change.
	ApplianceName *string `json:"applianceName,omitempty"`
}
