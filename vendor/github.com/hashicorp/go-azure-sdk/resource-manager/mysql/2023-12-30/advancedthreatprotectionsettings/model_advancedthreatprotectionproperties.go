package advancedthreatprotectionsettings

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvancedThreatProtectionProperties struct {
	CreationTime      *string                                    `json:"creationTime,omitempty"`
	ProvisioningState *AdvancedThreatProtectionProvisioningState `json:"provisioningState,omitempty"`
	State             *AdvancedThreatProtectionState             `json:"state,omitempty"`
}

func (o *AdvancedThreatProtectionProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AdvancedThreatProtectionProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}
