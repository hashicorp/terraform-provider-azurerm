package autoexportjob

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoExportJobProperties struct {
	AdminStatus        *AutoExportJobAdminStatus           `json:"adminStatus,omitempty"`
	AutoExportPrefixes *[]string                           `json:"autoExportPrefixes,omitempty"`
	ProvisioningState  *AutoExportJobProvisioningStateType `json:"provisioningState,omitempty"`
	Status             *AutoExportJobPropertiesStatus      `json:"status,omitempty"`
}
