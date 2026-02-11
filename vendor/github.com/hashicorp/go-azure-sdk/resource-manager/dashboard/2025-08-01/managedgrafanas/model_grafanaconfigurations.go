package managedgrafanas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GrafanaConfigurations struct {
	Security                   *Security                   `json:"security,omitempty"`
	Smtp                       *Smtp                       `json:"smtp,omitempty"`
	Snapshots                  *Snapshots                  `json:"snapshots,omitempty"`
	UnifiedAlertingScreenshots *UnifiedAlertingScreenshots `json:"unifiedAlertingScreenshots,omitempty"`
	Users                      *Users                      `json:"users,omitempty"`
}
