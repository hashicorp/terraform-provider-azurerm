package portalconfig

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortalConfigCspProperties struct {
	AllowedSources *[]string              `json:"allowedSources,omitempty"`
	Mode           *PortalSettingsCspMode `json:"mode,omitempty"`
	ReportUri      *[]string              `json:"reportUri,omitempty"`
}
