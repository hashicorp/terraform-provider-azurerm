package securitysettings

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityComplianceStatus struct {
	DataAtRestEncrypted    *ComplianceStatus `json:"dataAtRestEncrypted,omitempty"`
	DataInTransitProtected *ComplianceStatus `json:"dataInTransitProtected,omitempty"`
	LastUpdated            *string           `json:"lastUpdated,omitempty"`
	SecuredCoreCompliance  *ComplianceStatus `json:"securedCoreCompliance,omitempty"`
	WdacCompliance         *ComplianceStatus `json:"wdacCompliance,omitempty"`
}

func (o *SecurityComplianceStatus) GetLastUpdatedAsTime() (*time.Time, error) {
	if o.LastUpdated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *SecurityComplianceStatus) SetLastUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdated = &formatted
}
