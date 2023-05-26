package alertruletemplates

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties struct {
	AlertRulesCreatedByTemplateCount int64                          `json:"alertRulesCreatedByTemplateCount"`
	CreatedDateUTC                   string                         `json:"createdDateUTC"`
	Description                      string                         `json:"description"`
	DisplayName                      string                         `json:"displayName"`
	DisplayNamesExcludeFilter        *[]string                      `json:"displayNamesExcludeFilter,omitempty"`
	DisplayNamesFilter               *[]string                      `json:"displayNamesFilter,omitempty"`
	LastUpdatedDateUTC               *string                        `json:"lastUpdatedDateUTC,omitempty"`
	ProductFilter                    MicrosoftSecurityProductName   `json:"productFilter"`
	RequiredDataConnectors           *[]AlertRuleTemplateDataSource `json:"requiredDataConnectors,omitempty"`
	SeveritiesFilter                 *[]AlertSeverity               `json:"severitiesFilter,omitempty"`
	Status                           TemplateStatus                 `json:"status"`
}

func (o *MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties) GetCreatedDateUTCAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.CreatedDateUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties) SetCreatedDateUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDateUTC = formatted
}

func (o *MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties) GetLastUpdatedDateUTCAsTime() (*time.Time, error) {
	if o.LastUpdatedDateUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedDateUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties) SetLastUpdatedDateUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedDateUTC = &formatted
}
