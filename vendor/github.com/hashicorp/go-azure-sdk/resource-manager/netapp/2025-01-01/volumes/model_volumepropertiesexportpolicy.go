package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumePropertiesExportPolicy struct {
	Rules *[]ExportPolicyRule `json:"rules,omitempty"`
}
