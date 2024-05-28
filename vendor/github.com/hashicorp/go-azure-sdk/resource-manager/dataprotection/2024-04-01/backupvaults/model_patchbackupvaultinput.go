package backupvaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchBackupVaultInput struct {
	FeatureSettings                *FeatureSettings    `json:"featureSettings,omitempty"`
	MonitoringSettings             *MonitoringSettings `json:"monitoringSettings,omitempty"`
	ResourceGuardOperationRequests *[]string           `json:"resourceGuardOperationRequests,omitempty"`
	SecuritySettings               *SecuritySettings   `json:"securitySettings,omitempty"`
}
