package cacherules

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheRuleProperties struct {
	CreationDate            *string            `json:"creationDate,omitempty"`
	CredentialSetResourceId *string            `json:"credentialSetResourceId,omitempty"`
	ProvisioningState       *ProvisioningState `json:"provisioningState,omitempty"`
	SourceRepository        *string            `json:"sourceRepository,omitempty"`
	TargetRepository        *string            `json:"targetRepository,omitempty"`
}

func (o *CacheRuleProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *CacheRuleProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}
