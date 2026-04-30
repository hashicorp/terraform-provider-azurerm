package nginxdeploymentapikeyresponses

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentApiKeyRequestProperties struct {
	EndDateTime *string `json:"endDateTime,omitempty"`
	SecretText  *string `json:"secretText,omitempty"`
}

func (o *NginxDeploymentApiKeyRequestProperties) GetEndDateTimeAsTime() (*time.Time, error) {
	if o.EndDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *NginxDeploymentApiKeyRequestProperties) SetEndDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDateTime = &formatted
}
