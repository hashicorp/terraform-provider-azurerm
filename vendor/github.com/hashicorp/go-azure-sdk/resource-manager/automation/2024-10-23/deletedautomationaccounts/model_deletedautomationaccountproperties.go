package deletedautomationaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedAutomationAccountProperties struct {
	AutomationAccountId         *string `json:"automationAccountId,omitempty"`
	AutomationAccountResourceId *string `json:"automationAccountResourceId,omitempty"`
	DeletionTime                *string `json:"deletionTime,omitempty"`
	Location                    *string `json:"location,omitempty"`
}

func (o *DeletedAutomationAccountProperties) GetDeletionTimeAsTime() (*time.Time, error) {
	if o.DeletionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeletionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DeletedAutomationAccountProperties) SetDeletionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeletionTime = &formatted
}
