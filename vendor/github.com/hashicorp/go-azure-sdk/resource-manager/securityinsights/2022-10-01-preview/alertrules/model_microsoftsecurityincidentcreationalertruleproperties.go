package alertrules

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MicrosoftSecurityIncidentCreationAlertRuleProperties struct {
	AlertRuleTemplateName     *string                      `json:"alertRuleTemplateName,omitempty"`
	Description               *string                      `json:"description,omitempty"`
	DisplayName               string                       `json:"displayName"`
	DisplayNamesExcludeFilter *[]string                    `json:"displayNamesExcludeFilter,omitempty"`
	DisplayNamesFilter        *[]string                    `json:"displayNamesFilter,omitempty"`
	Enabled                   bool                         `json:"enabled"`
	LastModifiedUtc           *string                      `json:"lastModifiedUtc,omitempty"`
	ProductFilter             MicrosoftSecurityProductName `json:"productFilter"`
	SeveritiesFilter          *[]AlertSeverity             `json:"severitiesFilter,omitempty"`
}

func (o *MicrosoftSecurityIncidentCreationAlertRuleProperties) GetLastModifiedUtcAsTime() (*time.Time, error) {
	if o.LastModifiedUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *MicrosoftSecurityIncidentCreationAlertRuleProperties) SetLastModifiedUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedUtc = &formatted
}
