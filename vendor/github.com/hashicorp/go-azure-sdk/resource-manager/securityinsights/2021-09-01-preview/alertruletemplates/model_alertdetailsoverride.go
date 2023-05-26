package alertruletemplates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertDetailsOverride struct {
	AlertDescriptionFormat  *string `json:"alertDescriptionFormat,omitempty"`
	AlertDisplayNameFormat  *string `json:"alertDisplayNameFormat,omitempty"`
	AlertSeverityColumnName *string `json:"alertSeverityColumnName,omitempty"`
	AlertTacticsColumnName  *string `json:"alertTacticsColumnName,omitempty"`
}
