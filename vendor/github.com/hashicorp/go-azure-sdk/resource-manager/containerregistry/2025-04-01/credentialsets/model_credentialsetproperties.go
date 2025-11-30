package credentialsets

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CredentialSetProperties struct {
	AuthCredentials   *[]AuthCredential  `json:"authCredentials,omitempty"`
	CreationDate      *string            `json:"creationDate,omitempty"`
	LoginServer       *string            `json:"loginServer,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

func (o *CredentialSetProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *CredentialSetProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}
