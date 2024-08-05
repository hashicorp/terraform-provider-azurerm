package maintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListMaintenanceConfigurationsResult struct {
	Value *[]MaintenanceConfiguration `json:"value,omitempty"`
}
