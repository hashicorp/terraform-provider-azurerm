package cloudexadatainfrastructures

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudExadataInfrastructureUpdateProperties struct {
	ComputeCount      *int64             `json:"computeCount,omitempty"`
	CustomerContacts  *[]CustomerContact `json:"customerContacts,omitempty"`
	DisplayName       *string            `json:"displayName,omitempty"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`
	StorageCount      *int64             `json:"storageCount,omitempty"`
}
