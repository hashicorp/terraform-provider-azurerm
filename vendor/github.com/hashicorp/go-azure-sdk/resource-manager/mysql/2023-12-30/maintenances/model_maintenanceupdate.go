package maintenances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceUpdate struct {
	Properties *MaintenancePropertiesForUpdate `json:"properties,omitempty"`
}
