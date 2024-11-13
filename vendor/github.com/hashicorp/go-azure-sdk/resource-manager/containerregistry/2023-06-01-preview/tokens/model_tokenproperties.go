package tokens

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TokenProperties struct {
	CreationDate      *string                     `json:"creationDate,omitempty"`
	Credentials       *TokenCredentialsProperties `json:"credentials,omitempty"`
	ProvisioningState *ProvisioningState          `json:"provisioningState,omitempty"`
	ScopeMapId        *string                     `json:"scopeMapId,omitempty"`
	Status            *TokenStatus                `json:"status,omitempty"`
}

func (o *TokenProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *TokenProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}
