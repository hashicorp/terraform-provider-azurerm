package backupprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KPIResourceHealthDetails struct {
	ResourceHealthDetails *[]ResourceHealthDetails `json:"resourceHealthDetails,omitempty"`
	ResourceHealthStatus  *ResourceHealthStatus    `json:"resourceHealthStatus,omitempty"`
}
