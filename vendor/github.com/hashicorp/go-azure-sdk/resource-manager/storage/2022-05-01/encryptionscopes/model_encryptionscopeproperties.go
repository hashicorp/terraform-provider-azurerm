package encryptionscopes

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionScopeProperties struct {
	CreationTime                    *string                            `json:"creationTime,omitempty"`
	KeyVaultProperties              *EncryptionScopeKeyVaultProperties `json:"keyVaultProperties,omitempty"`
	LastModifiedTime                *string                            `json:"lastModifiedTime,omitempty"`
	RequireInfrastructureEncryption *bool                              `json:"requireInfrastructureEncryption,omitempty"`
	Source                          *EncryptionScopeSource             `json:"source,omitempty"`
	State                           *EncryptionScopeState              `json:"state,omitempty"`
}

func (o *EncryptionScopeProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *EncryptionScopeProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *EncryptionScopeProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *EncryptionScopeProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
